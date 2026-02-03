package api

import (
	"encoding/json"
	"net/http"
)

func registerStats(mux *http.ServeMux, deps Deps) {
	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		data := deps.Cache.Get()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	})
}
