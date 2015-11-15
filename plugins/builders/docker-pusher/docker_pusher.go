package dockerpusher

import (
	"github.com/carlosdp/supply-chain/builder"
	"github.com/carlosdp/supply-chain/options"
	"github.com/fsouza/go-dockerclient"
)

type dockerPusher struct {
	imageName string
}

func init() {
	builder.RegisterBuilder("docker-pusher", &dockerPusher{})
}

func (dp *dockerPusher) New() builder.Builder {
	return &dockerPusher{}
}

func (dp *dockerPusher) Build(workDir, image string, opts options.Options) (string, error) {
	username := opts.GetString("username")
	password := opts.GetString("password")

	client, err := docker.NewClientFromEnv()
	if err != nil {
		return "", err
	}

	err = client.PushImage(docker.PushImageOptions{
		Name: image,
	}, docker.AuthConfiguration{
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", err
	}

	return image, nil
}
