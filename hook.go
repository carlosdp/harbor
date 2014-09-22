package main

import (
	"net/http"
)

type Hook interface {
	HandleRequest(req *http.Request) string
	Name() string
	DeploymentID() string
	URI() string
}
