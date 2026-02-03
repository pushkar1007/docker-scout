package docker

import (
	"context"

	"github.com/moby/moby/client"
)

// ListNetworks queries the Docker daemon for all networks (both user and system-managed).
// Returns the raw NetworkListResult containing all network metadata and connectivity info.
func ListNetworks(
	ctx context.Context,
	cli *client.Client,
) (client.NetworkListResult, error) {

	return cli.NetworkList(ctx, client.NetworkListOptions{})
}

// CreateNetwork creates a new network with the given driver and options.
// CheckDuplicate is always true to prevent accidental duplicates.
// The caller must handle name uniqueness in the API layer.
func CreateNetwork(
	ctx context.Context,
	cli *client.Client,
	name string,
	driver string,
	options map[string]string,
) (client.NetworkCreateResult, error) {

	req := client.NetworkCreateOptions{
		Driver:  driver,
		Options: options,
	}

	return cli.NetworkCreate(ctx, name, req)
}

// RemoveNetwork deletes a network by ID or name. Docker requires the network to not be
// connected to any running containers; the daemon returns an error if removal is impossible.
func RemoveNetwork(
	ctx context.Context,
	cli *client.Client,
	id string,
) error {

	_, err := cli.NetworkRemove(ctx, id, client.NetworkRemoveOptions{})
	return err
}
