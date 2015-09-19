package config

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/carlosdp/harbor/builder"
	"github.com/carlosdp/harbor/chain"
	"github.com/carlosdp/harbor/hook"
	"github.com/carlosdp/harbor/puller"
	"github.com/carlosdp/harbor/scheduler"
)

type linkConfig struct {
	// Name of the Chain Link
	Name string
	// Type of Chain Link
	Type chain.LinkType
	// Top-level Link parameters
	Parameters map[string]interface{}
	// Link-specific options
	Options map[string]interface{}
	// The Sub Chain for this Link
	SubChain []linkConfig
}

// ParseConfig parses a JSON chain config file and returns a complete chain.
func ParseConfig(configFile io.Reader) ([]*chain.Chain, error) {
	var chainDefs map[string]interface{}
	decoder := json.NewDecoder(configFile)

	err := decoder.Decode(&chainDefs)
	if err != nil {
		return nil, err
	}

	var chains []*chain.Chain

	for chainName, chainDef := range chainDefs {
		chainLinkDefs, ok := chainDef.([]interface{})
		if !ok {
			return nil, errors.New("chain is not an array")
		}

		chain, err := parseChain(chainLinkDefs)
		if err != nil {
			return nil, err
		}

		chain.Name = chainName

		chains = append(chains, chain)
	}

	return chains, nil
}

func parseChain(linkDefs []interface{}) (*chain.Chain, error) {
	newChain := chain.NewChain()
	newChain.Links = make([]*chain.Link, 0)
	for _, linkDef := range linkDefs {
		linkMap, ok := linkDef.(map[string]interface{})
		if !ok {
			return nil, errors.New("link def is not map")
		}

		link := chain.NewLink()
		if _, ok = linkMap["hook"]; ok {
			name, ok := linkMap["hook"].(string)
			if !ok {
				return nil, errors.New("hook does not have name")
			}

			hookInt, err := hook.GetHook(name)
			if err != nil {
				return nil, err
			}

			endpoint, ok := linkMap["endpoint"]
			if !ok {
				return nil, errors.New("hook needs endpoint")
			}

			endpointStr, ok := endpoint.(string)
			if !ok {
				return nil, errors.New("endpoint needs to be string")
			}

			hookWrap := hook.NewHook(name, hookInt, endpointStr)
			link.Link = hookWrap
			link.Type = chain.HOOK
		} else if _, ok = linkMap["puller"]; ok {
			name, ok := linkMap["puller"].(string)
			if !ok {
				return nil, errors.New("puller does not have name")
			}

			pullerInt, err := puller.GetPuller(name)
			if err != nil {
				return nil, err
			}

			pullerWrap := puller.NewPuller(name, pullerInt)
			link.Link = pullerWrap
			link.Type = chain.PULLER
		} else if _, ok = linkMap["builder"]; ok {
			name, ok := linkMap["builder"].(string)
			if !ok {
				return nil, errors.New("builder does not have name")
			}

			builderInt, err := builder.GetBuilder(name)
			if err != nil {
				return nil, err
			}

			builderWrap := builder.NewBuilder(name, builderInt)
			link.Link = builderWrap
			link.Type = chain.BUILDER
		} else if _, ok = linkMap["scheduler"]; ok {
			name, ok := linkMap["scheduler"].(string)
			if !ok {
				return nil, errors.New("scheduler does not have name")
			}

			schedulerInt, err := scheduler.GetScheduler(name)
			if err != nil {
				return nil, err
			}

			schedulerWrap := scheduler.NewScheduler(name, schedulerInt)
			link.Link = schedulerWrap
			link.Type = chain.SCHEDULER
		} else if _, ok = linkMap["notifier"]; ok {
		} else {
			return nil, errors.New("link type not recognized")
		}

		newChain.Links = append(newChain.Links, link)
	}

	return newChain, nil
}
