package mutating

import (
	"context"
	"fmt"
	"github.com/chungeun-choi/webhook/service/kube_cluster"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	admissionregistrationv1 "k8s.io/client-go/kubernetes/typed/admissionregistration/v1"
	"log"
	"reflect"
	"sync"
)

// MutatingManager is a struct that contains the client and the singleflight.Group
type MutatingManager struct { // name is a string that represents the name of the mutating webhook
	Client             kube_cluster.ClientInterface                             // Client is an interface that contains the AdmissionregistrationV1 method
	admissionInitGroup singleflight.Group                                       // singleflight.Group is a struct that provides a duplicate function call suppression
	once               sync.Once                                                // once is a struct that provides a mechanism for performing exactly one action
	admissionV1Client  admissionregistrationv1.AdmissionregistrationV1Interface // AdmissionregistrationV1Interface is an interface that contains the MutatingWebhookConfigurations method
}

// NewMutateManager is a function that creates a new instance of MutatingManager
func NewMutateManager(client kube_cluster.ClientInterface) *MutatingManager {
	return &MutatingManager{
		Client: client,
	}
}

// Register is a method that registers the mutating webhook
func (m *MutatingManager) Register(config *ConfigBuilder) (*ConfigBuilder, error) {
	if old, err := m.getOldConfig(config.Name); err != nil {
		if apierrors.IsNotFound(err) {
			return m.Create(*config)
		} else {
			return nil, errors.Wrap(err, "failed to get old config")
		}
	} else {
		return m.Update(config, old)
	}
}

// Create is a method that creates a new configuration for the mutating webhook
func (m *MutatingManager) Create(new ConfigBuilder) (*ConfigBuilder, error) {
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
func (m *MutatingManager) Update(new, old *ConfigBuilder) (*ConfigBuilder, error) {
	if equalConfig(new, old) {
		log.Printf(" no need to update the configuration for the mutating webhook %s", old.Name)
		return old, nil
	} else {
		v1, err := m.GetAdmissionV1()
		if err != nil {
			return nil, err
		}

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

// getOldConfig is a method that retrieves the old configuration for the mutating webhook
func (m *MutatingManager) getOldConfig(name string) (*ConfigBuilder, error) {
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
func (m *MutatingManager) GetAdmissionV1() (admissionregistrationv1.AdmissionregistrationV1Interface, error) {
	m.once.Do(func() {
		// Do is a method that executes and returns the result of the function f.
		v, err, _ := m.admissionInitGroup.Do("admissionV1", func() (interface{}, error) {
			return m.Client.AdmissionregistrationV1(), nil
		})
		if err == nil {
			m.admissionV1Client = v.(admissionregistrationv1.AdmissionregistrationV1Interface)
		}
	})
	if m.admissionV1Client == nil {
		return nil, fmt.Errorf("failed to initialize AdmissionregistrationV1 client")
	}
	return m.admissionV1Client, nil
}
