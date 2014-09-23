package main

type ChainLinkType int

const (
	HOOK ChainLinkType = iota
	PULLER
	SCHEDULER
	NOTIFIER
)

type LinkType interface {
	Name() string
	Execute() error
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

func (c *Chain) LinksOfType(t ChainLinkType) []LinkType {
	links := make([]LinkType, 0)

	for _, link := range c.Links {
		if link.Type == t {
			links = append(links, link.Link)
		}
	}

	return links
}
