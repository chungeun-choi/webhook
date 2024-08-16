package server

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, statusCode int, rsp interface{}) {
	// Return the list of patch operations as JSON
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	if rsp != nil {
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
