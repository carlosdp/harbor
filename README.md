# SupplyChain
SupplyChain is a light-weight automated, declarative-configuration, service orchestration system that is designed to assist in the continuous deployment of Docker container-based clusters in a minimal configuration environment. This design encourages preventing vendor lock-in and spreading a cluster over multiple cloud service providers to reduce dependency.

It allows you to define "Deployment Chains" in a declarative manner, like this:

```json
{"web-chain": [
  {"hook": "github-deployment", "endpoint": "web"},
  {"puller": "git-puller", "options": {"allowed_host": "github.com"}},
  {"builder": "docker-builder"},
  {"scheduler": "docker-scheduler", "options": {"port": "3000-3999"}},
  {"notifier": "consul", "options": {"service": "web"}}
]}
```

Using this configuration, SupplyChain will listen for a Deployment webhook request from GitHub, and then work its way down the chain:

- Pull the repository and commit sent by the webhook.
- Build a Docker image from the Dockerfile in the repository.
- Schedule a Docker container using the local Docker daemon.
- Notify Consul of the new container node.

It will also detect any existing previous deploys and roll it back by:

- Notifying Consul to remove the container node.
- Rolling back the scheduled container by destroying it.

Here, we used GitHub, Git, Docker, and Consul "chain links", but SupplyChain is entirely plugin based, so virtually any framework or service can be configured as long as a plugin is written and loaded.

## What it is not
SupplyChain is **not** a health-check framework or a process/node monitoring service. It strictly deals with the actions that are involved in managing the deployment and placement of service nodes, and the rolling back of actions performed in a deployment.

For example, SupplyChain will not detect when a Docker container fails, but a service that does detect that can communicate to SupplyChain and have it perform an action, such as deploying another node or rolling back the latest deployment. SupplyChain is designed to do **one** thing very well, and that is making the deployment of services over a cluster (1 node or 1000+) easy, trustworthy, stable, and reversible in a fully automated fashion.

## Design
There are 5 different types of "Chain Links" (providers) in a SupplyChain deployment chain:

### Hooks
Hooks receive notifications from external services and provide information to the rest of the chain. They are also **always** the first Link in a chain. Commonly, hooks are receive notifications of a requested deployment. An example of a Hook is a GitHub Deployment API webhook. This Hook would receive a `create_deployment` event from GitHub, authenticate it, and then gather the information on the repository in question to pass to the next stage in the chain.

### Pullers
Pullers are simply responsible for grabbing the correct version of the repository or artifact being deployed. An obvious Puller is the Git Puller which pulls from a specified Git remote repository and branch. You could also have a Puller that pulls a pre-built RPM package and skip adding a builder step.

### Builders
Builders are responsible for taking source code and building it into a package that can be passed to the Schedulers. An example of a builder would be the Docker Builder that builds a container using a Dockerfile in the source code and passes the image to the Scheduler. Builders are not always necessary in a chain as it is generally better practice to perform a build using a CI or dedicated build system which would then pass the already built package or container image to SupplyChain via a Hook.

### Schedulers
Schedulers deploy the new artifact to the cluster. An example Scheduler is the Fleet Scheduler which uses CoreOS Fleet to deploy containers across a range of cluster machines.

### Notifiers
Notifiers notify some service on the status of a deployment. One example of a Notifier is a GitHub Deployment Status API notifier. Another example is a Consul notifier which changes the service configuration on deployment.

## Automatic Rollback
If a Chain Link runs into an error, it cancels the rest of the deployment chain and executes a `Rollback` event on the already-executed Chain Links.

## Configuration
SupplyChain only requires one simple JSON file to setup a deployment chain:

*Note:* Deployment Chains must be suffixed with '-chain'

```json
{"web-chain": [
  {"hook": "github-deployment", "endpoint": "web"},
  {"puller": "git-puller", "options": {"allowed_host": "github.com"}},
  {"builder": "docker-builder"},
  {"scheduler": "fleet", "cache_deploys": 1, "options":{
    "strategy": "full_replace"
  }},
  {"notifier": "github-deployment-status", "options": {
    "api_key": "xxx",
    "api_secret": "xxx"
  }}
]}
```

The providers will be inserted into the chain in the order they are given, with the only requirement being the first provider in the chain must be a `Hook`. You can have more than one `Hook` in the chain, however, to allow for more complicated deployment processes that require waiting for some external event.

This configuration allows for easily composing multi-step, automated deployment processes. For example, `FleetScheduler` has a deployment strategy called "canary" which deploys only one instance of a container. `ConsulNotifier` allows us to tag an instance, perhaps to configure it to only receive 10% traffic through whatever means of load balancer configuration the user provides (for example, a `confd` setup).

We can use this to create a deployment chain which deploys a canary, waits 10 minutes to receive a message on a `Hook` that indicates an error of some kind (such as an event from a logging service like Airbrake). If it does receive a message, the Hook causes a `Rollback` event which the `FleetScheduler` and `ConsulNotifier` catches and takes the canary out of the load balancer (from Consul) and out of the cluster (via Fleet) and cancels the commit. Finally, we have a final Notifier which triggers on any status that lets our GitHub Deployment Status API know the build failed.

If it doesn't receive an event in the timeout, `FleetScheduler` and `ConsulNotifier` execute a quiet `Rollback` and the chain continues with a full deployment with the subsequent `FleetScheduler` and `ConsulNotifier`, finally ending with the GitHub Deployment Status API notifier.

**Note:** Another cool example is you could have an `AWSScheduler` which checks if it is necessary to spin up a new instance before continuing with a deployment. The possibilities are endless with this light-weight framework.

```json
{"web-chain":[
  {"hook": "git-deployment", "endpoint": "/web"},
  {"puller": "git-puller", "options": {"allowed_host": "github.com"}},
  {"builder": "docker-builder"},
  {"scheduler": "fleet", "options": {"strategy": "canary"},
    "always_rollback": true, "chain":
    [
      {"notifier": "consul", "always_rollback": true, "options": {
        "service": "web",
        "tags": ["canary"],
        "chain": [
          {"hook": "airbrake", "timeout": 600, "rollback_deployment": true}
        ]
      }}
    ]
  },
  {"scheduler": "fleet", "strategy": "full_replace"},
  {"notifier": "github-deployment-status", "options": {
    "api_key": "xxx",
    "api_secret": "xxx"
  }}
]}
```

### Sub Chains
Here, we used a "sub chain" by adding the "chain" attribute to the Scheduler. Any rollbacks that are triggered within Sub Chains will stop after rolling back the chain link that defined it. The rest of the chain will then continue as normal. The only exception is if a chain link has "rollback_deployment" set to `true`. In this case, a rollback triggered by this chain link will always cause a full rollback of the entire deployment chain.

Sub Chains also allow for temporary chains, as in this case. When "always_rollback" is set to `true` on a chain link, it will "soft rollback" the chain link after complete execution, causing a rollback of that chain-link, but not the rest of the chain. If that chain link also has a Sub Chain, it will wait for that Sub Chain to complete before executing that soft rollback.

### Automatic Rollbacks of Previous Deploys
SupplyChain keeps track of deployments it executes and rolls-back old deployments upon successful new deployments, automatically. Sometimes, however, you want to keep the last version of a node up (but out of the load balancer or inactive) for speedy rollbacks. You can tell SupplyChain to stop at a certain point of the chain during a rollback, and keep X previous deployments by specifying `keep: X` on a chain link. For example, I could keep the last deploy as an active Docker container, but take it out of the load balancer, using this chain:

```json
{"web-chain": [
  {"hook": "github-deployment", "endpoint": "web"},
  {"puller": "git-puller", "options": {"allowed_host": "github.com"}},
  {"builder": "docker-builder"},
  {"scheduler": "docker-scheduler", "keep": 1, "options": {"port": "3000-3999"}},
  {"notifier": "consul", "options": {"service": "web"}}
]}
```

During a new deployment, SupplyChain will run the chain, grab the last deployment, rollback the ConsulNotifier for that deployment (taking it out of Consul's service discovery), and stop at the DockerScheduler. Next time I deploy a chain, this old deploy will be run through the DockerScheduler rollback, and the rest of the chain.

### Variables
In chain definitions, you have access to two kinds of variables:

- Environment variables, prepended with `$ENV_`
- Chain Link variables defined by previous links in the chain, prepended with `$<chain link name>_`. For example: `$DOCKER_SCHEDULER_NUM_HOSTS`.

## SupplyChain CLI
SupplyChain ships with a special hook called `chain-cli-hook`. It allows you to make chains that are initiated directly via the command line using the `supply-chain` executable.

```json
{"web-chain":[
  {"hook": "chain-cli-hook", "command": "cluster-tests"},
  {"builder": "jenkins-builder", "options": {
    "repo": "git@github.com:$SUPPLY_CHAIN_CLI_REPO",
    "branch": "$SUPPLY_CHAIN_CLI_BRANCH"
  }}
]}
```

The SupplyChainCLIHook populates its variables with the flag arguments it receives on the command line. So here, we can run:

```bash
> supply-chain cluster-tests -repo carlosdp/test-app -branch new-feature
```

SupplyChain will put the argument to `-repo` in `$SUPPLY_CHAIN_CLI_REPO` and `-branch` in `$SUPPLY_CHAIN_CLI_BRANCH`.

# Framework Design

## Plugins
### Options
Every Chain Link type can have a map of "options" associated with it at configuration. These will be passed in with every call to the chain link as a `options.Options` parameter.

The `Options` data structure allows for easy, type-safe access to the options passed into the JSON configuration. For example, the FleetScheduler has a "strategy" option that defines the method in which to deploy nodes. It accesses this parameter like this:

*Note:* Here, FleetScheduler has named the passed in `options.Options` parameter in the `Execute` function "params".

```go
strategy := params.GetString("strategy")
```

`GetString` will either return the passed in string, or it will return an empty string in the event no parameter was passed in or it was of the wrong data type. If you want to explicitly check if the parameter was set, you would use the `String` function:

```go
strategy, ok := params.String("strategy")
```

In addition to strings, options can be integers:

```go
strategy := params.GetInt("strategy")
```

or booleans:

```go
strategy := params.GetBool("use_strategy")
```

They could be an array of options:

```go
strategies := params.GetArray("strategy")

for _, strategyOpt := range strategies {
  strategy := strategyOpt.GetString()
  ...
}
```

or a map of strings to options:

```go
strategyMap := params.GetMap("strategy")

for serviceName, strategyOpt := range strategyMap {
  strategy := strategyOpt.GetString()
  ...
}
```

There are also some helpful pseudo-types we created that are useful for link configurations. For example, the `IntRange`:

```json
...
{"scheduler": "docker-scheduler", "options": {"port": "3000-3999"}},
...
```

```go
min, max, ok := params.IntRange("port")
if !ok {
  min = max = params.GetInt("port")
}
```

### State
While SupplyChain enforces compatibility between all links by rigidly defining what data each link has access to (name, image name, deployment ID), some links create more than one resource, which makes tracking the resources it created complicated using just conventions derived from the provided data. For example, a FleetScheduler may be configured to create multiple instances for each deployment, maybe across multiple cloud providers, regions, and hosts. When we execute an automatic rollback when launching a new deployment, the link needs to know where all these resources it created lie.

That information can be encoded in the return value of `Schedule` and `Notify` for schedulers and notifiers (hooks, builders, and pullers do nothing on rollbacks):

```go
func (m myScheduler) Schedule(image, name, id string, ops options.Options) (interface{}, error) {
  ...
  instances := []string {"i23433-3", "ie4t5g55g-2", "ir39r944"}
  return instances, nil
}
```

This return value of `interface{}` type is known as the link's "state". It is a generic interface to allow plugin writers to pass any type of information that makes sense, from a string with a container ID, to a slice of instance IDs (as we show here), to a map that contains hosts mapping to the containers on each host. This state is passed to the `Rollback` function as a `option.Option`:

```go
func (m myScheduler) Rollback(name, id string, ops options.Options, state options.Option) error {
  instances := state.GetArray()

  for _, instanceOpt := range instances {
    instanceID := instanceOpt.GetString()
    ...
  }

  return nil
}
```

State is persisted by SupplyChain, so even in the event SupplyChain is rebooted, it will remember the state of each link.

**Note:** State cannot be modified during a rollback. That means that if a rollback partially fails, the next time `Rollback` is run, it needs to account for the fact that some resources in the state my have already been rolled back. In other words, `Rollback` needs to be an idempotent operation.
