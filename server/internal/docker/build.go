package docker

import (
	"context"
	"io"

	"docker-scout/internal/model"
	"docker-scout/internal/util"

	"github.com/moby/moby/client"
)

func BuildImage(ctx context.Context, cli *client.Client, req model.BuildImageRequest) (io.ReadCloser, error) {
	tar, err := util.ArchiveDirectory(req.ContextPath)
	if err != nil {
		return nil, err
	}
	defer tar.Close()

	resp, err := cli.ImageBuild(
		ctx,
		tar,
		client.ImageBuildOptions{
			Tags:       []string{req.Tag},
			Dockerfile: req.Dockerfile,
			NoCache:    req.NoCache,
			Remove:     true,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
