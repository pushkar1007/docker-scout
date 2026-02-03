package docker

import (
	"context"
	"sort"

	"docker-scout/internal/model"

	"github.com/moby/moby/client"
)

func ListVolumes(ctx context.Context, cli *client.Client) ([]model.VolumeSummary, error) {
	resp, err := cli.VolumeList(ctx, client.VolumeListOptions{})
	if err != nil {
		return nil, err
	}

	vols := make([]model.VolumeSummary, 0, len(resp.Items))
	for _, v := range resp.Items {
		inUse := IsVolumeInUse(ctx, cli, v.Name)
		vols = append(vols, model.VolumeSummary{
			Name:       v.Name,
			Driver:     v.Driver,
			Mountpoint: v.Mountpoint,
			Scope:      v.Scope,
			Labels:     v.Labels,
			CreatedAt:  v.CreatedAt,
			InUse:      inUse,
		})
	}

	sort.Slice(vols, func(i, j int) bool {
		return vols[i].Name < vols[j].Name
	})

	return vols, nil
}

func InspectVolume(ctx context.Context, cli *client.Client, name string) (model.VolumeSummary, error) {
	vol, err := cli.VolumeInspect(ctx, name, client.VolumeInspectOptions{})
	if err != nil {
		return model.VolumeSummary{}, err
	}

	v := vol.Volume
	return model.VolumeSummary{
		Name:       v.Name,
		Driver:     v.Driver,
		Mountpoint: v.Mountpoint,
		Scope:      v.Scope,
		Labels:     v.Labels,
		CreatedAt:  v.CreatedAt,
		InUse:      IsVolumeInUse(ctx, cli, v.Name),
	}, nil
}

func CreateVolume(ctx context.Context, cli *client.Client, req model.CreateVolumeRequest) error {
	_, err := cli.VolumeCreate(ctx, client.VolumeCreateOptions{
		Name:       req.Name,
		Driver:     req.Driver,
		Labels:     req.Labels,
		DriverOpts: req.Options,
	})
	return err
}

func RemoveVolume(ctx context.Context, cli *client.Client, name string) error {
	_, err := cli.VolumeRemove(ctx, name, client.VolumeRemoveOptions{})
	return err
}

func PruneVolumes(ctx context.Context, cli *client.Client) (client.VolumePruneResult, error) {
	return cli.VolumePrune(ctx, client.VolumePruneOptions{})
}
