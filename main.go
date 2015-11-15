package main

import (
	"flag"
	"net/http"
	"os"
	"path"

	log "github.com/carlosdp/supply-chain/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/carlosdp/supply-chain/chain"
	"github.com/carlosdp/supply-chain/config"
	"github.com/carlosdp/supply-chain/hook"
)

var port string
var configPath string
var dataPath string

func init() {
	flag.StringVar(&port, "p", "3001", "Port webhooks should listen on.")
	flag.StringVar(&configPath, "c", "", "Path to chain config file.")
	flag.StringVar(&dataPath, "data", "", "Path where persistence data should be stored.")
}

func main() {
	flag.Parse()

	dataPath = path.Clean(dataPath)
	if dataPath != "/" {
		dataPath = dataPath + "/"
	}

	if configPath == "" {
		log.Error("[Config] You must specify a config file with -c")
		os.Exit(1)
	}
	configFile, err := os.Open(configPath)
	if err != nil {
		log.Errorf("[Config] %v", err)
		os.Exit(1)
	}
	chains, err := config.ParseConfig(configFile)
	if err != nil {
		log.Errorf("[Config] %v", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	requestChan := make(chan hook.Request)

	for _, c := range chains {
		log.Infof("[%v] Loading chain", c.Name)
		for _, link := range c.Links {
			log.Infof("[%v] Link loaded: %v", c.Name, link.Link.Name())
		}

		c.Load(dataPath)

		for _, hookLink := range c.LinksOfType(chain.HOOK) {
			hookWrap := hookLink.Link.(*hook.Wrapper)
			go hookWrap.Start(mux, requestChan, hookLink.Options, c, hookLink)
		}
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Info("Starting Server")
	go server.ListenAndServe()

	for {
		select {
		case request := <-requestChan:
			deploy, err := chain.NewDeployment(request.Chain, request.Link)
			if err != nil {
				log.Errorf("[Deployment] %v", err)
				return
			}
			deploy.SetName(request.Name)
			deploy.SetID(request.DeploymentID)
			deploy.SetURI(request.URI)
			deploy.SetImage(request.Image)

			go func(deploy *chain.Deployment) {
				err := deploy.Run()
				if err != nil {
					log.Errorf("[Deployment] %v", err)
					return
				}
				deploy.Chain.Persist(dataPath)
			}(deploy)
		}
	}
}
