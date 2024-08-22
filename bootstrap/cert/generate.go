package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"
)

const DefaultCertFilePath = "cert"

type CertPathInfo struct {
	CaCertPath  string
	CertPath    string
	CertKeyPath string
}

func GenerateCert(orgs, dnsNames []string, commonName string) (*CertPathInfo, error) {
	// init CA config
	ca := &x509.Certificate{
		SerialNumber:          big.NewInt(2022),
		Subject:               pkix.Name{Organization: orgs},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // expired in 1 year
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// generate private key for CA
	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	// create the CA certificate
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		return nil, err
	}

	// CA certificate with PEM encoded
	caPEM := new(bytes.Buffer)
	_ = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	// new certificate config
	newCert := &x509.Certificate{
		DNSNames:     dnsNames,
		SerialNumber: big.NewInt(1024),
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: orgs,
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(1, 0, 0), // expired in 1 year
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	// generate new private key
	newPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	// sign the new certificate
	newCertBytes, err := x509.CreateCertificate(rand.Reader, newCert, ca, &newPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		return nil, err
	}

	// new certificate with PEM encoded
	newCertPEM := new(bytes.Buffer)
	_ = pem.Encode(newCertPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: newCertBytes,
	})

	// new private key with PEM encoded
	newPrivateKeyPEM := new(bytes.Buffer)
	_ = pem.Encode(newPrivateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(newPrivateKey),
	})

	// 파일로 저장하는 기능 추가
	err = saveToFile(DefaultCertFilePath+"/ca.pem", caPEM)
	if err != nil {
		return nil, err
	}

	err = saveToFile(DefaultCertFilePath+"/cert.pem", newCertPEM)
	if err != nil {
		return nil, err
	}

	err = saveToFile(DefaultCertFilePath+"/key.pem", newPrivateKeyPEM)
	if err != nil {
		return nil, err
	}

	return &CertPathInfo{
		CaCertPath:  DefaultCertFilePath + "/ca.pem",
		CertPath:    DefaultCertFilePath + "/cert.pem",
		CertKeyPath: DefaultCertFilePath + "/key.pem",
	}, nil
}

func saveToFile(filename string, data *bytes.Buffer) error {
	// 디렉토리가 존재하지 않으면 생성
	err := os.MkdirAll("cert", 0755)
	if err != nil {
		return err
	}

	// 파일 생성 및 데이터 저장
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data.Bytes())
	return err
}
