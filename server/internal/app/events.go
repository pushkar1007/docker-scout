package app

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"strings"
	"time"

	"docker-scout/internal/docker"
	"docker-scout/internal/model"
	"docker-scout/internal/state"
	"docker-scout/internal/util"

	"github.com/moby/moby/api/types/events"
	"github.com/moby/moby/client"
)

type containerEventPayload struct {
	Type       string                `json:"type"`
	Action     string                `json:"action"`
	ID         string                `json:"id"`
	Time       string                `json:"time"`
	Attributes map[string]string     `json:"attributes,omitempty"`
	Container  *model.ContainerStats `json:"container,omitempty"`
}

func StartDockerEventStream(cli *client.Client, bcast *state.Broadcaster) {
	go func() {
		filters := make(client.Filters).Add("type", "container")
		for {
			ctx, cancel := context.WithCancel(context.Background())
			result := cli.Events(ctx, client.EventsListOptions{Filters: filters})

			for {
				select {
				case evt, ok := <-result.Messages:
					if !ok {
						cancel()
						goto restart
					}
					payload := buildContainerEvent(ctx, cli, evt)
					if payload == nil {
						continue
					}
					if data, err := json.Marshal(payload); err == nil {
						_ = bcast.Send(data)
					}
				case err := <-result.Err:
					cancel()
					if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, io.EOF) {
						log.Printf("docker events stream error: %v", err)
					}
					goto restart
				}
			}

		restart:
			time.Sleep(2 * time.Second)
		}
	}()
}

func buildContainerEvent(ctx context.Context, cli *client.Client, evt events.Message) *containerEventPayload {
	payload := &containerEventPayload{
		Type:       string(evt.Type),
		Action:     string(evt.Action),
		ID:         evt.Actor.ID,
		Time:       time.Unix(evt.Time, 0).UTC().Format(time.RFC3339),
		Attributes: evt.Actor.Attributes,
	}

	inspect, err := docker.InspectContainer(ctx, cli, evt.Actor.ID)
	if err != nil || inspect.Container.ID == "" {
		return payload
	}

	stats, statsErr := docker.ReadContainerStats(ctx, cli, evt.Actor.ID)

	name := ""
	if inspect.Container.Name != "" {
		name = strings.TrimPrefix(inspect.Container.Name, "/")
	}

	state := "unknown"
	if inspect.Container.State != nil {
		state = string(inspect.Container.State.Status)
	}

	payload.Container = &model.ContainerStats{
		ID:       evt.Actor.ID[:12],
		Name:     name,
		State:    state,
		LastUsed: util.ResolveLastUsed(inspect),
		Image:    inspect.Container.Image,
		Labels:   docker.ResolveLabels(nil, inspect),
		Ports:    "-",
		CPU:      util.FormatCPU(stats.CPUPercent, statsErr == nil),
		Memory:   util.FormatMemory(stats.MemoryBytes, stats.MemoryLimit, statsErr == nil),
		NetIO:    util.FormatNetIO(stats.NetRxBps, stats.NetTxBps, statsErr == nil),
		DiskIO:   util.FormatDiskIO(stats.DiskReadBps, stats.DiskWriteBps, statsErr == nil),
	}

	return payload
}
