package mutating_test

import (
	"fmt"
	"github.com/chungeun-choi/webhook/bootstrap/config"
	"os"
	"testing"

	"github.com/chungeun-choi/webhook/bootstrap"
	"github.com/chungeun-choi/webhook/service/mutating"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes/fake" // Adjust the import path as necessary.
)

// Define a test case struct for CreateClientSet tests
type createClientSetTestCase struct {
	name           string
	setup          func() // Function to setup the test environment
	expectError    bool   // Whether we expect an error from CreateClientSet
	expectedClient bool   // Whether we expect ClientCache to be initialized
}

// TestCreateClientSet tests different scenarios for CreateClientSet using test cases
func TestCreateClientSet(t *testing.T) {
	// Define the test cases
	testCases := []createClientSetTestCase{
		{
			name: "WithKubeConfig",
			setup: func() {
				// Create a mock kubeconfig file
				kubeConfigContent := `
apiVersion: v1
kind: Config
clusters:
- cluster:
    server: http://localhost:8080
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
  name: test-context
current-context: test-context
users:
- name: test-user
  user:
    token: fake-token
`
				kubeConfigPath := "test_kubeconfig"
				err := os.WriteFile(kubeConfigPath, []byte(kubeConfigContent), 0644)
				require.NoError(t, err, "failed to create mock kubeconfig file")
				os.Setenv("KUBECONFIG", kubeConfigPath)
			},
			expectError:    false,
			expectedClient: true,
		},
		{
			name: "WithoutKubeConfig",
			setup: func() {
				// Log current state
				fmt.Printf("Initial ServerConfigs: %+v\n", bootstrap.ServerConfigs)

				// Remove the KUBECONFIG environment variable
				os.Unsetenv("KUBECONFIG")

				// Mock bootstrap values
				if bootstrap.ServerConfigs == nil {
					bootstrap.ServerConfigs = &config.ServerConfig{} // Initialize if nil
				}
				bootstrap.ServerConfigs.KubeAPIServerURL = "http://localhost:8080"
				bootstrap.KubeAPIToken = "fake-token"
			},
			expectError:    false,
			expectedClient: true,
		},
		{
			name: "WithExistingClientCache",
			setup: func() {
				// Set a fake clientset in the ClientCache
				mutating.ClientCache = fake.NewSimpleClientset()
			},
			expectError:    false,
			expectedClient: true,
		},
	}

	// Execute each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Run the setup function for each test case
			tc.setup()

			// Ensure to cleanup resources like temporary files after each test
			if tc.name == "WithKubeConfig" {
				defer os.Remove("test_kubeconfig")
				defer os.Unsetenv("KUBECONFIG")
			}

			// Call CreateClientSet
			err := mutating.CreateClientSet()

			// Assert on error expectation
			if tc.expectError {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Did not expect error but got one")
			}

			// Assert on ClientCache expectation
			if tc.expectedClient {
				assert.NotNil(t, mutating.ClientCache, "ClientCache should be initialized")
			} else {
				assert.Nil(t, mutating.ClientCache, "ClientCache should not be initialized")
			}

			// Reset ClientCache after each test
			mutating.ClientCache = nil
		})
	}
}
