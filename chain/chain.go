package chain

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"

	log "github.com/carlosdp/harbor/Godeps/_workspace/src/github.com/Sirupsen/logrus"
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

type link interface {
	Name() string
	Execute(d *Deployment, ops options.Options) error
	Rollback(d *Deployment, ops options.Options) error
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
	ActiveDeployments []*Deployment
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

// Persist dumps the state of all active deployments to disc.
func (c *Chain) Persist() {
	output, err := json.Marshal(c.ActiveDeployments)
	if err != nil {
		log.Errorf("[%v] Could not store chain state: %v", c.Name, err)
	}

	f, err := os.Create(c.Name + ".json")
	if err != nil {
		log.Errorf("[%v] Could not store chain state: %v", c.Name, err)
	}
	defer f.Close()
	io.Copy(f, bytes.NewReader(output))
}

// Load loads the state of active deployments from disc.
func (c *Chain) Load() {
	if _, err := os.Stat(c.Name + ".json"); os.IsNotExist(err) {
		return
	}
	f, err := os.Open(c.Name + ".json")
	if err != nil {
		log.Errorf("[%v] Could not load chain state: %v", c.Name, err)
	}

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		log.Errorf("[%v] Could not load chain state: %v", c.Name, err)
	}

	json.Unmarshal(buf, &c.ActiveDeployments)

	for _, d := range c.ActiveDeployments {
		d.Chain = c
	}
}
