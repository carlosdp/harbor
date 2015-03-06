package scheduler

import (
	"github.com/carlosdp/harbor/chain"
)

var RegisteredSchedulers = make(map[string]Scheduler)

type Scheduler interface {
	New() Scheduler
	Schedule(image, name, id string) error
	Rollback(name, id string) error
}

type SchedulerWrapper struct {
	name      string
	Scheduler Scheduler
}

func NewScheduler(name string, scheduler Scheduler) *SchedulerWrapper {
	schedulerWrap := &SchedulerWrapper{
		name:      name,
		Scheduler: scheduler,
	}

	return schedulerWrap
}

func (s *SchedulerWrapper) Name() string {
	return s.name
}

func (s *SchedulerWrapper) Execute(d chain.Deployment) error {
	err := s.Scheduler.Schedule(d.Image(), d.Name(), d.ID())
	return err
}

func (s *SchedulerWrapper) Rollback() error {
	return nil
}
