package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

func registerStats(mux *http.ServeMux, deps Deps) {
	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		log.Println("New WebSocket client connected to /stats")

		// Create a ticker that sends stats every second
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		// Channel to signal when to stop
		done := make(chan struct{})

		// Read messages from client (to detect disconnect)
		go func() {
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					close(done)
					return
				}
			}
		}()

		// Send stats to client every second
		for {
			select {
			case <-ticker.C:
				data := deps.Cache.Get()
				if err := conn.WriteJSON(data); err != nil {
					log.Printf("WebSocket write error: %v", err)
					return
				}
			case <-done:
				log.Println("WebSocket client disconnected")
				return
			}
		}
	})
}
