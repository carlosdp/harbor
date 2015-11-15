package config_test

import (
	"strings"
	"testing"

	"github.com/carlosdp/supply-chain/config"
	_ "github.com/carlosdp/supply-chain/mocks"
)

func TestFailsOnEmptyConfig(t *testing.T) {
	t.Parallel()
	c := strings.NewReader("")

	_, err := config.ParseConfig(c)
	if err == nil {
		t.Fatal("should have returned error")
	}
}

func TestFailsWhenTopLevelIsArray(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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

func TestFailsUnlessStartsWithHook(t *testing.T) {
	t.Parallel()
	c := strings.NewReader(`
    {"web-chain": [
      {"puller": "fake-puller"},
      {"hook": "fake-hook", "endpoint": "/hook"},
      {"builder": "fake-builder"},
      {"scheduler": "fake-scheduler"}
    ]}
	`)

	_, err := config.ParseConfig(c)
	if err == nil {
		t.Fatal("should have returned error")
	}
}

func TestParsesOptions(t *testing.T) {
	t.Parallel()
	c := strings.NewReader(`
		{"web-chain": [
			{"hook": "fake-hook", "endpoint": "fake", "options": {"test": "success"}}
		]}
	`)

	chains, err := config.ParseConfig(c)
	if err != nil {
		t.Fatal(err)
	}

	link := chains[0].Links[0]
	if link.Options.GetString("test") != "success" {
		t.Fatal("option was not parsed correctly")
	}
}

func TestIgnoresNonMapOptions(t *testing.T) {
	t.Parallel()
	c := strings.NewReader(`
		{"web-chain": [
			{"hook": "fake-hook", "endpoint": "fake", "options": "invalid"}
		]}
	`)

	chains, err := config.ParseConfig(c)
	if err != nil {
		t.Fatal(err)
	}

	link := chains[0].Links[0]
	if link.Parameters.GetString("options") != "" {
		t.Fatal("should have ignored invalid options parameter")
	}
}

func TestParsesParameters(t *testing.T) {
	t.Parallel()
	c := strings.NewReader(`
		{"web-chain": [
			{"hook": "fake-hook", "endpoint": "fake", "test": "success"}
		]}
	`)

	chains, err := config.ParseConfig(c)
	if err != nil {
		t.Fatal(err)
	}

	link := chains[0].Links[0]
	if link.Parameters.GetString("test") != "success" {
		t.Fatal("paramter was not parsed correctly")
	}
}
