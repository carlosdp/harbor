package main

import (
	"flag"
	"net/http"
	"os"

	log "github.com/carlosdp/harbor/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/carlosdp/harbor/chain"
	"github.com/carlosdp/harbor/config"
	"github.com/carlosdp/harbor/deployment"
	"github.com/carlosdp/harbor/hook"
)

var port string
var configPath string

func init() {
	flag.StringVar(&port, "p", "3001", "The port webhooks should listen on.")
	flag.StringVar(&configPath, "c", "", "Path to chain config file")
}

func main() {
	flag.Parse()

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

	for _, c := range chains {
		log.Infof("[%v] Loading chain", c.Name)
		for _, link := range c.Links {
			log.Infof("[%v] Link loaded: %v", c.Name, link.Link.Name())
		}
		for _, hookLink := range c.LinksOfType(chain.HOOK) {
			hookWrap := hookLink.Link.(*hook.Wrapper)
			mux.HandleFunc(hookWrap.Endpoint, func(res http.ResponseWriter, req *http.Request) {
				deploy, err := deployment.NewDeployment(c, hookLink)
				if err != nil {
					log.Errorf("[Deployment] %v", err)
					return
				}

				err = hookWrap.HandleRequest(deploy, req)
				if err != nil {
					log.Errorf("[Deployment] %v", err)
					return
				}

				go func() {
					err := deploy.Run()
					if err != nil {
						log.Errorf("[Deployment] %v", err)
						return
					}
				}()

				res.WriteHeader(200)
			})
		}
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Info("Starting Server")
	server.ListenAndServe()
}
