package main

import (
	module_config "github.com/chungeun-choi/webhook/bootstrap/config"
	"github.com/chungeun-choi/webhook/bootstrap/server"
	"github.com/chungeun-choi/webhook/service/patch"
	"log"
)

func main() {
	var (
		config *module_config.ServerConfig
		err    error
	)
	if config, err = module_config.LoadConfig("./app.yml"); err != nil {
		log.Panicf("Error loading server configs: %v", err)
	}

	server := server.NewServer(config)
	patch.RegisterHandlers(server)

	server.Run()
}
