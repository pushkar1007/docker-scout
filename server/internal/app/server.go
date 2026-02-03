package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"docker-scout/internal/api"
	"docker-scout/internal/docker"
)

func RunServer() {
	cli, err := docker.NewClient()
	if err != nil {
		log.Fatalf("docker client init failed: %v", err)
	}

	mux := http.NewServeMux()
	api.RegisterRoutes(mux, api.Deps{
		Docker:        cli,
		DashboardPath: filepath.Join("web", "dashboard.html"),
	})

	addr := ":8089"
	if envAddr := strings.TrimSpace(os.Getenv("DOCKER_SCOUT_ADDR")); envAddr != "" {
		addr = envAddr
	}

	fmt.Printf("Dashboard running at http://localhost%s\n", addr)
	fmt.Printf("Stats WS at ws://localhost%s/stats\n", addr)
	fmt.Printf("Terminal WS at ws://localhost%s/terminal\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
