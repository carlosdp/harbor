package deployment

import (
	"errors"
	"fmt"

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

func (d *Deployment) Run() error {
	n := d.CurrentStep

	for i := d.CurrentStep; i < len(d.Chain.Links)-n+1; i++ {
		link := d.Chain.Links[i]
		fmt.Println("Running ", link.Link.Name())
		err := link.Link.Execute(d)
		if err != nil {
			rerr := d.Rollback()
			if rerr != nil {
				return errors.New("Deployment failed: " + err.Error() + ", Failed to Rollback: " + rerr.Error())
			}

			return errors.New("Deployment failed: " + err.Error())
		}

		fmt.Println("Ran ", link.Link.Name())

		d.CurrentStep++
	}

	return nil
}

func (d *Deployment) Rollback() error {
	d.CurrentStep--

	var rerr error

	for i := d.CurrentStep; i > d.StartStep; i-- {
		link := d.Chain.Links[i]
		err := link.Link.Rollback()
		if err != nil {
			fmt.Println("Rollback unsuccessful: " + err.Error())
			rerr = err
		}
	}

	return rerr
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
