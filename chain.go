package main

type HookType int

const (
	HOOK HookType = iota
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
	Type     HookType
	SubChain *Chain
}

type Chain struct {
	Name        string
	Links       []*ChainLink
	CurrentStep int
}

func NewChain() *Chain {
	return &Chain{}
}

func NewChainLink() *ChainLink {
	return &ChainLink{}
}
