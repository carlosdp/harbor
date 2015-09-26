package chain

import (
	"errors"

	"github.com/carlosdp/harbor/options"
)

// LinkType is an is an enum that identifies a link type.
type LinkType int

// LinkTypes for each type of chain link.
const (
	HOOK LinkType = iota
	PULLER
	BUILDER
	SCHEDULER
	NOTIFIER
)

// Deployment describes a single deployment that can be sent through a chain.
type Deployment interface {
	URI() string
	Name() string
	ID() string
	WorkDir() string
	Image() string
	SetURI(uri string)
	SetName(name string)
	SetID(id string)
	SetWorkDir(workDir string)
	SetImage(image string)
	SetState(linkName string, state interface{})
	GetState(linkName string) options.Option
}

type link interface {
	Name() string
	Execute(d Deployment, ops options.Options) error
	Rollback(d Deployment, ops options.Options) error
}

// Link is a link in the chain.
type Link struct {
	Link       link
	Type       LinkType
	Parameters options.Options
	Options    options.Options
	SubChain   *Chain
}

// Chain is a deployment chain.
type Chain struct {
	Name              string
	Links             []*Link
	ActiveDeployments []Deployment
}

// NewChain returns an empty deployment chain.
func NewChain() *Chain {
	return &Chain{}
}

// NewLink returns an empty chain link.
func NewLink() *Link {
	return &Link{}
}

// LinksOfType returns a slice of links in the chain of `t` type.
func (c *Chain) LinksOfType(t LinkType) []*Link {
	var links []*Link

	for _, link := range c.Links {
		if link.Type == t {
			links = append(links, link)
		}
	}

	return links
}

// LinkPosition returns the zero-indexed position of `link` in the chain.
func (c *Chain) LinkPosition(link *Link) (int, error) {
	for i, l := range c.Links {
		if l == link {
			return i, nil
		}
	}

	return -1, errors.New("link not in chain")
}
