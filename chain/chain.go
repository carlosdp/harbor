package chain

import (
	"errors"
)

type ChainLinkType int

const (
	HOOK ChainLinkType = iota
	PULLER
	BUILDER
	SCHEDULER
	NOTIFIER
)

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
}

type LinkType interface {
	Name() string
	Execute(d Deployment) error
	Rollback() error
}

type ChainLink struct {
	Link     LinkType
	Type     ChainLinkType
	SubChain *Chain
}

type Chain struct {
	Name  string
	Links []*ChainLink
}

func NewChain() *Chain {
	return &Chain{}
}

func NewChainLink() *ChainLink {
	return &ChainLink{}
}

func (c *Chain) LinksOfType(t ChainLinkType) []*ChainLink {
	links := make([]*ChainLink, 0)

	for _, link := range c.Links {
		if link.Type == t {
			links = append(links, link)
		}
	}

	return links
}

func (c *Chain) LinkPosition(link *ChainLink) (int, error) {
	for i, l := range c.Links {
		if l == link {
			return i, nil
		}
	}

	return -1, errors.New("link not in chain")
}
