package main

import (
	"net/http"
)

type Hook interface {
	Endpoint() string
	HandleRequest(req *http.Request) string
	Name() string
	DeploymentID() string
	URI() string
}
