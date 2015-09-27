package puller

import (
	"errors"

	"github.com/carlosdp/harbor/chain"
	"github.com/carlosdp/harbor/options"
)

// RegisteredPullers contains the pullers registered with the Harbor build.
var RegisteredPullers = make(map[string]Puller)

// Puller describes an interface for a Harbor Puller.
type Puller interface {
	New() Puller
	Pull(uri, id string, ops options.Options) (string, error)
}

// Wrapper is a wrapper struct for holding a
// name Puller in the registry.
type Wrapper struct {
	name   string
	Puller Puller
}

// NewPuller wraps a builder and returns a puller Wrapper.
func NewPuller(name string, puller Puller) *Wrapper {
	pullerWrap := &Wrapper{
		name:   name,
		Puller: puller,
	}

	return pullerWrap
}

// Name returns the name of the puller.
func (p *Wrapper) Name() string {
	return p.name
}

// Execute runs the pull operation for a deployment chain.
func (p *Wrapper) Execute(d *chain.Deployment, ops options.Options) error {
	workDir, err := p.Puller.Pull(d.URI, d.ID, ops)
	if err != nil {
		return err
	}

	d.SetWorkDir(workDir)

	return nil
}

// Rollback does nothing at the moment in a puller.
func (p *Wrapper) Rollback(d *chain.Deployment, ops options.Options) error {
	return nil
}

// RegisterPuller registers a puller with `name`.
func RegisterPuller(name string, puller Puller) {
	RegisteredPullers[name] = puller
}

// GetPuller returns a puller registered as `name`, if it exists.
// It returns an error if it does not exist.
func GetPuller(name string) (Puller, error) {
	puller, ok := RegisteredPullers[name]
	if !ok {
		return nil, errors.New("puller not found")
	}

	return puller, nil
}
