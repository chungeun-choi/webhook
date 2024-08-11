package mutating

import (
	"context"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
	"testing"
)

func strPtr(s string) *string {
	return &s
}

// setUp initializes the test environment with a fake clientset and MutatingManager
func setUp(t *testing.T, name string) (*fake.Clientset, *MutatingManager) {
	clientset := fake.NewSimpleClientset()
	mutateManager := NewMutateManager(clientset)

	// Add reactors for create requests
	clientset.Fake.PrependReactor("create", "mutatingwebhookconfigurations", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		createAction := action.(k8stesting.CreateAction)
		webhookConfig := createAction.GetObject().(*admissionregistrationv1.MutatingWebhookConfiguration)
		if webhookConfig.Name == "" {
			return true, nil, apierrors.NewInvalid(
				schema.GroupKind{Group: "admissionregistration.k8s.io", Kind: "MutatingWebhookConfiguration"},
				webhookConfig.Name,
				nil,
			)
		}

		// Add the configuration to the tracker's state
		err = clientset.Tracker().Add(webhookConfig)
		if err != nil {
			return true, nil, err
		}

		return true, webhookConfig, nil
	})

	// Add reactors for update requests
	clientset.Fake.PrependReactor("update", "mutatingwebhookconfigurations", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		updateAction := action.(k8stesting.UpdateAction)
		webhookConfig := updateAction.GetObject().(*admissionregistrationv1.MutatingWebhookConfiguration)
		if webhookConfig.Name == "non-existent-webhook" {
			return true, nil, apierrors.NewNotFound(
				schema.GroupResource{Group: "admissionregistration.k8s.io", Resource: "mutatingwebhookconfigurations"},
				webhookConfig.Name,
			)
		}

		// Update the configuration in the tracker's state
		err = clientset.Tracker().Update(admissionregistrationv1.SchemeGroupVersion.WithResource("mutatingwebhookconfigurations"), webhookConfig, "")
		if err != nil {
			return true, nil, err
		}

		return true, webhookConfig, nil
	})

	// Add reactors for get requests
	clientset.Fake.PrependReactor("get", "mutatingwebhookconfigurations", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		getAction := action.(k8stesting.GetAction)
		obj, err := clientset.Tracker().Get(admissionregistrationv1.SchemeGroupVersion.WithResource("mutatingwebhookconfigurations"), "", getAction.GetName())
		if err != nil {
			return true, nil, apierrors.NewNotFound(
				schema.GroupResource{Group: "admissionregistration.k8s.io", Resource: "mutatingwebhookconfigurations"},
				getAction.GetName(),
			)
		}
		return true, obj, nil
	})

	// Add reactors for delete requests
	clientset.Fake.PrependReactor("delete", "mutatingwebhookconfigurations", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		deleteAction := action.(k8stesting.DeleteAction)
		err = clientset.Tracker().Delete(admissionregistrationv1.SchemeGroupVersion.WithResource("mutatingwebhookconfigurations"), "", deleteAction.GetName())
		if err != nil {
			return true, nil, apierrors.NewNotFound(
				schema.GroupResource{Group: "admissionregistration.k8s.io", Resource: "mutatingwebhookconfigurations"},
				deleteAction.GetName(),
			)
		}
		return true, nil, nil
	})

	return clientset, mutateManager
}

// tearDown cleans up after each test
func tearDown(clientset *fake.Clientset) {
	clientset.ClearActions()
}

func TestCreate(t *testing.T) {
	clientset, mutateManager := setUp(t, "")
	defer tearDown(clientset)

	tests := []struct {
		name        string
		config      *ConfigBuilder
		expectError bool
	}{
		{
			name: "Create New Webhook",
			config: &ConfigBuilder{
				MutatingWebhookConfiguration: admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: meta.ObjectMeta{
						Name: "test-webhook",
					},
					Webhooks: []admissionregistrationv1.MutatingWebhook{
						{
							Name: "example.webhook",
							ClientConfig: admissionregistrationv1.WebhookClientConfig{
								Service: &admissionregistrationv1.ServiceReference{
									Name:      "webhook-service",
									Namespace: "default",
									Path:      strPtr("/mutate"),
								},
								CABundle: []byte("fake-ca-bundle"),
							},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Fail to Create Webhook - Invalid Config",
			config: &ConfigBuilder{
				MutatingWebhookConfiguration: admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: meta.ObjectMeta{
						Name: "", // Invalid because name is required
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdConfig, err := mutateManager.Create(*tt.config)
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got none")
				}
			} else {
				if err != nil {
					t.Fatalf("Create() error = %v", err)
				}
				if createdConfig.Name != tt.config.Name {
					t.Fatalf("expected webhook name to be %s, but got %s", tt.config.Name, createdConfig.Name)
				}
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	clientset, mutateManager := setUp(t, "")
	defer tearDown(clientset)

	tests := []struct {
		name        string
		oldConfig   *ConfigBuilder
		newConfig   *ConfigBuilder
		expectError bool
	}{
		{
			name: "Update Existing Webhook",
			oldConfig: &ConfigBuilder{
				MutatingWebhookConfiguration: admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: meta.ObjectMeta{
						Name: "test-webhook",
					},
					Webhooks: []admissionregistrationv1.MutatingWebhook{
						{
							Name: "example.webhook",
							ClientConfig: admissionregistrationv1.WebhookClientConfig{
								Service: &admissionregistrationv1.ServiceReference{
									Name:      "webhook-service",
									Namespace: "default",
									Path:      strPtr("/mutate"),
								},
								CABundle: []byte("fake-ca-bundle"),
							},
						},
					},
				},
			},
			newConfig: &ConfigBuilder{
				MutatingWebhookConfiguration: admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: meta.ObjectMeta{
						Name: "test-webhook",
					},
					Webhooks: []admissionregistrationv1.MutatingWebhook{
						{
							Name: "example.webhook",
							ClientConfig: admissionregistrationv1.WebhookClientConfig{
								Service: &admissionregistrationv1.ServiceReference{
									Name:      "webhook-service",
									Namespace: "default",
									Path:      strPtr("/mutate-v2"),
								},
								CABundle: []byte("fake-ca-bundle"),
							},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Fail to Update - Webhook Not Found",
			oldConfig: &ConfigBuilder{
				MutatingWebhookConfiguration: admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: meta.ObjectMeta{
						Name: "non-existent-webhook",
					},
				},
			},
			newConfig: &ConfigBuilder{
				MutatingWebhookConfiguration: admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: meta.ObjectMeta{
						Name: "non-existent-webhook",
					},
					Webhooks: []admissionregistrationv1.MutatingWebhook{
						{
							Name: "example.webhook",
							ClientConfig: admissionregistrationv1.WebhookClientConfig{
								Service: &admissionregistrationv1.ServiceReference{
									Name:      "webhook-service",
									Namespace: "default",
									Path:      strPtr("/mutate-v2"),
								},
								CABundle: []byte("fake-ca-bundle"),
							},
						},
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the old configuration if it exists
			if tt.oldConfig != nil && tt.oldConfig.Name != "non-existent-webhook" {
				_, _ = clientset.AdmissionregistrationV1().MutatingWebhookConfigurations().Create(context.TODO(), &tt.oldConfig.MutatingWebhookConfiguration, meta.CreateOptions{})
			}

			updatedConfig, err := mutateManager.Update(tt.newConfig, tt.oldConfig)
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got none")
				}
			} else {
				if err != nil {
					t.Fatalf("Update() error = %v", err)
				}
				if updatedConfig.Webhooks[0].ClientConfig.Service.Path == tt.oldConfig.Webhooks[0].ClientConfig.Service.Path {
					t.Fatalf("expected webhook path to be updated from %v to %v", tt.oldConfig.Webhooks[0].ClientConfig.Service.Path, updatedConfig.Webhooks[0].ClientConfig.Service.Path)
				}
			}
		})
	}
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name        string
		oldConfig   *ConfigBuilder
		newConfig   *ConfigBuilder
		expectError bool
	}{
		{
			name:      "Register New Webhook",
			oldConfig: nil,
			newConfig: &ConfigBuilder{
				MutatingWebhookConfiguration: admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: meta.ObjectMeta{
						Name: "test-webhook",
					},
					Webhooks: []admissionregistrationv1.MutatingWebhook{
						{
							Name: "example.webhook",
							ClientConfig: admissionregistrationv1.WebhookClientConfig{
								Service: &admissionregistrationv1.ServiceReference{
									Name:      "webhook-service",
									Namespace: "default",
									Path:      strPtr("/mutate"),
								},
								CABundle: []byte("fake-ca-bundle"),
							},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "Update Existing Webhook",
			oldConfig: &ConfigBuilder{
				MutatingWebhookConfiguration: admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: meta.ObjectMeta{
						Name: "test-webhook",
					},
					Webhooks: []admissionregistrationv1.MutatingWebhook{
						{
							Name: "example.webhook",
							ClientConfig: admissionregistrationv1.WebhookClientConfig{
								Service: &admissionregistrationv1.ServiceReference{
									Name:      "webhook-service",
									Namespace: "default",
									Path:      strPtr("/mutate"),
								},
								CABundle: []byte("fake-ca-bundle"),
							},
						},
					},
				},
			},
			newConfig: &ConfigBuilder{
				MutatingWebhookConfiguration: admissionregistrationv1.MutatingWebhookConfiguration{
					ObjectMeta: meta.ObjectMeta{
						Name: "test-webhook",
					},
					Webhooks: []admissionregistrationv1.MutatingWebhook{
						{
							Name: "example.webhook",
							ClientConfig: admissionregistrationv1.WebhookClientConfig{
								Service: &admissionregistrationv1.ServiceReference{
									Name:      "webhook-service",
									Namespace: "default",
									Path:      strPtr("/mutate-v2"),
								},
								CABundle: []byte("fake-ca-bundle"),
							},
						},
					},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientset, mutateManager := setUp(t, "test-webhook")
			defer tearDown(clientset)
			if tt.oldConfig != nil {
				_, err := clientset.AdmissionregistrationV1().MutatingWebhookConfigurations().Create(context.TODO(), &tt.oldConfig.MutatingWebhookConfiguration, meta.CreateOptions{})
				if err != nil {
					t.Fatalf("failed to create old config: %v", err)
				}
			}

			registeredConfig, err := mutateManager.Register(tt.newConfig)
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got none")
				}
			} else {
				if err != nil {
					t.Fatalf("Register() error = %v", err)
				}

				// Check if the webhook was registered correctly
				if tt.oldConfig != nil && *registeredConfig.Webhooks[0].ClientConfig.Service.Path != *tt.newConfig.Webhooks[0].ClientConfig.Service.Path {
					t.Fatalf("expected webhook path to be updated to %v, but got %v", *tt.newConfig.Webhooks[0].ClientConfig.Service.Path, *registeredConfig.Webhooks[0].ClientConfig.Service.Path)
				} else if tt.oldConfig == nil && registeredConfig.Name != tt.newConfig.Name {
					t.Fatalf("expected webhook name to be %v, but got %v", tt.newConfig.Name, registeredConfig.Name)
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	clientset, mutateManager := setUp(t, "mutating-webhook")
	defer tearDown(clientset)

	// Create a configuration to delete later
	_, err := clientset.AdmissionregistrationV1().MutatingWebhookConfigurations().Create(context.TODO(), &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: meta.ObjectMeta{
			Name: "mutating-webhook",
		},
	}, meta.CreateOptions{})

	if err != nil {
		t.Fatalf("failed to create configuration: %v", err)
	}

	// Test the Delete method
	err = mutateManager.Delete("mutating-webhook")
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Try to get the deleted configuration to ensure it's gone
	_, err = clientset.AdmissionregistrationV1().MutatingWebhookConfigurations().Get(context.TODO(), "mutating-webhook", meta.GetOptions{})
	if err == nil {
		t.Fatal("expected error, got none")
	}
	if !apierrors.IsNotFound(err) {
		t.Fatalf("expected not found error, got %v", err)
	}
}
