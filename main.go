package main

import (
	"fmt"
	"net/http"

	"github.com/carlosdp/harbor/chain"
	"github.com/carlosdp/harbor/deployment"
	"github.com/carlosdp/harbor/hook"
)

func main() {
	mux := http.NewServeMux()
	chains, err := ParseConfig("config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, c := range chains {
		for _, link := range c.Links {
			fmt.Println("Link loaded: ", link.Link.Name())
		}
		for _, hookLink := range c.LinksOfType(chain.HOOK) {
			hookWrap := hookLink.Link.(*hook.HookWrapper)
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

	fmt.Println("starting server")
	server.ListenAndServe()
}
