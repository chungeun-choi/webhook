package config

import (
	"fmt"
	"github.com/chungeun-choi/webhook/bootstrap/cert"
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
	CaFile                 string   `yaml:"ca_file"`                  // path to the x509 certificate authority file
	ServiceName            string   `yaml:"service_name"`             // webhook pkg name in k8s
	KubeAPIServerURL       string   `yaml:"kube_api_server_url"`      // k8s cluster host
	AdmissionFailurePolicy string   `yaml:"admission_failure_policy"` // admission failure policy
	TokenPath              string   `yaml:"token_path"`
	Token                  string
	IsPod                  bool
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

	if config.GenerateCert() != nil {
		return nil, fmt.Errorf("failed to generate cert: %w", err)
	}

	// Load the token.txt.
	if config.TokenPath != "" {
		if config.loadToken() != nil {
			return nil, fmt.Errorf("failed to load token.txt: %w", err)
		}
	}

	if config.AdmissionFailurePolicy == "" {
		config.AdmissionFailurePolicy = "Ignore"
	}

	// Check if running in a test_patch.
	config.IsPod = checkRunningInPod()

	return &config, nil
}

// LoadToken reads a token.txt file and returns the token.txt.
func (c *ServerConfig) loadToken() error {
	// Read the token.txt file.
	data, err := os.ReadFile(c.TokenPath)
	if err != nil {
		log.Printf("Error reading token.txt file: %v", err)
		return fmt.Errorf("failed to read token.txt file: %w", err)
	}

	c.Token = string(data)

	return nil
}

func checkRunningInPod() bool {
	if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token"); err == nil {
		return true
	}

	return false
}

// GenerateCert generates a self-signed certificate if the key and cert files are not provided.
func (c *ServerConfig) GenerateCert() error {
	if c.KeyFile == "" || c.CertFile == "" {
		if info, err := cert.GenerateCert([]string{"self-signed-cert"}, []string{"*"}, "localhost"); err != nil {
			log.Fatalf(" Failed to generate cert: %v", err)
			return err
		} else {
			c.KeyFile = info.CertKeyPath
			c.CertFile = info.CertPath
			c.CaFile = info.CaCertPath
		}
	}

	return nil
}
