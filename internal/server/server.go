package server

import (
	"fmt"
	"github.com/chungeun-choi/webhook/internal/config"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	Config *config.ServerConfig
}

func NewServer(config *config.ServerConfig) *Server {
	return &Server{
		Router: mux.NewRouter(),
		Config: config,
	}
}

func (s *Server) Run() {
	// Print the registered routes
	s.printRoutes()

	address := fmt.Sprintf("%v:%v", s.Config.Hostname, s.Config.Port)
	log.Printf("Starting server on %s\n", address)

	// Run the server
	if err := http.ListenAndServeTLS(address, s.Config.CertFile, s.Config.KeyFile, s.Router); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
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
			wrappedHandler := s.logRequestDetails(handlerFunc)
			subRouter.HandleFunc(path, wrappedHandler).Methods(method)
		}
	}
	return s
}

func (s *Server) logRequestDetails(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		log.Printf("Handling %s request for %s", r.Method, r.URL.Path)
		handler.ServeHTTP(w, r)
		log.Printf("Completed handling %s request for %s in %v", r.Method, r.URL.Path, time.Since(startTime))
	}
}
