package service

import (
	"log"
	"net/http"
)

type Handler struct {
	Server *http.Server
	log    *log.Logger
}

func NewHandler(server *http.Server, log *log.Logger) *Handler {
	return &Handler{
		Server: server,
		log:    log,
	}
}
