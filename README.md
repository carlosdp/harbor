# Harbor
Harbor is a light-weight automated service orchestration system that is designed to assist in the continuous deployment of Docker container-based clusters in a minimal configuration environment. This design encourages preventing vendor lock-in and spreading a cluster over multiple cloud service providers to reduce dependency.

## Design
There are 4 different types of providers in a Harbor deployment chain:

### Hooks
Hooks are the providers that receive notification of a requested deployment and gather any necessary information. An example of a Hook is a GitHub Deployment API webhook. This Hook would receive a `create_deployment` event from GitHub, authenticate it, and then gather the information on the repository in question to pass to the next stage in the chain.

### Pullers
Pullers are simply responsible for grabbing the correct version of the repository being deployed. An obvious Puller is the Git Puller which pulls from a specified Git remote repository and branch.

### Schedulers
Schedulers deploy the new container to the cluster. An example Scheduler is the Fleet Scheduler which uses CoreOS Fleet to deploy containers across a range of cluster machines.

### Notifiers
Notifiers notify some service on the status of a deployment. One example of a Notifier is a GitHub Deployment Status API notifier. Another example is a Consul notifier which changes the service configuration on deployment.

If a provider runs into an error, it cancels the rest of the deployment chain and executes a `Rollback` event on the chain. More on this later.

## Configuration
Harbor only requires one simple json file to setup a deployment chain:

*Note:* Deployment Chains must be suffixed with '-chain'

```json
{"web-chain": [ 
  {"hook": "git_deployment", "endpoint": "web"},
  {"puller": "git-puller", "allowed_host": "github.com"},
  {"scheduler": "fleet", "strategy": "full_replace"},
  {"notifier": "github_deployment_status", "api_key": "xxx"
    "api_secret": "xxx"}
]}
```

The providers will be inserted into the chain in the order they are given, with the only requirement being the first provider in the chain must be a `Hook`. You can have more than one `Hook` in the chain, however, to allow for more complicated deployment processes that require waiting for some external event.

This configuration allows for easily composing multi-step, automated deployment processes. For example, `FleetScheduler` has a deployment strategy called "canary" which deploys only one instance of a container. `ConsulNotifier` allows us to tag an instance, perhaps to configure it to only receive 10% traffic through whatever means of load balancer configuration the user provides (for example, a `confd` setup).

We can use this to create a deployment chain which deploys a canary, waits 10 minutes to receive a message on a `Hook` that indicates an error of some kind (such as an event from a logging service like Airbrake). If it does receive a message, the Hook causes a `Rollback` event which the `FleetScheduler` and `ConsulNotifier` catches and takes the canary out of the load balancer (from Consul) and out of the cluster (via Fleet) and cancels the commit. Finally, we have a final Notifier which triggers on any status that lets our GitHub Deployment Status API know the build failed.

If it doesn't receive an event in the timeout, `FleetScheduler` and `ConsulNotifier` execute a quiet `Rollback` and the chain continues with a full deployment with the subsequent `Fleet-Scheduler` and `Consul-Notifier`, finally ending with the GitHub Deployment Status API notifier.

```json
{"web-chain": [
  {"hook": "git_deployment", "endpoint": "web"},
  {"puller": "git-puller", "allowed_host": "github.com"},
  {"scheduler": "fleet", "strategy": "canary",
    "always_rollback": true, "then":
    [
      {"notifier": "consul", "service": "web",
        "tags": ["canary"], "always_rollback": true,
        "then": [
          {"hook": "airbrake", "timeout": 600, "execute_rollback" true}
        ]
      }
    ]
  },
  {"scheduler": "fleet", "strategy": "full_replace"},
  {"notifier": "github_deployment_status", "api_key": "xxx"
    "api_secret": "xxx"}
]}
```

Another cool example is you could have an `AWSScheduler` which checks if it is neccessary to spin up a new instance before continuing with a deployment. The possibilities are endless with this light-weight framework.
