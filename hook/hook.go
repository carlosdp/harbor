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
	Start(mux *http.ServeMux, queue chan<- Request, ops options.Options) error
	DeploymentID() string
	URI() string
}

// Wrapper is a wrapper struct for holding a
// names Hook in the registry.
type Wrapper struct {
	name string
	Hook Hook
}

// Request is a request object sent to the deployment runner by a hook.
type Request struct {
	Name         string
	DeploymentID string
	URI          string
	Image        string
	Chain        *chain.Chain
	Link         *chain.Link
}

// NewRequest creates a new hook request.
func NewRequest(name, id, uri, image string) Request {
	return Request{
		Name:         name,
		DeploymentID: id,
		URI:          uri,
		Image:        image,
	}
}

// NewHook wraps a hook and returns a hook Wrapper.
func NewHook(name string, hook Hook) *Wrapper {
	hookWrap := &Wrapper{
		name: name,
		Hook: hook.New(),
	}

	return hookWrap
}

// Name returns the name of the hook.
func (hw *Wrapper) Name() string {
	return hw.name
}

// Execute does nothing at the moment in a hook.
func (hw *Wrapper) Execute(d *chain.Deployment, ops options.Options) error {
	return nil
}

// Start starts up the hook with the given channel for passing deployment information.
func (hw *Wrapper) Start(mux *http.ServeMux, queue chan<- Request, ops options.Options, c *chain.Chain, link *chain.Link) {
	q := make(chan Request)
	go hw.Hook.Start(mux, q, ops)

	for {
		select {
		case r := <-q:
			r.Chain = c
			r.Link = link
			queue <- r
		}
	}
}

// Rollback does nothing at the moment in a hook.
func (hw *Wrapper) Rollback(d *chain.Deployment, ops options.Options) error {
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
