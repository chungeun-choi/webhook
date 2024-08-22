package mutating

import (
	"context"
	"fmt"
	"github.com/chungeun-choi/webhook/bootstrap/kubernetes"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	admissionregistration "k8s.io/client-go/kubernetes/typed/admissionregistration/v1"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
)

type MutatingConfig struct {
	URL              string
	Client           kubernetes.ClientInterface
	AdmissionVersion []string
	FailurePolicy    admissionregistrationv1.FailurePolicyType
	CAPath           string
}

// MutatingManager is a struct that contains the client and the single flight.Group
type MutatingManager struct {
	Config             *MutatingConfig                                        // MutatingConfig is a struct that contains the URL, client, admission version, and failure policy
	admissionInitGroup singleflight.Group                                     // single flight.Group is a struct that provides a duplicate function call suppression
	admissionV1Client  admissionregistration.AdmissionregistrationV1Interface // Admission registrationV1Interface is an interface that contains the MutatingWebhookConfigurations method
	once               sync.Once                                              // once is a struct that provides a mechanism for performing exactly one action
	CAByte             []byte
}

// NewMutateManager is a function that creates a new instance of MutatingManager
func NewMutateManager(config *MutatingConfig) *MutatingManager {
	// If CAPath is empty, use the default CA
	if config.CAPath == "" {
		log.Printf("CAPath is empty, using the default CA")
		return &MutatingManager{
			Config: config,
		}
	}

	// Try to read the CA file from the specified path
	caBytes, err := os.ReadFile(config.CAPath)
	if err != nil {
		log.Printf("Failed to read CA file from %s: %v", config.CAPath, err)
		return &MutatingManager{
			Config: config,
		}
	}

	return &MutatingManager{
		Config: config,
		CAByte: caBytes,
	}
}

// Register is a method that registers the mutating webhook
func (m *MutatingManager) Register(req RequestAddRulesBody) (*ConfigBuilder, error) {
	req.Name = strings.ToLower(req.Name)

	// NewMutatingConfigBuilder is a function that creates a new instance of ConfigBuilder
	mutatingConfig := NewMutatingConfigBuilder().WithMetaInfo(req.Name)
	mutatingConfig.WithWebhook(NewWebhookConfigBuilder().
		WithName(req.Name+".admission"+".webhook").                  // required
		WithSideEffect(admissionregistrationv1.SideEffectClassNone). // required
		WithAdmissionReviewVersions(m.Config.AdmissionVersion...).   // required
		WithClientConfig(m.Config.URL, req.Name, m.CAByte).          // required
		WithRoles(req.Rules...).WithFailurePolicy(m.Config.FailurePolicy),
	)

	// getOldConfig is a method that retrieves the old configuration for the mutating webhook
	if old, err := m.Get(mutatingConfig.Name); err != nil {
		if apierrors.IsNotFound(err) {
			return m.create(*mutatingConfig)
		} else {
			return nil, errors.Wrap(err, "failed to get old config")
		}
	} else {
		return m.update(mutatingConfig, old)
	}
}

// Create is a method that creates a new configuration for the mutating webhook
func (m *MutatingManager) create(new ConfigBuilder) (*ConfigBuilder, error) {
	v1, err := m.GetAdmissionV1()
	if err != nil {
		return nil, err
	}

	newConfig, err := v1.MutatingWebhookConfigurations().Create(context.TODO(), &new.MutatingWebhookConfiguration, meta.CreateOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new config.")
	}

	return &ConfigBuilder{
		*newConfig,
	}, nil
}

// Update is a method that updates the configuration for the mutating webhook
func (m *MutatingManager) update(new, old *ConfigBuilder) (*ConfigBuilder, error) {
	if equalConfig(new, old) {
		log.Printf(" no need to update the configuration for the mutating webhook %s", old.Name)
		return old, nil
	} else {
		v1, err := m.GetAdmissionV1()
		if err != nil {
			return nil, err
		}

		new.MutatingWebhookConfiguration.ResourceVersion = old.ResourceVersion

		result, err := v1.MutatingWebhookConfigurations().Update(context.TODO(), &new.MutatingWebhookConfiguration, meta.UpdateOptions{})
		if err != nil {
			return nil, errors.Wrap(err, "failed to update the configuration for the mutating webhook")
		} else {
			return &ConfigBuilder{
				*result,
			}, nil
		}

	}
}

// Delete is a method that deletes the configuration for the mutating webhook
func (m *MutatingManager) Delete(name string) error {
	v1, err := m.GetAdmissionV1()
	if err != nil {
		return err
	}

	return v1.MutatingWebhookConfigurations().Delete(context.TODO(), name, meta.DeleteOptions{})
}

// Get is a method that retrieves the old configuration for the mutating webhook
func (m *MutatingManager) Get(name string) (*ConfigBuilder, error) {
	v1, err := m.GetAdmissionV1()
	if err != nil {
		return nil, err
	}

	oldConfig, err := v1.MutatingWebhookConfigurations().Get(context.TODO(), name, meta.GetOptions{})
	if err != nil {
		return nil, err
	}

	return &ConfigBuilder{
		*oldConfig,
	}, nil
}

func equalConfig(cur, old *ConfigBuilder) bool {
	// Use reflect.DeepEqual for deep comparison
	return reflect.DeepEqual(cur, old)
}

// GetAdmissionV1 is a method that returns the AdmissionregistrationV1Interface
func (m *MutatingManager) GetAdmissionV1() (admissionregistration.AdmissionregistrationV1Interface, error) {
	m.once.Do(func() {
		// Do is a method that executes and returns the result of the function f.
		v, err, _ := m.admissionInitGroup.Do("admissionV1", func() (interface{}, error) {
			return m.Config.Client.AdmissionregistrationV1(), nil
		})
		if err == nil {
			m.admissionV1Client = v.(admissionregistration.AdmissionregistrationV1Interface)
		}
	})
	if m.admissionV1Client == nil {
		return nil, fmt.Errorf("failed to initialize AdmissionregistrationV1 client")
	}
	return m.admissionV1Client, nil
}
