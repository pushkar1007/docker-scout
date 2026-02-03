package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"docker-scout/internal/docker"

	"github.com/gorilla/websocket"
)

var statsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func registerStats(mux *http.ServeMux, deps Deps) {
	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		conn, err := statsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("stats ws upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Send stats every second
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		// Handle ping/pong for connection health
		go func() {
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					cancel()
					return
				}
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Get all running containers
				containers, err := docker.ListContainers(ctx, deps.Docker, false)
				if err != nil {
					log.Printf("failed to list containers: %v", err)
					continue
				}

				// Collect stats for all containers
				allStats := make(map[string]interface{})
				for _, c := range containers.Items {
					stats, err := docker.ReadContainerStats(ctx, deps.Docker, c.ID)
					if err != nil {
						log.Printf("failed to read stats for %s: %v", c.ID, err)
						continue
					}

					allStats[c.ID] = map[string]interface{}{
						"id":            c.ID,
						"name":          c.Names[0],
						"cpu_percent":   stats.CPUPercent,
						"memory_bytes":  stats.MemoryBytes,
						"memory_limit":  stats.MemoryLimit,
						"net_rx_bps":    stats.NetRxBps,
						"net_tx_bps":    stats.NetTxBps,
						"disk_read_bps": stats.DiskReadBps,
						"disk_write_bps": stats.DiskWriteBps,
					}
				}

				// Send stats to frontend
				data, err := json.Marshal(map[string]interface{}{
					"type":  "stats",
					"stats": allStats,
					"timestamp": time.Now().Unix(),
				})
				if err != nil {
					log.Printf("failed to marshal stats: %v", err)
					continue
				}

				if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
					log.Printf("failed to write stats: %v", err)
					return
				}
			}
		}
	})
}
