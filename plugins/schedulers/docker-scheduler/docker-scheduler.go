package dockerscheduler

import (
	"os"
	"strings"

	"github.com/carlosdp/harbor/options"
	"github.com/carlosdp/harbor/scheduler"
	"github.com/fsouza/go-dockerclient"
)

func init() {
	scheduler.RegisterScheduler("docker-scheduler", &dockerScheduler{})
}

type dockerScheduler struct {
}

func (ds *dockerScheduler) New() scheduler.Scheduler {
	return &dockerScheduler{}
}

func (ds *dockerScheduler) Schedule(image, name, id string, ops options.Options) (interface{}, error) {
	name = strings.Replace(name, "/", "-", -1)
	cID, err := createContainer(image, name+"-"+id)
	return cID, err
}

func (ds *dockerScheduler) Rollback(name, id string, ops options.Options, state options.Option) error {
	cID := state.GetString()

	dockerHost := os.Getenv("DOCKER_HOST")

	var client *docker.Client
	var err error

	if dockerHost == "" {
		dockerHost = "unix:///var/run/docker.sock"
		client, err = docker.NewClient(dockerHost)
	} else {
		client, err = docker.NewClientFromEnv()
	}

	if err != nil {
		return err
	}

	client.StopContainer(cID, 10)

	err = client.RemoveContainer(docker.RemoveContainerOptions{
		ID: cID,
	})

	return err
}

func createContainer(imageName, name string) (string, error) {
	dockerHost := os.Getenv("DOCKER_HOST")

	var client *docker.Client
	var err error

	if dockerHost == "" {
		dockerHost = "unix:///var/run/docker.sock"
		client, err = docker.NewClient(dockerHost)
	} else {
		client, err = docker.NewClientFromEnv()
	}

	if err != nil {
		return "", err
	}

	c, err := client.CreateContainer(docker.CreateContainerOptions{
		Name: name,
		Config: &docker.Config{
			Image: imageName,
		},
	})

	if err != nil {
		return "", err
	}

	err = client.StartContainer(c.ID, &docker.HostConfig{})

	if err != nil {
		return "", err
	}

	return c.ID, nil
}
