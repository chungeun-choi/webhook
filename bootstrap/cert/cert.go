package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"github.com/chungeun-choi/webhook/common/errors"
	"math/big"
	"time"
)

type CertificateManager struct {
	key  *rsa.PrivateKey
	ca   *CACert
	info *Info
}

func NewCert(key *rsa.PrivateKey, info *Info, ca *CACert) *CertificateManager {
	return &CertificateManager{
		key:  key,
		info: info,
		ca:   ca,
	}
}

// Generate generates the cert
func (c *CertificateManager) Generate() (*Certificate, error) {
	var (
		result *Certificate = &Certificate{
			CA: c.ca.Bytes,
		}
		err error
	)

	newPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, errors.FailedGenerateCertPEM(err)
	}

	// new certificate config
	newCert := &x509.Certificate{
		DNSNames:     c.info.DNSNames,
		SerialNumber: big.NewInt(1024),
		Subject: pkix.Name{
			CommonName:   c.info.Commons,
			Organization: c.info.Org,
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(1, 0, 0), // expired in 1 year
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	// sign the new certificate
	newCertBytes, err := x509.CreateCertificate(rand.Reader, newCert, c.ca.CAInfo, &newPrivateKey.PublicKey, c.ca.Key)
	if err != nil {
		return nil, errors.FailedGenerateCertPEM(err)
	}

	// new certificate with PEM encoded
	if result.Cert, err = EncodeToPem(newCertBytes, "CERTIFICATE"); err != nil {
		return nil, errors.FailedGenerateCertPEM(err)
	}

	if result.PrivateKey, err = EncodeToPem(x509.MarshalPKCS1PrivateKey(newPrivateKey), "RSA"); err != nil {
		return nil, errors.FailedGenerateCertPEM(err)
	}

	return result, nil
}
