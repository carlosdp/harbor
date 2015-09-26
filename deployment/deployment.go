package deployment

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/carlosdp/harbor/chain"
	"github.com/carlosdp/harbor/options"
)

// Deployment is a single deployment that can be run.
type Deployment struct {
	Chain          *chain.Chain
	StartStep      int
	CurrentStep    int
	CompletedLinks []*chain.Link
	State          map[string]interface{}

	uri     string
	name    string
	id      string
	workDir string
	image   string
	failure bool
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
		State:          make(map[string]interface{}),
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

	log.Infof("Deployment %v complete, rolling back active deployments", d.ID())

	activeDeploys := d.Chain.ActiveDeployments

	for _, deploy := range activeDeploys {
		dp, ok := deploy.(*Deployment)
		if !ok {
			log.Error("Invalid deployment in deployment log")
		}

		log.Infof("Rolling back deployment %v", dp.ID())
		dp.Rollback()
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
		if !d.failure && keep > 0 && len(d.Chain.ActiveDeployments) > keep && index > keep {
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

// GetIndex indicates where the deployment lies in the deployment log.
func (d *Deployment) GetIndex() int {
	for i, dp := range d.Chain.ActiveDeployments {
		if dp.Name() == d.Name() && dp.ID() == d.ID() {
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
