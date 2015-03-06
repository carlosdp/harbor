package builder

import (
	"errors"

	"github.com/carlosdp/harbor/chain"
)

var RegisteredBuilders = make(map[string]Builder)

type Builder interface {
	New() Builder
	Build(workDir, image string) (string, error)
}

type BuilderWrapper struct {
	name    string
	Builder Builder
}

func NewBuilder(name string, builder Builder) *BuilderWrapper {
	builderWrap := &BuilderWrapper{
		name:    name,
		Builder: builder,
	}

	return builderWrap
}

func (b *BuilderWrapper) Name() string {
	return b.name
}

func (b *BuilderWrapper) Execute(d chain.Deployment) error {
	newImage, err := b.Builder.Build(d.WorkDir(), d.Image())
	if err == nil {
		d.SetImage(newImage)
	}
	return err
}

func (b *BuilderWrapper) Rollback() error {
	return nil
}

func RegisterBuilder(name string, builder Builder) {
	RegisteredBuilders[name] = builder
}

func GetBuilder(name string) (Builder, error) {
	builder, ok := RegisteredBuilders[name]
	if !ok {
		return nil, errors.New("builder not found")
	}

	return builder, nil
}
