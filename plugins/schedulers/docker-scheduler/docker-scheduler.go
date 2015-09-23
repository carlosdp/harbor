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

func (ds *dockerScheduler) Schedule(image, name, id string, ops options.Options) error {
	name = strings.Replace(name, "/", "-", -1)
	err := createContainer(image, name+"-"+id)
	return err
}

func (ds *dockerScheduler) Rollback(name, id string, ops options.Options) error {
	return nil
}

func createContainer(imageName, name string) error {
	dockerHost := os.Getenv("DOCKER_HOST")

	if dockerHost == "" {
		dockerHost = "unix:///var/run/docker.sock"
	}

	client, err := docker.NewClient(dockerHost)

	if err != nil {
		return err
	}

	c, err := client.CreateContainer(docker.CreateContainerOptions{
		Name: name,
		Config: &docker.Config{
			Image: imageName,
		},
	})

	if err != nil {
		return err
	}

	err = client.StartContainer(c.ID, &docker.HostConfig{})

	if err != nil {
		return err
	}

	return nil
}
