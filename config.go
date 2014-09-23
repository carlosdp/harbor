package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/carlosdp/harbor/hook"
	"os"
)

func ParseConfig(configPath string) ([]*Chain, error) {
	fmt.Println("Parsing chain")
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	var chainDefs map[string]interface{}
	decoder := json.NewDecoder(f)

	err = decoder.Decode(&chainDefs)
	if err != nil {
		return nil, err
	}

	chains := make([]*Chain, 0)

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

func parseChain(linkDefs []interface{}) (*Chain, error) {
	chain := NewChain()
	chain.Links = make([]*ChainLink, 0)
	for _, linkDef := range linkDefs {
		linkMap, ok := linkDef.(map[string]interface{})
		if !ok {
			return nil, errors.New("link def is not map")
		}

		link := NewChainLink()
		if _, ok = linkMap["hook"]; ok {
			fmt.Println("hook detected")
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
			link.Type = HOOK
		} else if _, ok = linkMap["puller"]; ok {
			fmt.Println("puller detected")
		} else if _, ok = linkMap["scheduler"]; ok {
			schName, ok := linkMap["scheduler"].(string)

			if ok {
				fmt.Println("scheduler detected: ", schName)
			}
		} else if _, ok = linkMap["notifier"]; ok {
			fmt.Println("notifier detected")
		} else {
			return nil, errors.New("link type not recognized")
		}

		chain.Links = append(chain.Links, link)
	}

	return chain, nil
}
