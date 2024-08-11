package mutating_test

import (
	"testing"

	"github.com/chungeun-choi/webhook/bootstrap"
	"github.com/chungeun-choi/webhook/bootstrap/config"
	"github.com/chungeun-choi/webhook/service/mutating"
	"github.com/stretchr/testify/assert"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
)

// Test case struct for ConfigBuilder tests
type configBuilderTestCase struct {
	name          string
	setup         func() mutating.ConfigBuilder
	expectedName  string
	expectedCount int
}

// Test case struct for WebhookConfigBuilder tests
type webhookConfigTestCase struct {
	name               string
	setup              func() mutating.WebhookConfigBuilder
	expectedCA         []byte
	expectedService    *admissionregistrationv1.ServiceReference
	expectedURL        *string
	expectedOperations []admissionregistrationv1.OperationType
}

// TestConfigBuilder tests different scenarios for ConfigBuilder using test cases
func TestConfigBuilder(t *testing.T) {
	// Define the test cases
	testCases := []configBuilderTestCase{
		{
			name: "WithMetaInfo",
			setup: func() mutating.ConfigBuilder {
				builder := mutating.NewMutatingConfigBuilder()
				builder.WithMetaInfo("test-webhook")
				return *builder
			},
			expectedName:  "test-webhook",
			expectedCount: 0,
		},
		{
			name: "WithWebhook",
			setup: func() mutating.ConfigBuilder {
				builder := mutating.NewMutatingConfigBuilder()
				builder.WithMetaInfo("test-webhook") // Ensure the name is set
				webhookBuilder := mutating.WebhookConfigBuilder{}
				webhookBuilder.WithName("test-webhook")
				builder.WithWebhook(webhookBuilder)
				return *builder
			},
			expectedName:  "test-webhook",
			expectedCount: 1,
		},
	}

	// Execute each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := tc.setup()

			assert.Equal(t, tc.expectedName, builder.ObjectMeta.Name, "Expected name does not match")
			assert.Equal(t, tc.expectedCount, len(builder.Webhooks), "Webhook count does not match")
			if tc.expectedCount > 0 {
				assert.Equal(t, tc.expectedName, builder.Webhooks[0].Name, "Webhook name does not match")
			}
		})
	}
}

// TestWebhookConfigBuilder_WithClientConfig tests WebhookConfigBuilder's client config setup
func TestWebhookConfigBuilder_WithClientConfig(t *testing.T) {
	// Define the test cases
	testCases := []webhookConfigTestCase{
		{
			name: "RunningInK8S",
			setup: func() mutating.WebhookConfigBuilder {
				bootstrap.IsRunningInK8S = true
				bootstrap.ServerConfigs = new(config.ServerConfig)
				bootstrap.ServerConfigs.ServiceName = "test-service"
				bootstrap.Namespace = "default"

				webhookBuilder := mutating.WebhookConfigBuilder{}
				webhookBuilder.WithClientConfig("/test", []byte("mock-ca"))
				return webhookBuilder
			},
			expectedCA: []byte("mock-ca"),
			expectedService: &admissionregistrationv1.ServiceReference{
				Name:      "test-service",
				Namespace: "default",
				Path:      ptr("/test"),
			},
			expectedURL: nil,
		},
		{
			name: "RunningOutsideK8S",
			setup: func() mutating.WebhookConfigBuilder {
				bootstrap.IsRunningInK8S = false
				bootstrap.ServerConfigs.Hostname = "localhost"
				bootstrap.ServerConfigs.Port = 8080

				webhookBuilder := mutating.WebhookConfigBuilder{}
				webhookBuilder.WithClientConfig("/test", []byte("mock-ca"))
				return webhookBuilder
			},
			expectedCA:      []byte("mock-ca"),
			expectedService: nil,
			expectedURL:     ptr("https://localhost:8080/test"),
		},
	}

	// Execute each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			webhookBuilder := tc.setup()

			assert.Equal(t, tc.expectedCA, webhookBuilder.ClientConfig.CABundle)
			assert.Equal(t, tc.expectedService, webhookBuilder.ClientConfig.Service)
			assert.Equal(t, tc.expectedURL, webhookBuilder.ClientConfig.URL)
		})
	}
}

// TestWebhookConfigBuilder_WithRoles tests WebhookConfigBuilder's role setup
func TestWebhookConfigBuilder_WithRoles(t *testing.T) {
	// Define a simple test case
	testCases := []webhookConfigTestCase{
		{
			name: "WithRoles",
			setup: func() mutating.WebhookConfigBuilder {
				webhookBuilder := mutating.WebhookConfigBuilder{}
				rules := []mutating.Rule{
					{
						APIGroups:   []string{"apps"},
						APIVersions: []string{"v1"},
						Resources:   []string{"deployments"},
						Operations:  []string{"CREATE", "UPDATE"},
					},
				}
				webhookBuilder.WithRoles(rules...)
				return webhookBuilder
			},
			expectedOperations: []admissionregistrationv1.OperationType{
				admissionregistrationv1.Create,
				admissionregistrationv1.Update,
			},
		},
	}

	// Execute each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			webhookBuilder := tc.setup()

			assert.Equal(t, 1, len(webhookBuilder.Rules))
			assert.Equal(t, []string{"apps"}, webhookBuilder.Rules[0].APIGroups)
			assert.Equal(t, []string{"v1"}, webhookBuilder.Rules[0].APIVersions)
			assert.Equal(t, []string{"deployments"}, webhookBuilder.Rules[0].Resources)
			assert.Equal(t, tc.expectedOperations, webhookBuilder.Rules[0].Operations)
		})
	}
}

// Helper function to get a pointer to a string
func ptr(s string) *string {
	return &s
}
