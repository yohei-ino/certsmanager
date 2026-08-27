package pki

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// --- Constants and Structures ---

const caKeyFile = "ca.key"
const caCertFile = "ca.crt"
const serverDir = "servers"

// InfoCA holds the data of a functional CA.
type InfoCA struct {
	CommonName string
}

// InitializeCA generates the private key and certificate for the Certificate Authority (CA).
func InitializeCA(projectName string, projectDir string) error {
	fmt.Printf("-> Configuring directory: %s\n", projectDir)
	if err := os.Remove(caKeyFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to clean up old keys: %v", err)
	}
	if err := os.MkdirAll(projectDir, 0700); err != nil {
		return fmt.Errorf("failed to create PKI directory: %v", err)
	}

	// 1. Generate RSA private key for the CA
	fmt.Println("-> Generating RSA private key for CA... (2048 bits)")
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate CA key: %v", err)
	}

	// 2. Create PEM block for the private key and save it
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey)})
	if err := os.WriteFile(filepath.Join(projectDir, caKeyFile), privPEM, 0600); err != nil {
		return fmt.Errorf("failed to save CA private key (%s): %v", caKeyFile, err)
	}
	fmt.Printf("   [OK] CA private key saved in %s/%s (Keep it secure!)\n", projectName, caKeyFile)

	// 3. Create the CA certificate
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Root CA"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // Valid for 10 years
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certDERBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return fmt.Errorf("failed to create self-signed certificate: %v", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDERBytes})
	if err := os.WriteFile(filepath.Join(projectDir, caCertFile), certPEM, 0644); err != nil {
		return fmt.Errorf("failed to save CA certificate (%s): %v", caCertFile, err)
	}
	fmt.Printf("   [OK] CA certificate saved in %s/%s\n", projectName, caCertFile)

	return nil
}

// AddServerCertificate generates a new key pair and a certificate signed by the CA.
func AddServerCertificate(projectName string, serverName string, sansList string, projectDir string) error {
	// 1. Generate private key for the new server
	fmt.Println("-> Generating private key for the server...")
	serverPrivKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate server key: %v", err)
	}

	// 2. Define the saving path
	serverPath := filepath.Join(projectDir, serverDir)
	if err := os.MkdirAll(serverPath, 0700); err != nil {
		return fmt.Errorf("failed to create servers directory: %v", err)
	}

	keyFilePath := filepath.Join(serverPath, serverName+"_server.key")
	certFilePath := filepath.Join(serverPath, serverName+"_server.crt")

	// 3. Save server private key
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(serverPrivKey)})
	if err := os.WriteFile(keyFilePath, privPEM, 0600); err != nil {
		return fmt.Errorf("failed to save server private key: %v", err)
	}
	fmt.Printf("   [OK] Server private key saved in %s\n", keyFilePath)

	// 4. Configure certificate information for the server
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: serverName},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // Valid for 1 year
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Process SANs
	if sansList != "" {
		sansItems := strings.Split(sansList, ",")
		for _, san := range sansItems {
			san = strings.TrimSpace(san)
			if san != "" {
				template.DNSNames = append(template.DNSNames, san)
				fmt.Printf("      [INFO] SAN '%s' added to DNSNames.\n", san)
			}
		}
	}

	// 5. Read CA certificate and key for signing
	caCertBytes, err := os.ReadFile(filepath.Join(projectDir, caCertFile))
	if err != nil {
		return fmt.Errorf("unable to read CA certificate: %v", err)
	}

	caPrivKeyBytes, err := os.ReadFile(filepath.Join(projectDir, caKeyFile))
	if err != nil {
		return fmt.Errorf("unable to read CA private key: %v", err)
	}

	block, _ := pem.Decode(caPrivKeyBytes)
	if block == nil {
		return fmt.Errorf("failed to decode CA private key PEM")
	}

	caPrivKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse CA private key: %v", err)
	}

	certBlock, _ := pem.Decode(caCertBytes)
	if certBlock == nil {
		return fmt.Errorf("failed to decode CA certificate PEM")
	}
	caCertParsed, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse CA certificate: %v", err)
	}

	// Sign the server certificate using the CA
	certDERBytes, err := x509.CreateCertificate(rand.Reader, template, caCertParsed, &serverPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return fmt.Errorf("failed to sign server certificate with CA: %v", err)
	}

	// 6. Save the final certificate
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDERBytes})
	if err := os.WriteFile(certFilePath, certPEM, 0644); err != nil {
		return fmt.Errorf("failed to save server certificate: %v", err)
	}
	fmt.Printf("   [OK] Server certificate saved in %s\n", certFilePath)

	return nil
}

