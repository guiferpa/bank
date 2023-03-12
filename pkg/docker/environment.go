package docker

import "github.com/docker/docker/client"

type Environment struct {
	client *client.Client
}

func NewEnvironment() (*Environment, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Environment{client: cli}, nil
}
