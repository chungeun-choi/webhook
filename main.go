package main

import (
	"github.com/chungeun-choi/webhook/internal/config"
	"github.com/chungeun-choi/webhook/internal/handlers"
	"github.com/chungeun-choi/webhook/internal/server"
	"log"
)

func main() {
	// Load the configuration
	cfg, err := config.LoadConfig("app.yml")
	if err != nil {
		log.Fatalf("Failed to load the configuration: %v", err)
	}

	// Create a new newServer
	newServer := server.NewServer(cfg)

	// Add the handlers
	if err = handlers.InitHandler(newServer); err != nil {
		log.Fatalf("Failed to initialize handlers: %v", err)
	}

	// Run the newServer
	newServer.Run()
}
