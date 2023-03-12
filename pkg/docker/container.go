package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

func (e *Environment) RunContainer(ctx context.Context, image, port, bindingPort string, env []string) (string, error) {
	containerConfig := &container.Config{
		Image:        image,
		ExposedPorts: nat.PortSet{nat.Port(bindingPort): struct{}{}},
		Env:          env,
	}
	containerHostConfig := &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{nat.Port(bindingPort): {{HostIP: "127.0.0.1", HostPort: port}}},
	}
	resp, err := e.client.ContainerCreate(ctx, containerConfig, containerHostConfig, nil, nil, "")
	if err != nil {
		return "", err
	}
	if err := e.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (e *Environment) KillContainer(ctx context.Context, containerID string) error {
	if err := e.client.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true}); err != nil {
		return err
	}

	return nil
}
