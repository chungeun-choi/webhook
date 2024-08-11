package cert_test

import (
	"crypto/x509"
	"github.com/chungeun-choi/webhook/bootstrap/cert"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCACert(t *testing.T) {
	info := &cert.Info{
		Org:    []string{"Test Organization"},
		CaType: cert.SelfSigned,
	}

	caManager := cert.NewCACert(info)
	assert.NotNil(t, caManager)
	assert.Equal(t, caManager.Info, info)
}

func TestGenerateSelfSignedCACert(t *testing.T) {
	info := &cert.Info{
		Org:    []string{"Test Organization"},
		CaType: cert.SelfSigned,
	}

	caManager := cert.NewCACert(info)
	caCert, err := caManager.Get()
	assert.NoError(t, err)
	assert.NotNil(t, caCert)

	// Check the certificate info
	assert.Equal(t, caCert.CAInfo.Subject.Organization[0], "Test Organization")
	assert.True(t, caCert.CAInfo.IsCA)
	assert.True(t, caCert.CAInfo.BasicConstraintsValid)
	assert.Contains(t, caCert.CAInfo.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
	assert.Contains(t, caCert.CAInfo.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	assert.Equal(t, caCert.CAInfo.KeyUsage, x509.KeyUsageDigitalSignature|x509.KeyUsageCertSign)
	assert.WithinDuration(t, caCert.CAInfo.NotBefore, time.Now(), time.Minute)
	assert.WithinDuration(t, caCert.CAInfo.NotAfter, time.Now().AddDate(1, 0, 0), time.Minute)
}
