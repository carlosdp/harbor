package container

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"os"
)

func NewContainer(imageName string, deployId string) {
	containerName := imageName + "-" + deployId
	dockerHost := os.Getenv("DOCKER_HOST")

	if dockerHost == "" {
		dockerHost = "unix:///var/run/docker.sock"
	}

	client, err := docker.NewClient(dockerHost)

	if err != nil {
		panic(err)
	}

	containers, err := client.ListContainers(docker.ListContainersOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}

	for _, c := range containers {
		err = client.StopContainer(c.ID, 5)
		if err != nil {
			fmt.Println(err)
		}
		err = client.RemoveContainer(docker.RemoveContainerOptions{
			ID:            c.ID,
			RemoveVolumes: true,
		})
		if err != nil {
			panic(err)
		}
	}

	c, err := client.CreateContainer(docker.CreateContainerOptions{
		Name: containerName,
		Config: &docker.Config{
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Image:        imageName,
			Cmd:          []string{"/bin/echo", "Hail Voltron"},
		},
	})

	if err != nil {
		panic(err)
	}

	err = client.StartContainer(c.ID, &docker.HostConfig{})

	if err != nil {
		panic(err)
	}

	err = client.AttachToContainer(docker.AttachToContainerOptions{
		Container:    c.ID,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		Logs:         true,
		Stream:       false,
		Stdout:       true,
		Stderr:       true,
	})

	if err != nil {
		panic(err)
	}
}
