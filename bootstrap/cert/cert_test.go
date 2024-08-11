package cert_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/chungeun-choi/webhook/bootstrap/cert"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func generateTestCACert() *cert.CACert {
	key, _ := rsa.GenerateKey(rand.Reader, 4096)
	info := &x509.Certificate{
		SerialNumber:          big.NewInt(2022),
		Subject:               pkix.Name{Organization: []string{"Test CA Organization"}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	caBytes, _ := x509.CreateCertificate(rand.Reader, info, info, &key.PublicKey, key)
	pemBytes, _ := cert.EncodeToPem(caBytes, "CERTIFICATE")

	return &cert.CACert{
		CAInfo: info,
		Bytes:  pemBytes,
		Key:    key,
	}
}

func TestGenerateCert(t *testing.T) {
	caCert := generateTestCACert()
	info := &cert.Info{
		Commons:  "Test Common Name",
		Org:      []string{"Test Organization"},
		DNSNames: []string{"test.example.com"},
	}

	manager := cert.NewCert(caCert.Key, info, caCert)
	certificate, err := manager.Generate()
	assert.NoError(t, err)
	assert.NotNil(t, certificate)

	// Check the generated certificate
	assert.Contains(t, certificate.CA.String(), "CERTIFICATE")
	assert.Contains(t, certificate.Cert.String(), "CERTIFICATE")
	assert.Contains(t, certificate.PrivateKey.String(), "RSA")

	// Decode the PEM block to get the DER-encoded certificate
	block, _ := pem.Decode(certificate.Cert.Bytes())
	if block == nil {
		t.Fatal("failed to decode PEM block containing the certificate")
	}

	// Parse the DER-encoded certificate
	parsedCert, err := x509.ParseCertificate(block.Bytes)
	assert.NoError(t, err)
	assert.Equal(t, "Test Common Name", parsedCert.Subject.CommonName)
	assert.Equal(t, []string{"Test Organization"}, parsedCert.Subject.Organization)
	assert.Contains(t, parsedCert.DNSNames, "test.example.com")
	assert.WithinDuration(t, parsedCert.NotBefore, time.Now(), time.Minute)
	assert.WithinDuration(t, parsedCert.NotAfter, time.Now().AddDate(1, 0, 0), time.Minute)
	assert.Contains(t, parsedCert.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
	assert.Contains(t, parsedCert.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	assert.Equal(t, parsedCert.KeyUsage, x509.KeyUsageDigitalSignature)
}
