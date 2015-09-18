package scheduler

import (
	"errors"

	"github.com/carlosdp/harbor/chain"
)

// RegisteredSchedulers contains the schedulers registered with the Harbor build.
var RegisteredSchedulers = make(map[string]Scheduler)

// Scheduler describes an interface for a Harbor Scheduler.
type Scheduler interface {
	New() Scheduler
	Schedule(image, name, id string) error
	Rollback(name, id string) error
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
func (s *Wrapper) Execute(d chain.Deployment) error {
	err := s.Scheduler.Schedule(d.Image(), d.Name(), d.ID())
	return err
}

// Rollback reverts the scheduling operation for a deployment chain.
// TODO: Implement this
func (s *Wrapper) Rollback() error {
	return nil
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
