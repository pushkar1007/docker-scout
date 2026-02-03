package docker

import (
	"context"

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
