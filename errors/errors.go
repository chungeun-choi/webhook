package errors

import (
	"github.com/pkg/errors"
)

var (
	// cert file
	ErrFailedLoadCACert      = errors.New("failed to load CA cert file with by using path")
	ErrFailedGenerateCACert  = errors.New("failed to generate CA cert")
	ERrFailedGenerateCertPEM = errors.New("failed to generate cert pem")
	ErrKubeConfigNotFound    = errors.New("kubeconfig not found")
	ErrFailedCreateClientSet = errors.New("failed to create clientset")
	ErrNotFound              = errors.New("not found")
)

func FailedLoadCACert(err error) error {
	return errors.Wrapf(ErrFailedLoadCACert, "failed to load CA cert file with by using path: %v", err)
}

func FailedGenerateCACert(err error) error {
	return errors.Wrapf(ErrFailedGenerateCACert, "failed to generate CA cert: %v", err)
}

func FailedGenerateCertPEM(err error) error {
	return errors.Wrapf(ERrFailedGenerateCertPEM, "failed to generate cert pem: %v", err)
}
