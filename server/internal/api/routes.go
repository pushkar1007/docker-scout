package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/moby/moby/client"
)

type Deps struct {
	Docker        *client.Client
	DashboardPath string
}

func RegisterRoutes(mux *http.ServeMux, deps Deps) {
	registerStats(mux, deps)
	registerTerminal(mux, deps)
	registerNetworks(mux, deps)
	registerContainers(mux, deps)
	registerImages(mux, deps)
	registerVolumes(mux, deps)
	registerSystem(mux, deps)
	registerDashboard(mux, deps)
}

func registerDashboard(mux *http.ServeMux, deps Deps) {
	dashboardPath := deps.DashboardPath
	if dashboardPath == "" {
		dashboardPath = filepath.Join("web", "dashboard.html")
	}
	echoPath := filepath.Join("web", "index.html")

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/echo" || r.URL.Path == "/index" {
			body, err := os.ReadFile(echoPath)
			if err != nil {
				http.Error(w, "echo page not found", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write(body)
			return
		}

		body, err := os.ReadFile(dashboardPath)
		if err != nil {
			http.Error(w, "dashboard not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(body)
	})
}
