package patch

import (
	"encoding/json"
	"fmt"
	"github.com/chungeun-choi/webhook/bootstrap/server"
	"log"
	"net/http"
)

func RegisterHandlers(s *server.Server) {
	s.AddHandler("/patch", map[string]map[string]http.HandlerFunc{
		"":                  {"POST": AddPatchHandler, "GET": GetPatchHandler},
		"/{endpoint}":       {"POST": UpdatePatchHandler, "DELETE": DeletePatchManager},
		"/{endpoint}/clear": {"GET": ClearPatchHandler},
	})
}

// GetPatchHandler returns the patch operations for the given endpoint
func GetPatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// If no endpoint query parameter is provided, return the list of patch operations
	if param := r.URL.Query().Get("endpoint"); param == "" {
		var rsp *ResponsePatchList = new(ResponsePatchList)

		if len(ManagerMap) == 0 {
			w.WriteHeader(http.StatusNoContent)
		}

		for k, v := range ManagerMap {
			rsp.PatchList = append(rsp.PatchList, ResponsePatch{
				EndpointPath:    k,
				PatchOperations: v.GetPatchOperations(),
			})
		}

		// Return the list of patch operations
		server.WriteJson(w, http.StatusOK, rsp)
	} else {
		var rsp *ResponsePatch = new(ResponsePatch)

		if result, ok := ManagerMap[param]; !ok {
			http.Error(w, "Patch manager not found", http.StatusNotFound)
			return
		} else {
			rsp.EndpointPath = param
			rsp.PatchOperations = result.GetPatchOperations()
		}

		// Return the list of patch operations as JSON
		server.WriteJson(w, http.StatusOK, rsp)
	}

}

// AddPatchHandler adds a new patch operation
func AddPatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RequestPatch
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

	server.WriteJson(w, http.StatusOK, req)
}

// UpdatePatchHandler updates the patch operation with the given id
func UpdatePatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RequestPatch
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	pm, ok := ManagerMap[req.EndpointPath]
	if !ok {
		log.Printf("Patch manager not found for endpoint %s", req.EndpointPath)
		http.Error(w, fmt.Sprintf("Patch manager not found for endpoint %s", req.EndpointPath), http.StatusBadRequest)
		return
	}

	err = pm.updatePatchOperation(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	server.WriteJson(w, http.StatusOK, req)
}

// ClearPatchHandler clears the patch operations for the given endpoint
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

	server.WriteJson(w, http.StatusOK, nil)
}

func DeletePatchManager(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	endpointPath := r.URL.Query().Get("endpoint")
	if endpointPath == "" {
		http.Error(w, "Missing endpoint query parameter", http.StatusBadRequest)
		return
	}

	_, ok := ManagerMap[endpointPath]
	if !ok {
		http.Error(w, "Patch manager not found", http.StatusNotFound)
		return
	}

	delete(ManagerMap, endpointPath)

	server.WriteJson(w, http.StatusOK, nil)
}
