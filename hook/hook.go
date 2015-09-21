package hook

import (
	"errors"
	"net/http"

	"github.com/carlosdp/harbor/chain"
	"github.com/carlosdp/harbor/options"
)

// RegisteredHooks contains the hooks registered with the Harbor build.
var RegisteredHooks = make(map[string]Hook)

// Hook describes an interface for a Harbor Hook.
type Hook interface {
	New() Hook
	Name() string
	HandleRequest(req *http.Request, ops options.Options) error
	DeploymentID() string
	URI() string
}

// Wrapper is a wrapper struct for holding a
// names Hook in the registry.
type Wrapper struct {
	name     string
	Endpoint string
	Hook     Hook
}

// NewHook wraps a hook and returns a hook Wrapper.
func NewHook(name string, hook Hook, endpoint string) *Wrapper {
	hookWrap := &Wrapper{
		name:     name,
		Hook:     hook.New(),
		Endpoint: endpoint,
	}

	return hookWrap
}

// Name returns the name of the hook.
func (hw *Wrapper) Name() string {
	return hw.name
}

// Execute does nothing at the moment in a hook.
func (hw *Wrapper) Execute(d chain.Deployment, ops options.Options) error {
	return nil
}

// HandleRequest handles an incoming web hook request.
func (hw *Wrapper) HandleRequest(d chain.Deployment, req *http.Request, ops options.Options) error {
	err := hw.Hook.HandleRequest(req, ops)
	if err != nil {
		return err
	}

	d.SetName(hw.Hook.Name())
	d.SetURI(hw.Hook.URI())
	d.SetID(hw.Hook.DeploymentID())
	d.SetImage(hw.Hook.Name() + "-" + hw.Hook.DeploymentID())
	return nil
}

// Rollback does nothing at the moment in a hook.
func (hw *Wrapper) Rollback(ops options.Options) error {
	return nil
}

// RegisterHook registers a hook with `name`.
func RegisterHook(name string, hook Hook) {
	RegisteredHooks[name] = hook
}

// GetHook returns a hook registered as `name`, if it exists.
// It returns an error if it does not exist.
func GetHook(name string) (Hook, error) {
	hook, ok := RegisteredHooks[name]
	if !ok {
		return nil, errors.New("hook not found")
	}

	return hook, nil
}
