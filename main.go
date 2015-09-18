package main

import (
	"fmt"
	"net/http"

	log "github.com/carlosdp/harbor/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/carlosdp/harbor/chain"
	"github.com/carlosdp/harbor/config"
	"github.com/carlosdp/harbor/deployment"
	"github.com/carlosdp/harbor/hook"
)

func main() {
	mux := http.NewServeMux()
	chains, err := config.ParseConfig("config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, c := range chains {
		for _, link := range c.Links {
			log.Info("Link loaded: ", link.Link.Name())
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
		Addr:    ":3002",
		Handler: mux,
	}

	log.Info("Starting Server")
	server.ListenAndServe()
}
