package docker

import (
	"context"
	"io"

	"github.com/moby/moby/client"
)

func ListImages(ctx context.Context, cli *client.Client, all bool) (client.ImageListResult, error) {
	return cli.ImageList(ctx, client.ImageListOptions{All: all})
}

func RemoveImage(ctx context.Context, cli *client.Client, name string, force bool) error {
	_, err := cli.ImageRemove(
		ctx,
		name,
		client.ImageRemoveOptions{
			Force:         force,
			PruneChildren: true,
		},
	)
	return err
}

func PullImage(ctx context.Context, cli *client.Client, name string) error {
	rc, err := cli.ImagePull(ctx, name, client.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer rc.Close()
	_, _ = io.Copy(io.Discard, rc)
	return nil
}
