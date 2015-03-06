package hook

import (
	"errors"
	"net/http"

	"github.com/carlosdp/harbor/chain"
)

var RegisteredHooks = make(map[string]Hook)

type Hook interface {
	New() Hook
	Name() string
	HandleRequest(req *http.Request) error
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

func (hw *HookWrapper) Execute(d chain.Deployment) error {
	return nil
}

func (hw *HookWrapper) HandleRequest(d chain.Deployment, req *http.Request) error {
	err := hw.Hook.HandleRequest(req)
	if err != nil {
		return err
	}

	d.SetName(hw.Hook.Name())
	d.SetURI(hw.Hook.URI())
	d.SetID(hw.Hook.DeploymentID())
	d.SetImage(hw.Hook.Name() + "-" + hw.Hook.DeploymentID())
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
