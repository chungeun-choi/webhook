package handlers

import (
	"github.com/chungeun-choi/webhook/internal/server"
	"github.com/pkg/errors"
)

func InitHandler(server *server.Server) error {
	if RegisterMutatingHandler(server) != nil {
		return errors.New("failed to register mutating handler")
	}
	RegisterPatchHandlers(server)

	return nil
}
