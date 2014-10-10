package main

import (
	"fmt"
	// "github.com/carlosdp/harbor/container"
	// "github.com/carlosdp/harbor/git-puller"
	// "github.com/carlosdp/harbor/image"
	"github.com/carlosdp/harbor/chain"
	"github.com/carlosdp/harbor/deployment"
	"github.com/carlosdp/harbor/hook"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	// mux.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
	// 	hook := &githubhook.GithubHook{}
	// 	hook.HandleRequest(request)
	// 	workDir := gitpuller.Pull(hook.URI())
	// 	image.NewImage(hook.Name(), workDir)
	// 	container.NewContainer(hook.Name(), hook.DeploymentID())
	// })

	chains, err := ParseConfig("config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, c := range chains {
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
