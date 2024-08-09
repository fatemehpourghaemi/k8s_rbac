package certificate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

func GenerateClientCreds(username string, caCertPath string, caKeyPath string) (string, string, string, error) {
	caCertBytes, err := os.ReadFile(caCertPath)
	if err != nil {
		return "", "", "", fmt.Errorf("Failed to read CA certificate file: %v", err)
	}

	caKeyBytes, err := os.ReadFile(caKeyPath)
	if err != nil {
		return "", "", "", fmt.Errorf("Failed to read CA private key file: %v", err)
	}

	caCertBlock, _ := pem.Decode(caCertBytes)
	if caCertBlock == nil {
		return "", "", "", fmt.Errorf("failed to decode CA certificate")
	}
	caKeyBlock, _ := pem.Decode(caKeyBytes)
	if caKeyBlock == nil {
		return "", "", "", fmt.Errorf("failed to decode CA private key")
	}

	caCertParsed, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to parse CA certificate: %v", err)
	}

	caKeyParsed, err := x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to parse CA private key: %v", err)
	}

	clientKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate client private key: %v", err)
	}

	csrTemplate := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: username,
		},
	}

	_, err = x509.CreateCertificateRequest(rand.Reader, &csrTemplate, clientKey)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create CSR: %v", err)
	}

	clientCert, err := x509.CreateCertificate(rand.Reader, &x509.Certificate{
		Subject:      csrTemplate.Subject,
		PublicKey:    clientKey.Public(),
		SerialNumber: big.NewInt(1658),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(2, 0, 0), // 2 years
	}, caCertParsed, &clientKey.PublicKey, caKeyParsed)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to sign CSR: %v", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: clientCert})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(clientKey)})

	caB64 := base64.StdEncoding.EncodeToString(caCertBytes)
	certB64 := base64.StdEncoding.EncodeToString(certPEM)
	keyB64 := base64.StdEncoding.EncodeToString(keyPEM)

	return caB64, certB64, keyB64, nil
}
