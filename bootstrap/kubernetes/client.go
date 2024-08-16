package kubernetes

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ClientInterface abstracts the methods of the Kubernetes client interface.
type ClientInterface interface {
	kubernetes.Interface // Embeds the core methods of the Kubernetes client interface.
	// Add additional methods here if you need to mock specific clientset behavior.
}

var clientCache ClientInterface

// CreateClientSet initializes the ClientCache if it is not already initialized.
func CreateClientSet(url, token string) (ClientInterface, error) {
	if clientCache != nil {
		return clientCache, nil
	}

	clientSet, err := initializeClientSet(url, token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes clientset: %w", err)
	}

	clientCache = clientSet
	return clientCache, nil
}

// initializeClientSet creates a Kubernetes clientset based on the available configuration.
func initializeClientSet(url, token string) (ClientInterface, error) {
	kubernetesConfig := createRestConfig(url, token)

	// Create the clientset using the configuration.
	clientSet, err := kubernetes.NewForConfig(kubernetesConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset from config: %w", err)
	}

	return clientSet, nil
}

// createRestConfig creates a REST configuration using bootstrap package settings.
func createRestConfig(url, token string) *rest.Config {
	//TODO: Implement the logic to read the certificate, key, and CA data from the files.

	return &rest.Config{
		Host:        url,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
}
