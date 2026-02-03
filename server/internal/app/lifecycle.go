package app

import (
	"time"

	"docker-scout/internal/state"

	"github.com/moby/moby/client"
)

func StartStatsUpdater(cli *client.Client, cache *state.Cache, interval time.Duration) {
	go func() {
		state.UpdateCache(cli, cache)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			state.UpdateCache(cli, cache)
		}
	}()
}
