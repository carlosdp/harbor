package main

import (
	"fmt"
	"github.com/carlosdp/harbor/container"
	"github.com/carlosdp/harbor/git-puller"
	"github.com/carlosdp/harbor/github-hook"
	"github.com/carlosdp/harbor/image"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		hook := &githubhook.GithubHook{}
		hook.HandleRequest(request)
		workDir := gitpuller.Pull(hook.URI())
		image.NewImage(hook.Name(), workDir)
		container.NewContainer(hook.Name(), hook.DeploymentID())
	})

	server := &http.Server{
		Addr:    ":3002",
		Handler: mux,
	}

	fmt.Println("starting server")
	server.ListenAndServe()
}
