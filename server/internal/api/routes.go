package api

import (
	"net/http"
	"os"
	"path/filepath"

	"docker-scout/internal/state"

	"github.com/moby/moby/client"
)

type Deps struct {
	Docker        *client.Client
	Cache         *state.Cache
	Broadcaster   *state.Broadcaster
	DashboardPath string
}

func RegisterRoutes(mux *http.ServeMux, deps Deps) {
	registerEvents(mux, deps)
	registerStats(mux, deps)
	registerNetworks(mux, deps)
	registerContainers(mux, deps)
	registerImages(mux, deps)
	registerVolumes(mux, deps)
	registerSystem(mux, deps)
	registerDashboard(mux, deps)
}

func registerDashboard(mux *http.ServeMux, deps Deps) {
	path := deps.DashboardPath
	if path == "" {
		path = filepath.Join("web", "dashboard.html")
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := os.ReadFile(path)
		if err != nil {
			http.Error(w, "dashboard not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(body)
	})
}
