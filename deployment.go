package main

import (
	"github.com/carlosdp/harbor/hook"
)

type Deployment struct {
	Chain          *Chain
	CurrentStep    int
	CompletedLinks []*ChainLink

	URI string
	ID  string
}

func NewDeployment(chain *Chain, hookLink *ChainLink) (*Deployment, error) {
	currentStep, err := chain.LinkPosition(hookLink)
	if err != nil {
		return nil, err
	}

	hook := hookLink.Link.(hook.Hook)
	uri := hook.URI()
	id := hook.DeploymentID()

	d := &Deployment{
		Chain:          chain,
		CurrentStep:    currentStep + 1,
		CompletedLinks: make([]*ChainLink, 0),
		URI:            uri,
		ID:             id,
	}

	return d, nil
}
