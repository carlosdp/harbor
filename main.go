package main

import (
	"flag"
	"fmt"
	"net/http"

	log "github.com/carlosdp/harbor/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/carlosdp/harbor/chain"
	"github.com/carlosdp/harbor/config"
	"github.com/carlosdp/harbor/deployment"
	"github.com/carlosdp/harbor/hook"
)

var port string

func init() {
	flag.StringVar(&port, "p", "3001", "The port webhooks should listen on.")
}

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	chains, err := config.ParseConfig("config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

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
					panic(err)
				}

				err = hookWrap.HandleRequest(deploy, req)
				if err != nil {
					panic(err)
				}

				go func() {
					err := deploy.Run()
					if err != nil {
						fmt.Println(err)
					}
				}()

				fmt.Println(deploy.ID())
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
