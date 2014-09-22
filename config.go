package main

import (
	"encoding/json"
  "os"
  "errors"
  "fmt"
)

func ParseConfig(configPath string) ([]*Chain, error) {
  fmt.Println("Parsing chain")
  f, err := os.Open(configPath)
  if err != nil {
    return nil, err
  }

  var chains map[string]interface{}
  decoder := json.NewDecoder(f)

  err = decoder.Decode(&chains)
  if err != nil {
    return nil, err
  }

  for chainName, chainDef := range chains {
    chainLinkDefs, ok := chainDef.([]interface{})
    if !ok {
      return nil, errors.New("chain is not an array")
    }

    chain, err := parseChain(chainLinkDefs)
    if err != nil {
      return nil, err
    }

    chain.Name = chainName
  }

	return nil, nil
}

func parseChain(linkDefs []interface{}) (*Chain, error) {
  chain := NewChain()
  for _, linkDef := range linkDefs {
    linkMap, ok := linkDef.(map[string]interface{})
    if !ok {
      return nil, errors.New("link def is not map")
    }

    if _, ok = linkMap["hook"]; ok {
      fmt.Println("hook detected")
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
  }

  return chain, nil
}
