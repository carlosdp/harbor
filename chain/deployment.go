package chain

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/carlosdp/harbor/options"
)

// Deployment is a single deployment that can be run.
type Deployment struct {
	Chain       *Chain `json:"-"`
	StartStep   int
	CurrentStep int
	State       map[string]interface{}

	URI     string
	Name    string
	ID      string
	WorkDir string
	Image   string
	failure bool
}

// NewDeployment creates a deployment from a `chain` and `hookLink`
// and returns a deployment that is ready to be run or rolled back.
func NewDeployment(dChain *Chain, hookLink *Link) (*Deployment, error) {
	currentStep, err := dChain.LinkPosition(hookLink)
	if err != nil {
		return nil, err
	}

	d := &Deployment{
		Chain:       dChain,
		StartStep:   currentStep,
		CurrentStep: currentStep + 1,
		State:       make(map[string]interface{}),
	}

	return d, nil
}

// Run runs the deployment through the chain.
func (d *Deployment) Run() error {
	n := d.CurrentStep

	for i := d.CurrentStep; i < len(d.Chain.Links)-n+1; i++ {
		link := d.Chain.Links[i]
		log.Info("Running ", link.Link.Name())
		err := link.Link.Execute(d, link.Options)
		if err != nil {
			d.failure = true
			rerr := d.Rollback()
			if rerr != nil {
				return errors.New("Deployment failed: " + err.Error() + ", Failed to Rollback: " + rerr.Error())
			}

			return errors.New("Deployment failed: " + err.Error())
		}

		log.Info("Ran ", link.Link.Name())

		d.CurrentStep++
	}

	log.Infof("Deployment %v complete, rolling back active deployments", d.ID)

	activeDeploys := d.Chain.ActiveDeployments

	for _, deploy := range activeDeploys {
		log.Infof("Rolling back deployment %v", deploy.ID)
		deploy.Rollback()
	}

	d.Chain.ActiveDeployments = append(d.Chain.ActiveDeployments, d)

	log.Info("Deployment chain complete")

	return nil
}

// Rollback executes a rolling back of the deployment backward
// through the chain.
func (d *Deployment) Rollback() error {
	d.CurrentStep--

	var rerr error
	index := d.GetIndex()

	for i := d.CurrentStep; i > d.StartStep; i-- {
		link := d.Chain.Links[i]
		keep := link.Parameters.GetInt("keep")
		if !d.failure && keep > 0 && len(d.Chain.ActiveDeployments) >= keep && index < keep {
			return rerr
		}
		err := link.Link.Rollback(d, link.Options)
		if err != nil {
			log.Error("Rollback unsuccessful: " + err.Error())
			rerr = err
		}
		d.CurrentStep--
	}

	if !d.failure && index >= 0 {
		d.Chain.ActiveDeployments = append(d.Chain.ActiveDeployments[:index], d.Chain.ActiveDeployments[index+1:]...)
	}

	return rerr
}

// SetURI sets the uri for the resource being deployed.
func (d *Deployment) SetURI(uri string) {
	d.URI = uri
}

// SetName sets the name of the resource being deployed.
func (d *Deployment) SetName(name string) {
	d.Name = name
}

// SetID sets the deployment id.
func (d *Deployment) SetID(id string) {
	d.ID = id
}

// SetWorkDir sets the local working directory for the resource
// being deployed.
func (d *Deployment) SetWorkDir(workDir string) {
	d.WorkDir = workDir
}

// SetImage sets the image identifier for the built artifact.
func (d *Deployment) SetImage(image string) {
	d.Image = image
}

// GetIndex indicates where the deployment lies in the deployment log.
func (d *Deployment) GetIndex() int {
	for i, dp := range d.Chain.ActiveDeployments {
		if dp.Name == d.Name && dp.ID == d.ID {
			return i
		}
	}

	return -1
}

// SetState sets the state for a link in a chain for use in a rollback.
func (d *Deployment) SetState(linkName string, state interface{}) {
	d.State[linkName] = state
}

// GetState gets the state for a link as an option.Option.
func (d *Deployment) GetState(linkName string) options.Option {
	state, ok := d.State[linkName]
	if !ok {
		return options.NewOption(nil)
	}

	return options.NewOption(state)
}
