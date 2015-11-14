package main

import (
	_ "github.com/carlosdp/harbor/plugins/builders/docker-builder"
	_ "github.com/carlosdp/harbor/plugins/builders/docker-pusher"
	_ "github.com/carlosdp/harbor/plugins/hooks/github-hook"
	_ "github.com/carlosdp/harbor/plugins/pullers/git-puller"
	_ "github.com/carlosdp/harbor/plugins/schedulers/docker-scheduler"
	_ "github.com/carlosdp/harbor/plugins/schedulers/nomad-scheduler"
)
