package puller

import (
	"errors"
	"github.com/carlosdp/harbor/chain"
)

var RegisteredPullers = make(map[string]Puller)

type Puller interface {
	New() Puller
	Pull(uri string) (string, error)
}

type PullerWrapper struct {
	name   string
	Puller Puller
}

func NewPuller(name string, puller Puller) *PullerWrapper {
	pullerWrap := &PullerWrapper{
		name:   name,
		Puller: puller,
	}

	return pullerWrap
}

func (p *PullerWrapper) Name() string {
	return p.name
}

func (p *PullerWrapper) Execute(d chain.Deployment) error {
	workDir, err := p.Puller.Pull(d.URI())
	if err != nil {
		return err
	}

	d.SetWorkDir(workDir)

	return nil
}

func (p *PullerWrapper) Rollback() error {
	return nil
}

func RegisterPuller(name string, puller Puller) {
	RegisteredPullers[name] = puller
}

func GetPuller(name string) (Puller, error) {
	puller, ok := RegisteredPullers[name]
	if !ok {
		return nil, errors.New("puller not found")
	}

	return puller, nil
}
