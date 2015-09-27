package scheduler

import (
	"errors"

	"github.com/carlosdp/harbor/chain"
	"github.com/carlosdp/harbor/options"
)

// RegisteredSchedulers contains the schedulers registered with the Harbor build.
var RegisteredSchedulers = make(map[string]Scheduler)

// Scheduler describes an interface for a Harbor Scheduler.
type Scheduler interface {
	New() Scheduler
	Schedule(image, name, id string, ops options.Options) (interface{}, error)
	Rollback(name, id string, ops options.Options, state options.Option) error
}

// Wrapper is a wrapper struct for holding a
// named Scheduler in the registry.
type Wrapper struct {
	name      string
	Scheduler Scheduler
}

// NewScheduler wraps a scheduler and returns a scheduler Wrapper.
func NewScheduler(name string, scheduler Scheduler) *Wrapper {
	schedulerWrap := &Wrapper{
		name:      name,
		Scheduler: scheduler,
	}

	return schedulerWrap
}

// Name returns the name of the scheduler.
func (s *Wrapper) Name() string {
	return s.name
}

// Execute runs the schedule operation for a deployment chain.
func (s *Wrapper) Execute(d *chain.Deployment, ops options.Options) error {
	state, err := s.Scheduler.Schedule(d.Image, d.Name, d.ID, ops)
	if err != nil {
		return err
	}

	d.SetState(s.Name(), state)

	return nil
}

// Rollback reverts the scheduling operation for a deployment chain.
func (s *Wrapper) Rollback(d *chain.Deployment, ops options.Options) error {
	state := d.GetState(s.Name())
	return s.Scheduler.Rollback(d.Name, d.ID, ops, state)
}

// RegisterScheduler registers a scheduler with `name`.
func RegisterScheduler(name string, scheduler Scheduler) {
	RegisteredSchedulers[name] = scheduler
}

// GetScheduler returns a scheduler registered with `name`, if it exists.
// It returns an error if it does not exist.
func GetScheduler(name string) (Scheduler, error) {
	scheduler, ok := RegisteredSchedulers[name]
	if !ok {
		return nil, errors.New("scheduler not found")
	}

	return scheduler, nil
}
