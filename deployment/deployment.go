package deployment

import (
	"github.com/carlosdp/harbor/chain"
	"github.com/carlosdp/harbor/hook"
)

type Deployment struct {
	Chain          *chain.Chain
	CurrentStep    int
	CompletedLinks []*chain.ChainLink

	URI     string
	ID      string
	WorkDir string
}

func NewDeployment(dChain *chain.Chain, hookLink *chain.ChainLink) (*Deployment, error) {
	currentStep, err := dChain.LinkPosition(hookLink)
	if err != nil {
		return nil, err
	}

	hook := hookLink.Link.(hook.Hook)
	uri := hook.URI()
	id := hook.DeploymentID()

	d := &Deployment{
		Chain:          dChain,
		CurrentStep:    currentStep + 1,
		CompletedLinks: make([]*chain.ChainLink, 0),
		URI:            uri,
		ID:             id,
	}

	return d, nil
}

func (d *Deployment) URI() string {
	return d.URI
}

func (d *Deployment) ID() string {
	return d.ID
}

func (d *Deployment) WorkDir() string {
	return d.WorkDir
}

func (d *Deployment) SetURI(uri string) {
	d.URI = uri
}

func (d *Deployment) SetID(id string) {
	d.ID = id
}

func (d *Deployment) SetWorkDir(workDir string) {
	d.WorkDir = workDir
}
