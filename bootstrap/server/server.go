package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chungeun-choi/webhook/bootstrap/config"
	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	config *config.ServerConfig
}

func NewServer(config *config.ServerConfig) *Server {
	return &Server{
		Router: mux.NewRouter(),
		config: config,
	}
}

func (s *Server) Run() {
	// 서버 시작 전 라우트 출력
	s.printRoutes()

	address := fmt.Sprintf("%v:%v", s.config.Hostname, s.config.Port)
	log.Printf("Starting server on %s\n", address)
	if err := http.ListenAndServe(address, s.Router); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

func (s *Server) printRoutes() {
	err := s.Router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return err
		}

		methods, err := route.GetMethods()
		if err != nil {
			// If there are no methods associated with the route, it might be a middleware
			methods = []string{"ANY"}
		}

		log.Printf("Registered route: %s %v", pathTemplate, methods)
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to walk routes: %v", err)
	}
}

func (s *Server) AddHandler(prefix string, handlerMap map[string]map[string]http.HandlerFunc) *Server {
	subRouter := s.Router.PathPrefix(prefix).Subrouter()
	for path, methods := range handlerMap {
		for method, handlerFunc := range methods {
			subRouter.HandleFunc(path, handlerFunc).Methods(method)
		}
	}
	return s
}
