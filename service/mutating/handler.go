package mutating

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	GroupName = "mutating"
)

// Define a handler for registering or updating a mutating webhook configuration
func registerHandler(mutateManager *MutatingManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var config ConfigBuilder
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		var (
			result *ConfigBuilder
			err    error
		)
		if r.Method == http.MethodPost {
			result, err = mutateManager.Register(&config)
		} else {
			existingConfig, err := mutateManager.getOldConfig(config.Name)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to get existing config: %v", err), http.StatusInternalServerError)
				return
			}
			result, err = mutateManager.Update(&config, existingConfig)
		}

		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to register or update config: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

// Define a handler for deleting a mutating webhook configuration
func deleteHandler(mutateManager *MutatingManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}

		err := mutateManager.Delete(name)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete config: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Deleted successfully")
	}
}

// NewMux creates a new HTTP server mux for mutating webhook operations
func NewMux(mutateManager *MutatingManager) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(GroupName+"/register", registerHandler(mutateManager))
	mux.HandleFunc(GroupName+"/delete", deleteHandler(mutateManager))
	return mux
}
