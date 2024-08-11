package patch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func UpdatePatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RequestOpBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	endpointPath := r.URL.Query().Get("endpoint")
	if endpointPath == "" {
		http.Error(w, "Missing endpoint query parameter", http.StatusBadRequest)
		return
	}

	pm, ok := ManagerMap[endpointPath]
	if !ok {
		log.Printf("Patch manager not found for endpoint %s", endpointPath)
		http.Error(w, fmt.Sprintf("Patch manager not found for endpoint %s", endpointPath), http.StatusBadRequest)
		return
	}

	err = pm.updatePatchOperation(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func AddPatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RequestOpBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	pm := NewPatchManager(req.EndpointPath)
	err = pm.AddPatchOperations(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetPatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	endpointPath := r.URL.Query().Get("endpoint")
	if endpointPath == "" {
		http.Error(w, "Missing endpoint query parameter", http.StatusBadRequest)
		return
	}

	pm, ok := ManagerMap[endpointPath]
	if !ok {
		http.Error(w, "Patch manager not found", http.StatusNotFound)
		return
	}

	patchOps := pm.GetPatchOperations()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(patchOps)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func ClearPatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	endpointPath := r.URL.Query().Get("endpoint")
	if endpointPath == "" {
		http.Error(w, "Missing endpoint query parameter", http.StatusBadRequest)
		return
	}

	pm, ok := ManagerMap[endpointPath]
	if !ok {
		http.Error(w, "Patch manager not found", http.StatusNotFound)
		return
	}

	pm.ClearPatchOperations()
	w.WriteHeader(http.StatusOK)
}
