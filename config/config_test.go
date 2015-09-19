package config_test

import (
	"strings"
	"testing"

	"github.com/carlosdp/harbor/config"
	_ "github.com/carlosdp/harbor/mocks"
)

func TestFailsOnEmptyConfig(t *testing.T) {
	c := strings.NewReader("")

	_, err := config.ParseConfig(c)
	if err == nil {
		t.Fatal("should have returned error")
	}
}

func TestFailsWhenTopLevelIsArray(t *testing.T) {
	c := strings.NewReader(`
    [
      {"hook": "fake-hook", "endpoint": "/hook"},
      {"puller": "fake-puller"},
      {"builder": "fake-builder"},
      {"scheduler": "fake-scheduler"}
    ]
  `)

	_, err := config.ParseConfig(c)
	if err == nil {
		t.Fatal("should have returned error")
	}
}

func TestReadsBasicChain(t *testing.T) {
	c := strings.NewReader(`
    {"web-chain": [
      {"hook": "fake-hook", "endpoint": "/hook"},
      {"puller": "fake-puller"},
      {"builder": "fake-builder"},
      {"scheduler": "fake-scheduler"}
    ]}
  `)

	chains, err := config.ParseConfig(c)
	if err != nil {
		t.Fatal(err)
	}

	if len(chains) != 1 {
		t.Fatal("should have parsed 1 chain, parsed: ", len(chains))
	}

	chain := chains[0]
	if chain.Name != "web-chain" {
		t.Fatal("could not find parsed chain")
	}

	if len(chain.Links) != 4 {
		t.Fatal("should have parsed 4 chain links, parsed: ", len(chains))
	}
}