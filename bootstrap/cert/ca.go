package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/chungeun-choi/webhook/common/errors"
	"math/big"
	"time"
)

type CACertificateManager struct {
	Info *Info
}

type CACert struct {
	CAInfo *x509.Certificate
	Bytes  *bytes.Buffer
	Key    *rsa.PrivateKey
}

// NewCACert creates a new CA cert
func NewCACert(info *Info) *CACertificateManager {
	return &CACertificateManager{
		Info: info,
	}
}

// Generate generates the CA cert - if it is self-signed, it generates a new one, otherwise it loads the existing one
func (ca *CACertificateManager) Get() (*CACert, error) {
	var (
		result *CACert
		key    *rsa.PrivateKey
		err    error
	)

	// if the CA is self-signed, generate a new onex
	if ca.Info.CaType == SelfSigned {
		if key, err = rsa.GenerateKey(rand.Reader, 4096); err != nil {
			return nil, errors.FailedGenerateCACert(err)
		}

		if result, err = ca.generate(key); err != nil {
			return nil, errors.FailedGenerateCACert(err)
		} else {
			return result, nil
		}
		// otherwise, load the existing one
	} else {
		if result, err = ca.load(); err != nil {
			return nil, errors.FailedLoadCACert(err)
		} else {
			return result, nil
		}
	}
}

// Load loads the CA cert
func (ca *CACertificateManager) load() (*CACert, error) {
	// TODO: implement loading the CA cert
	return nil, nil
}

// Generate generates a new CA cert
func (ca *CACertificateManager) generate(key *rsa.PrivateKey) (*CACert, error) {
	info := &x509.Certificate{
		SerialNumber:          big.NewInt(2022),
		Subject:               pkix.Name{Organization: ca.Info.Org},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // expired in 1 year
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	bytesData, err := x509.CreateCertificate(rand.Reader, info, info, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}

	if pemBytes, err := EncodeToPem(bytesData, "CERTIFICATE"); err != nil {
		return nil, err
	} else {
		return &CACert{
			CAInfo: info,
			Bytes:  pemBytes,
			Key:    key,
		}, nil
	}
}
