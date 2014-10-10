package deployment

import (
	"github.com/carlosdp/harbor/chain"
)

type Deployment struct {
	Chain          *chain.Chain
	StartStep      int
	CurrentStep    int
	CompletedLinks []*chain.ChainLink

	uri     string
	name    string
	id      string
	workDir string
	image   string
}

func NewDeployment(dChain *chain.Chain, hookLink *chain.ChainLink) (*Deployment, error) {
	currentStep, err := dChain.LinkPosition(hookLink)
	if err != nil {
		return nil, err
	}

	d := &Deployment{
		Chain:          dChain,
		StartStep:      currentStep,
		CurrentStep:    currentStep + 1,
		CompletedLinks: make([]*chain.ChainLink, 0),
	}

	return d, nil
}

func (d *Deployment) URI() string {
	return d.uri
}

func (d *Deployment) Name() string {
	return d.name
}

func (d *Deployment) ID() string {
	return d.id
}

func (d *Deployment) WorkDir() string {
	return d.workDir
}

func (d *Deployment) Image() string {
	return d.image
}

func (d *Deployment) SetURI(uri string) {
	d.uri = uri
}

func (d *Deployment) SetName(name string) {
	d.name = name
}

func (d *Deployment) SetID(id string) {
	d.id = id
}

func (d *Deployment) SetWorkDir(workDir string) {
	d.workDir = workDir
}

func (d *Deployment) SetImage(image string) {
	d.image = image
}
