package deployment

import (
	"errors"
	"fmt"

	"github.com/carlosdp/harbor/chain"
)

// Deployment is a single deployment that can be run.
type Deployment struct {
	Chain          *chain.Chain
	StartStep      int
	CurrentStep    int
	CompletedLinks []*chain.Link

	uri     string
	name    string
	id      string
	workDir string
	image   string
}

// NewDeployment creates a deployment from a `chain` and `hookLink`
// and returns a deployment that is ready to be run or rolled back.
func NewDeployment(dChain *chain.Chain, hookLink *chain.Link) (*Deployment, error) {
	currentStep, err := dChain.LinkPosition(hookLink)
	if err != nil {
		return nil, err
	}

	d := &Deployment{
		Chain:          dChain,
		StartStep:      currentStep,
		CurrentStep:    currentStep + 1,
		CompletedLinks: make([]*chain.Link, 0),
	}

	return d, nil
}

// Run runs the deployment through the chain.
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

// Rollback executes a rolling back of the deployment backward
// through the chain.
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

// URI is the uri of the resource being deployed.
func (d *Deployment) URI() string {
	return d.uri
}

// Name is the name of the resource being deployed.
func (d *Deployment) Name() string {
	return d.name
}

// ID is the deployment id that identifies the deployment for the
// particular version of the resource being deployed.
func (d *Deployment) ID() string {
	return d.id
}

// WorkDir is the local working directory containing the resource
// being deployed.
func (d *Deployment) WorkDir() string {
	return d.workDir
}

// Image is the image identifier for the built artifact.
func (d *Deployment) Image() string {
	return d.image
}

// SetURI sets the uri for the resource being deployed.
func (d *Deployment) SetURI(uri string) {
	d.uri = uri
}

// SetName sets the name of the resource being deployed.
func (d *Deployment) SetName(name string) {
	d.name = name
}

// SetID sets the deployment id.
func (d *Deployment) SetID(id string) {
	d.id = id
}

// SetWorkDir sets the local working directory for the resource
// being deployed.
func (d *Deployment) SetWorkDir(workDir string) {
	d.workDir = workDir
}

// SetImage sets the image identifier for the built artifact.
func (d *Deployment) SetImage(image string) {
	d.image = image
}
