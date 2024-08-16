package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chungeun-choi/webhook/bootstrap/kubernetes"
	"github.com/chungeun-choi/webhook/internal/server"
	"github.com/chungeun-choi/webhook/pkg/mutating"
	"github.com/gorilla/mux"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"net/http"
)

var mutatingManger *mutating.MutatingManager

func RegisterMutatingHandler(s *server.Server) error {
	kubeClient, err := kubernetes.CreateClientSet(s.Config.KubeAPIServerURL, s.Config.Token)
	if err != nil {
		return errors.New("failed to create Kubernetes clientset")
	}

	mutatingManger = mutating.NewMutateManager(
		kubeClient,
		s.Config.AdmissionReviewVersion,
		fmt.Sprintf("http://%s:%d", s.Config.Hostname, s.Config.Port),
	)

	s.AddHandler("/mutating", map[string]map[string]http.HandlerFunc{
		"":        {"POST": createUpdateMutatingHandler},
		"/{name}": {"DELETE": deleteMutatingHandler, "GET": getMutatingHandler},
	})

	return nil
}

func getMutatingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var (
		req string
		rsp *mutating.ResponseGetRulesBody
	)

	vars := mux.Vars(r)
	if req = vars["name"]; req == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	} else {
		// Get the list of mutating webhook configurations
		if result, err := mutatingManger.Get(req); err != nil {
			// If the mutating webhook configuration is not found, return 204
			if apierrors.IsNotFound(err) {
				server.WriteJson(w, http.StatusNoContent,
					&mutating.ResponseBody{
						Message: fmt.Sprintf("Mutating webhook configuration %s not found", req),
					},
				)
			}
			server.WriteJson(w, http.StatusInternalServerError,
				&mutating.ResponseBody{
					Message: fmt.Sprintf("Failed to list mutating webhooks: %s", err.Error()),
				},
			)
			return
		} else {
			rsp = new(mutating.ResponseGetRulesBody)
			rsp.Message = "List of mutating webhook configurations"
			rsp.WebhookConfiguration = result.MutatingWebhookConfiguration
			server.WriteJson(w, http.StatusOK, rsp)
		}
	}
}

// addMutatingHandler adds a new mutating webhook configuration
func createUpdateMutatingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var (
		req *mutating.RequestAddRulesBody
		rsp *mutating.ResponseAddRulesBody
	)

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Register the mutating webhook configuration
	if result, err := mutatingManger.Register(*req); err != nil {
		server.WriteJson(w, http.StatusInternalServerError,
			&mutating.ResponseBody{
				Message: fmt.Sprintf("Failed to register mutating webhook: %s", err.Error()),
			},
		)
		return
	} else {
		rsp = new(mutating.ResponseAddRulesBody)
		rsp.Message = "Mutating webhook configuration added successfully"
		rsp.ConfigBuilder = *result
		server.WriteJson(w, http.StatusOK, result)
	}
}

// deleteMutatingHandler deletes an existing mutating webhook configuration
func deleteMutatingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var (
		req string
		rsp *mutating.ResponseBody
	)

	// Get the name of the mutating webhook configuration from the request
	vars := mux.Vars(r)
	if req = vars["name"]; req == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	} else {
		// Delete the mutating webhook configuration
		if err := mutatingManger.Delete(req); err != nil {
			server.WriteJson(w, http.StatusInternalServerError,
				mutating.ResponseBody{Message: fmt.Sprintf("Failed to delete mutating webhook: %s", err.Error())},
			)
			http.Error(w, "Failed to delete mutating webhook", http.StatusInternalServerError)
			return
		}

		rsp = new(mutating.ResponseBody)
		rsp.Message = "Mutating webhook configuration deleted successfully"
		server.WriteJson(w, http.StatusOK, rsp)
	}
}
