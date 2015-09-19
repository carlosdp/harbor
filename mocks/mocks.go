package config_test

import (
	"net/http"

	"github.com/carlosdp/harbor/builder"
	"github.com/carlosdp/harbor/hook"
	"github.com/carlosdp/harbor/puller"
	"github.com/carlosdp/harbor/scheduler"
)

type fakeHook struct {
}

type fakePuller struct {
}

type fakeBuilder struct {
}

type fakeScheduler struct {
}

func init() {
	hook.RegisterHook("fake-hook", &fakeHook{})
	puller.RegisterPuller("fake-puller", &fakePuller{})
	builder.RegisterBuilder("fake-builder", &fakeBuilder{})
	scheduler.RegisterScheduler("fake-scheduler", &fakeScheduler{})
}

func (h *fakeHook) New() hook.Hook {
	return &fakeHook{}
}

func (h *fakeHook) HandleRequest(req *http.Request) error {
	return nil
}

func (h *fakeHook) Name() string {
	return ""
}

func (h *fakeHook) DeploymentID() string {
	return ""
}

func (h *fakeHook) URI() string {
	return ""
}

func (p *fakePuller) New() puller.Puller {
	return &fakePuller{}
}

func (p *fakePuller) Pull(uri, id string) (string, error) {
	return "", nil
}

func (d *fakeBuilder) New() builder.Builder {
	return &fakeBuilder{}
}

func (d *fakeBuilder) Build(workDir, image string) (string, error) {
	return "", nil
}

func (s *fakeScheduler) New() scheduler.Scheduler {
	return &fakeScheduler{}
}

func (s *fakeScheduler) Schedule(image, name, id string) error {
	return nil
}

func (s *fakeScheduler) Rollback(name, id string) error {
	return nil
}
