package hook

import (
	"errors"
	"net/http"
)

var RegisteredHooks = make(map[string]Hook)

type Hook interface {
	New() Hook
	Name() string
	HandleRequest(req *http.Request)
	DeploymentID() string
	URI() string
}

type HookWrapper struct {
	name     string
	Endpoint string
	Hook     Hook
}

func NewHook(name string, hook Hook, endpoint string) *HookWrapper {
	hookWrap := &HookWrapper{
		name:     name,
		Hook:     hook.New(),
		Endpoint: endpoint,
	}

	return hookWrap
}

func (hw *HookWrapper) Name() string {
	return hw.name
}

func (hw *HookWrapper) Execute() error {
	return nil
}

func (hw *HookWrapper) Rollback() error {
	return nil
}

func RegisterHook(name string, hook Hook) {
	RegisteredHooks[name] = hook
}

func GetHook(name string) (Hook, error) {
	hook, ok := RegisteredHooks[name]
	if !ok {
		return nil, errors.New("hook not found")
	}

	return hook, nil
}
