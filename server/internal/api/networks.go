package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"docker-scout/internal/docker"
	"docker-scout/internal/model"
)

// registerNetworks wires the /networks, /networks/create, and /networks/remove endpoints.
// All operations include a 10-15 second timeout; requests blocking beyond this are terminated.
func registerNetworks(mux *http.ServeMux, deps Deps) {
	mux.HandleFunc("/networks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// GET /networks lists all Docker networks.
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		nets, err := docker.ListNetworks(ctx, deps.Docker)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(nets)
	})

	mux.HandleFunc("/networks/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// POST /networks/create creates a new network. Request must include:
		// - name: unique network name
		// - driver: network driver (e.g., "bridge", "overlay")
		// - options: driver-specific options (optional)
		var req model.CreateNetworkRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid payload"})
			return
		}

		if req.Name == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "network name required"})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()

		resp, err := docker.CreateNetwork(ctx, deps.Docker, req.Name, req.Driver, req.Options)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/networks/remove", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// DELETE /networks/remove?id=<id> or ?name=<name> removes a network.
		// Docker forbids removal if the network is connected to any running containers.
		id := r.URL.Query().Get("id")
		if id == "" {
			id = r.URL.Query().Get("name")
		}
		if id == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "network id or name required"})
			return
		}

		if err := docker.RemoveNetwork(r.Context(), deps.Docker, id); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "removed", "id": id})
	})
}
