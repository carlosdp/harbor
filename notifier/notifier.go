package notifier

import (
	"errors"

	"github.com/carlosdp/harbor/chain"
	"github.com/carlosdp/harbor/options"
)

// RegisteredNotifiers contains the notifiers registered with the Harbor build.
var RegisteredNotifiers = make(map[string]Notifier)

// Notifier describes an interface for a Harbor Notifier.
type Notifier interface {
	New() Notifier
	Notify(name, id string, ops options.Options) (interface{}, error)
	Rollback(name, id string, ops options.Options, state options.Option) error
}

// Wrapper is a wrapper struct for holding a
// named Notifier in the registry.
type Wrapper struct {
	name     string
	Notifier Notifier
}

// NewNotifier wraps a notifier and returns a notifier Wrapper.
func NewNotifier(name string, notifier Notifier) *Wrapper {
	notifierWrap := &Wrapper{
		name:     name,
		Notifier: notifier,
	}

	return notifierWrap
}

// Name returns the name of the notifier.
func (s *Wrapper) Name() string {
	return s.name
}

// Execute runs the schedule operation for a deployment chain.
func (s *Wrapper) Execute(d *chain.Deployment, ops options.Options) error {
	state, err := s.Notifier.Notify(d.Name, d.ID, ops)
	if err != nil {
		return err
	}

	d.SetState(s.Name(), state)

	return nil
}

// Rollback reverts the scheduling operation for a deployment chain.
func (s *Wrapper) Rollback(d *chain.Deployment, ops options.Options) error {
	state := d.GetState(s.Name())
	return s.Notifier.Rollback(d.Name, d.ID, ops, state)
}

// RegisterNotifier registers a notifier with `name`.
func RegisterNotifier(name string, notifier Notifier) {
	RegisteredNotifiers[name] = notifier
}

// GetNotifier returns a notifier registered with `name`, if it exists.
// It returns an error if it does not exist.
func GetNotifier(name string) (Notifier, error) {
	notifier, ok := RegisteredNotifiers[name]
	if !ok {
		return nil, errors.New("notifier not found")
	}

	return notifier, nil
}
