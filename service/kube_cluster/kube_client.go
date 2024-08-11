package kube_cluster

import (
	"fmt"
	"os"

	"github.com/chungeun-choi/webhook/bootstrap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ClientInterface abstracts the methods of the Kubernetes client interface.
type ClientInterface interface {
	kubernetes.Interface // Embeds the core methods of the Kubernetes client interface.
	// Add additional methods here if you need to mock specific clientset behavior.
}

var ClientCache ClientInterface

// CreateClientSet initializes the ClientCache if it is not already initialized.
func CreateClientSet() error {
	if ClientCache != nil {
		return nil
	}

	clientSet, err := initializeClientSet()
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes clientset: %w", err)
	}

	ClientCache = clientSet
	return nil
}

// initializeClientSet creates a Kubernetes clientset based on the available configuration.
func initializeClientSet() (ClientInterface, error) {
	kubeConfig := os.Getenv("KUBECONFIG")

	// Try to build the config from the KUBECONFIG environment variable.
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		// Fall back to in-cluster config or default REST config.
		config = createRestConfig()
	}

	// Create the clientset using the configuration.
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset from config: %w", err)
	}

	return clientSet, nil
}

// createRestConfig creates a REST configuration using bootstrap package settings.
func createRestConfig() *rest.Config {
	return &rest.Config{
		Host:        bootstrap.ServerConfigs.KubeAPIServerURL,
		BearerToken: bootstrap.KubeAPIToken,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true, // Consider setting this to false in production environments.
		},
	}
}
