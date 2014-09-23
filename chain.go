package main

import (
	"errors"
)

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

func (c *Chain) LinkPosition(link *ChainLink) (int, error) {
	for i, l := range c.Links {
		if l == link {
			return i, nil
		}
	}

	return -1, errors.New("link not in chain")
}
