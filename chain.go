package main

type LinkType interface {
	Name() string
	Execute() error
	Rollback() error
}

type ChainLink struct {
  Link LinkType
  SubChain *Chain
}

type Chain struct {
  Name string
	Links       []*ChainLink
	CurrentStep int
}

func NewChain() *Chain {
  chain := &Chain{}

  return chain
}
