package bootstrap

import (
	"fmt"
	"github.com/chungeun-choi/webhook/bootstrap/config"
	"log"
	"os"
)

var (
	Namespace      string
	KubeConfigPath string
	KubeAPIToken   string
	IsRunningInK8S bool
	ServerConfigs  *config.ServerConfig
)

func Init() {
	var err error

	// load envs from kubernetes
	Namespace = os.Getenv("POD_NAMESPACE")
	if Namespace != "" {
		IsRunningInK8S = true
	}
	KubeConfigPath = os.Getenv("KUBECONFIG")

	// load server configs
	if ServerConfigs, err = config.LoadConfig("app.yml"); err != nil {
		log.Fatalf("Error loading server configs: %v", err)
	}

	// load kube api token
	if err = loadKubeAPIToken(); err != nil {
		log.Fatalf("Error loading kube api token: %v", err)
	}

	// create cert and key files

}

func loadKubeAPIToken() error {
	// load kube api token
	token, err := os.ReadFile(ServerConfigs.TokenPath)
	if err != nil {
		log.Printf("Error reading token file: %v", err)
		return fmt.Errorf("failed to read token file: %w", err)
	}
	KubeAPIToken = string(token)

	return nil
}
