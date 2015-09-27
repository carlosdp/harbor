package builder

import (
	"errors"

	"github.com/carlosdp/harbor/chain"
	"github.com/carlosdp/harbor/options"
)

// RegisteredBuilders contains the builders registered with the Harbor build.
var RegisteredBuilders = make(map[string]Builder)

// Builder describes an interface for a Harbor Builder.
type Builder interface {
	// Initializes a Builder
	New() Builder
	// Performs a build in the `workDir` source directory
	// and outputing `image` image.
	Build(workDir, image string, ops options.Options) (string, error)
}

// Wrapper is a wrapper struct for holding a
// named Builder in the registry.
type Wrapper struct {
	name    string
	Builder Builder
}

// NewBuilder wraps a builder and returns a builder Wrapper.
func NewBuilder(name string, builder Builder) *Wrapper {
	builderWrap := &Wrapper{
		name:    name,
		Builder: builder,
	}

	return builderWrap
}

// Name returns the name of the builder.
func (b *Wrapper) Name() string {
	return b.name
}

// Execute runs the build operation for a deployment chain.
func (b *Wrapper) Execute(d *chain.Deployment, ops options.Options) error {
	newImage, err := b.Builder.Build(d.WorkDir, d.Image, ops)
	if err == nil {
		d.SetImage(newImage)
	}
	return err
}

// Rollback does nothing at the moment in a builder.
func (b *Wrapper) Rollback(d *chain.Deployment, ops options.Options) error {
	return nil
}

// RegisterBuilder registers a builder with `name`.
func RegisterBuilder(name string, builder Builder) {
	RegisteredBuilders[name] = builder
}

// GetBuilder returns a builder registered as `name`, if it exists.
// It returns an error if it does not exist.
func GetBuilder(name string) (Builder, error) {
	builder, ok := RegisteredBuilders[name]
	if !ok {
		return nil, errors.New("builder not found")
	}

	return builder, nil
}
