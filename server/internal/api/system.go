package api

import (
	"encoding/json"
	"net/http"

	"docker-scout/internal/docker"

	"github.com/moby/moby/client"
)

func registerSystem(mux *http.ServeMux, deps Deps) {
	mux.HandleFunc("/system/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	mux.HandleFunc("/system/nuke", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.URL.Query().Get("confirm") != "true" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "confirmation required - pass ?confirm=true"})
			return
		}

		ctx := r.Context()

		containers, _ := docker.ListContainers(ctx, deps.Docker, true)
		for _, c := range containers.Items {
			_ = docker.RemoveContainer(ctx, deps.Docker, c.ID, &client.ContainerRemoveOptions{
				Force:         true,
				RemoveVolumes: true,
			})
		}

		images, _ := docker.ListImages(ctx, deps.Docker, true)
		for _, img := range images.Items {
			for _, tag := range img.RepoTags {
				_ = docker.RemoveImage(ctx, deps.Docker, tag, true)
			}
		}

		_, _ = docker.PruneVolumes(ctx, deps.Docker)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "system cleaned"})
	})
}
