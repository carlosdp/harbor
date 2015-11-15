package main

import (
	_ "github.com/carlosdp/supply-chain/plugins/builders/docker-builder"
	_ "github.com/carlosdp/supply-chain/plugins/builders/docker-pusher"
	_ "github.com/carlosdp/supply-chain/plugins/hooks/github-hook"
	_ "github.com/carlosdp/supply-chain/plugins/pullers/git-puller"
	_ "github.com/carlosdp/supply-chain/plugins/schedulers/docker-scheduler"
	_ "github.com/carlosdp/supply-chain/plugins/schedulers/nomad-scheduler"
)
