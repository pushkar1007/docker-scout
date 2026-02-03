package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"docker-scout/internal/api"
	"docker-scout/internal/docker"
	"docker-scout/internal/state"
)

func RunServer() {
	cli, err := docker.NewClient()
	if err != nil {
		log.Fatalf("docker client init failed: %v", err)
	}

	cache := state.NewCache()
	bcast := state.NewBroadcaster(256)
	StartStatsUpdater(cli, cache, 5*time.Second)
	StartDockerEventStream(cli, bcast)

	mux := http.NewServeMux()
	api.RegisterRoutes(mux, api.Deps{
		Docker:        cli,
		Cache:         cache,
		Broadcaster:   bcast,
		DashboardPath: filepath.Join("web", "dashboard.html"),
	})

	addr := "0.0.0.0:3000"
	if envAddr := strings.TrimSpace(os.Getenv("DOCKER_SCOUT_ADDR")); envAddr != "" {
		addr = envAddr
	}

	fmt.Printf("Dashboard running at http://localhost%s\n", addr)
	fmt.Printf("SSE stream at http://localhost%s/events\n", addr)
	fmt.Printf("WebSocket stats at ws://localhost%s/stats\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
