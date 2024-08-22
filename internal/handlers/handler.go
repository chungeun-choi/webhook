package handlers

import (
	"github.com/chungeun-choi/webhook/internal/server"
	"github.com/pkg/errors"
	"net/http"
)

func HandlerMain(server *server.Server) {
	server.AddHandler("/main", map[string]map[string]http.HandlerFunc{
		"": {"GET": getMainHandler},
	})

}

func getMainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Hello, World!"))
}

func InitHandler(server *server.Server) error {
	HandlerMain(server)
	if RegisterMutatingHandler(server) != nil {
		return errors.New("failed to register mutating handler")
	}
	RegisterPatchHandlers(server)

	return nil
}
