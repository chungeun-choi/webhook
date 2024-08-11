package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// ServerConfig represents the configuration for the server.
type ServerConfig struct {
	Name                   string   `yaml:"name"`                     // webhook server name
	Hostname               string   `yaml:"hostname"`                 // webhook server host name
	AdmissionReviewVersion []string `yaml:"admission_review_version"` // webhook server api version
	Port                   int      `yaml:"port"`                     // webhook server port
	CertFile               string   `yaml:"cert_file"`                // path to the x509 certificate for https
	KeyFile                string   `yaml:"key_file"`                 // path to the x509 private key matching `CertFile`
	ServiceName            string   `yaml:"service_name"`             // webhook service name in k8s
	KubeAPIServerURL       string   `yaml:"kube_api_server_url"`      // k8s cluster host
	TokenPath              string   `yaml:"token"`                    // token for k8s cluster
}

// LoadConfig reads a YAML file and unmarshals its content into a ServerConfig struct.
func LoadConfig(filePath string) (*ServerConfig, error) {
	// Read the YAML file.
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading YAML file: %v", err)
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Create an instance of ServerConfig.
	var config ServerConfig

	// Unmarshal the YAML data into the ServerConfig struct.
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Printf("Error unmarshaling YAML data: %v", err)
		return nil, fmt.Errorf("failed to unmarshal YAML data: %w", err)
	}

	return &config, nil
}
