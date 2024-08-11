package config_test

import (
	"github.com/chungeun-choi/webhook/bootstrap/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test data for app.yaml
const testYAMLContent = `
name: my-webhook-server
hostname: example.com
admission_review_version:
  - v1
  - v1beta1
port: 8443
cert_file: /path/to/certfile.crt
key_file: /path/to/keyfile.key
service_name: my-webhook-service
`

// Test file path
const testFilePath = "app.yaml"

// setUp creates the test YAML file.
func setUp(t *testing.T) {
	err := os.WriteFile(testFilePath, []byte(testYAMLContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}
}

// tearDown deletes the test YAML file.
func tearDown(t *testing.T) {
	err := os.Remove(testFilePath)
	if err != nil {
		t.Fatalf("Failed to delete test YAML file: %v", err)
	}
}

// TestLoadConfig tests the LoadConfig function.
func TestLoadConfig(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	// Load the configuration
	cfg, err := config.LoadConfig(testFilePath)
	assert.NoError(t, err, "Loading config should not produce an error")

	// Verify the configuration
	assert.Equal(t, "my-webhook-server", cfg.Name)
	assert.Equal(t, "example.com", cfg.Hostname)
	assert.Equal(t, []string{"v1", "v1beta1"}, cfg.AdmissionReviewVersion)
	assert.Equal(t, 8443, cfg.Port)
	assert.Equal(t, "/path/to/certfile.crt", cfg.CertFile)
	assert.Equal(t, "/path/to/keyfile.key", cfg.KeyFile)
	assert.Equal(t, "my-webhook-service", cfg.ServiceName)
}
