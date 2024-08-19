package handlers

import (
	"encoding/json"
	"fmt"
	server2 "github.com/chungeun-choi/webhook/internal/server"
	"github.com/chungeun-choi/webhook/pkg/patch"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

func RegisterPatchHandlers(s *server2.Server) {
	s.AddHandler("/patch", map[string]map[string]http.HandlerFunc{
		"":                    {"POST": addPatchHandler, "GET": getPatchHandler},
		"/{endpoint}":         {"POST": updatePatchHandler, "DELETE": deletePatchHandler},
		"/{endpoint}/clear":   {"GET": clearPatchHandler},
		"/{endpoint}/trigger": {"POST": triggerPatchHandler},
	})
}

// GetPatchHandler returns the patch operations for the given endpoint
func getPatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req string

	// If no endpoint query parameter is provided, return the list of patch operations
	if req = r.URL.Query().Get("endpoint"); req == "" {
		rsp := new(patch.ResponsePatchList)

		if len(patch.ManagerMap) == 0 {
			w.WriteHeader(http.StatusNoContent)
		}

		for k, v := range patch.ManagerMap {
			rsp.PatchList = append(rsp.PatchList, patch.ResponsePatch{
				EndpointPath:    k,
				PatchOperations: v.GetPatchOperations(),
			})
		}

		// Return the list of patch operations
		server2.WriteJson(w, http.StatusOK, rsp)
	} else {
		rsp := new(patch.ResponsePatch)

		if result, ok := patch.ManagerMap[req]; !ok {
			http.Error(w, "Patch manager not found", http.StatusNotFound)
			return
		} else {
			rsp.EndpointPath = req
			rsp.PatchOperations = result.GetPatchOperations()
		}

		// Return the list of patch operations as JSON
		server2.WriteJson(w, http.StatusOK, rsp)
	}

}

// AddPatchHandler adds a new patch operation
func addPatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var (
		req *patch.RequestPatch
		rsp *patch.ResponseBody
	)

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	pm := patch.NewPatchManager(req.EndpointPath)
	err = pm.AddPatchOperations(*req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rsp = new(patch.ResponseBody)
	rsp.Message = "Patch operation added successfully"
	server2.WriteJson(w, http.StatusOK, rsp)
}

// UpdatePatchHandler updates the patch operation with the given id
func updatePatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req patch.RequestPatch
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	pm, ok := patch.ManagerMap[req.EndpointPath]
	if !ok {
		log.Printf("Patch manager not found for endpoint %s", req.EndpointPath)
		http.Error(w, fmt.Sprintf("Patch manager not found for endpoint %s", req.EndpointPath), http.StatusBadRequest)
		return
	}

	err = pm.UpdatePatchOperation(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	server2.WriteJson(w, http.StatusOK, req)
}

// ClearPatchHandler clears the patch operations for the given endpoint
func clearPatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	endpointPath := r.URL.Query().Get("endpoint")
	if endpointPath == "" {
		http.Error(w, "Missing endpoint query parameter", http.StatusBadRequest)
		return
	}

	pm, ok := patch.ManagerMap[endpointPath]
	if !ok {
		http.Error(w, "Patch manager not found", http.StatusNotFound)
		return
	}

	pm.ClearPatchOperations()

	server2.WriteJson(w, http.StatusOK, nil)
}

func deletePatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	endpointPath := r.URL.Query().Get("endpoint")
	if endpointPath == "" {
		http.Error(w, "Missing endpoint query parameter", http.StatusBadRequest)
		return
	}

	_, ok := patch.ManagerMap[endpointPath]
	if !ok {
		http.Error(w, "Patch manager not found", http.StatusNotFound)
		return
	}

	delete(patch.ManagerMap, endpointPath)

	server2.WriteJson(w, http.StatusOK, nil)
}

func triggerPatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req string
	if req = r.URL.Query().Get("endpoint"); req == "" {
		http.Error(w, "Missing endpoint query parameter", http.StatusBadRequest)
		return
	}

	response := patch.Patch(req)
	if rsp, err := json.Marshal(response); err != nil {
		server2.WriteJson(w, http.StatusInternalServerError, &patch.ResponseBody{
			Message: errors.Wrap(err, "Failed to encode response").Error(),
		})
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(rsp); err != nil {
			fmt.Printf("Failed to write response: %v", err)
		}
	}
}
