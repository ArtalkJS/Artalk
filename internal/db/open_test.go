package db

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadMySQLTLSConfig(t *testing.T) {
	certDir := t.TempDir()
	caCert, caKey := createTestCertificate(t, nil, nil, x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Test CA"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
	})
	serverCert, _ := createTestCertificate(t, caCert, caKey, x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: "db.example.internal"},
		DNSNames:     []string{"db.example.internal"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	})
	clientCert, clientKey := createTestCertificate(t, caCert, caKey, x509.Certificate{
		SerialNumber: big.NewInt(3),
		Subject:      pkix.Name{CommonName: "client"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	})

	caPath := writeTestCertificate(t, certDir, "server-ca.pem", caCert.Raw)
	clientCertPath := writeTestCertificate(t, certDir, "client-cert.pem", clientCert.Raw)
	clientKeyPath := writeTestPrivateKey(t, certDir, "client-key.pem", clientKey)

	tlsConfig, configured, err := loadMySQLTLSConfig(&config.DBConf{
		Host:           "db.example.internal",
		ServerCaPath:   caPath,
		ClientCertPath: clientCertPath,
		ClientKeyPath:  clientKeyPath,
	})
	require.NoError(t, err)
	assert.True(t, configured)
	assert.False(t, tlsConfig.InsecureSkipVerify)
	assert.Equal(t, "db.example.internal", tlsConfig.ServerName)
	assert.Equal(t, uint16(tls.VersionTLS12), tlsConfig.MinVersion)
	assert.Len(t, tlsConfig.Certificates, 1)

	_, err = serverCert.Verify(x509.VerifyOptions{
		DNSName:   tlsConfig.ServerName,
		Roots:     tlsConfig.RootCAs,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	})
	require.NoError(t, err)

	_, err = serverCert.Verify(x509.VerifyOptions{
		DNSName:   "attacker.example.internal",
		Roots:     tlsConfig.RootCAs,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	})
	assert.Error(t, err, "A certificate for a different host must be rejected")
}

func TestLoadMySQLTLSConfigRequiresCompleteCertificateSet(t *testing.T) {
	_, _, err := loadMySQLTLSConfig(&config.DBConf{ServerCaPath: "server-ca.pem"})
	assert.EqualError(t, err, "server CA, client certificate, and client key must be configured together")
}

func createTestCertificate(
	t *testing.T,
	parent *x509.Certificate,
	parentKey ed25519.PrivateKey,
	template x509.Certificate,
) (*x509.Certificate, ed25519.PrivateKey) {
	t.Helper()
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	if parent == nil {
		parent = &template
		parentKey = privateKey
	}
	der, err := x509.CreateCertificate(rand.Reader, &template, parent, publicKey, parentKey)
	require.NoError(t, err)
	cert, err := x509.ParseCertificate(der)
	require.NoError(t, err)
	return cert, privateKey
}

func writeTestCertificate(t *testing.T, dir string, name string, der []byte) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o600))
	return path
}

func writeTestPrivateKey(t *testing.T, dir string, name string, key ed25519.PrivateKey) string {
	t.Helper()
	der, err := x509.MarshalPKCS8PrivateKey(key)
	require.NoError(t, err)
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), 0o600))
	return path
}
