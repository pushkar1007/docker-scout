package docker

import (
	"context"

	"docker-scout/internal/model"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

func ListContainers(ctx context.Context, cli *client.Client, all bool) (client.ContainerListResult, error) {
	return cli.ContainerList(ctx, client.ContainerListOptions{All: all})
}

func InspectContainer(ctx context.Context, cli *client.Client, id string) (client.ContainerInspectResult, error) {
	return cli.ContainerInspect(ctx, id, client.ContainerInspectOptions{})
}

func PauseContainer(ctx context.Context, cli *client.Client, id string) error {
	_, err := cli.ContainerPause(ctx, id, client.ContainerPauseOptions{})
	return err
}

func UnpauseContainer(ctx context.Context, cli *client.Client, id string) error {
	_, err := cli.ContainerUnpause(ctx, id, client.ContainerUnpauseOptions{})
	return err
}

func StartContainer(ctx context.Context, cli *client.Client, id string) error {
	_, err := cli.ContainerStart(ctx, id, client.ContainerStartOptions{})
	return err
}

func StopContainer(ctx context.Context, cli *client.Client, id string) error {
	_, err := cli.ContainerStop(ctx, id, client.ContainerStopOptions{})
	return err
}

func RemoveContainer(ctx context.Context, cli *client.Client, id string, opts *client.ContainerRemoveOptions) error {
	if opts == nil {
		opts = &client.ContainerRemoveOptions{}
	}
	_, err := cli.ContainerRemove(ctx, id, *opts)
	return err
}

func ResolveLabels(summaryLabels map[string]string, inspect client.ContainerInspectResult) map[string]string {
	if inspect.Container.Config != nil && len(inspect.Container.Config.Labels) > 0 {
		return inspect.Container.Config.Labels
	}
	return summaryLabels
}

func IsVolumeInUse(ctx context.Context, cli *client.Client, volumeName string) bool {
	containers, err := ListContainers(ctx, cli, true)
	if err != nil {
		return false
	}

	for _, c := range containers.Items {
		for _, m := range c.Mounts {
			if m.Type == "volume" && m.Name == volumeName {
				return true
			}
		}
	}
	return false
}

func CreateContainer(ctx context.Context, cli *client.Client, req model.CreateContainerRequest) (string, error) {
	// Build container config with all provided options
	config := &container.Config{
		Image: req.Image,
	}

	// Add labels if provided
	if len(req.Labels) > 0 {
		config.Labels = req.Labels
	}

	// Add environment variables if provided
	if len(req.Env) > 0 {
		config.Env = req.Env
	}

	// Add command if provided
	if len(req.Cmd) > 0 {
		config.Cmd = req.Cmd
	}

	// Build host config with networking and restart policy
	hostConfig := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: container.RestartPolicyMode(req.RestartPolicy),
		},
	}

	// Add network mode if provided
	if req.NetworkMode != "" {
		hostConfig.NetworkMode = container.NetworkMode(req.NetworkMode)
	}

	// Add volumes if provided
	if len(req.Volumes) > 0 {
		hostConfig.Binds = req.Volumes
	}

	opts := client.ContainerCreateOptions{
		Name:       req.Name,
		Config:     config,
		HostConfig: hostConfig,
	}

	resp, err := cli.ContainerCreate(ctx, opts)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}
