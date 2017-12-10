/*
 * https://tools.ietf.org/html/rfc5280
 */

package pki

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/gob"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	pemKey  = "RSA PRIVATE KEY"
	pemCert = "CERTIFICATE"
	pemCSR  = "CERTIFICATE REQUEST"
)

type (
	// Manager provides methods for managing certificates
	Manager interface {
		GenCert(string, Config, Claim, Metadata) error
		GenCSR(string, Config, Claim, Metadata) error
		SignCSR(string, string, Policy, Metadata) error
	}

	// x509Manager provides methods for managing x509 certificates
	x509Manager struct{}
)

// NewX509Manager creates a new X509Manager
func NewX509Manager() Manager {
	return &x509Manager{}
}

func checkName(name string) error {
	if name == "" {
		return errors.New("Name is not set")
	}

	pattern := "./*/" + name + ".*"
	files, _ := filepath.Glob(pattern) // Glob ignores file system errors
	if len(files) > 0 {
		return errors.New(name + " already exists")
	}

	return nil
}

// genKeyPair generates a new public-private key pair
func genKeyPair(length int) (*rsa.PublicKey, *rsa.PrivateKey, error) {
	private, err := rsa.GenerateKey(rand.Reader, length)
	if err != nil {
		return nil, nil, err
	}
	public := &private.PublicKey

	return public, private, nil
}

func writePrivateKey(private *rsa.PrivateKey, password, path string) (err error) {
	var keyPem *pem.Block
	keyData := x509.MarshalPKCS1PrivateKey(private)

	// Encrypt private key if a password set
	if password == "" {
		keyPem = &pem.Block{
			Type:  pemKey,
			Bytes: keyData,
		}
	} else {
		keyPem, err = x509.EncryptPEMBlock(rand.Reader, pemKey, keyData, []byte(password), x509.PEMCipherAES256)
		if err != nil {
			return err
		}
	}

	keyFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	err = pem.Encode(keyFile, keyPem)
	if err != nil {
		return err
	}

	err = keyFile.Close()
	if err != nil {
		return err
	}

	return nil
}

func writePemFile(pemType string, pemData []byte, path string) error {
	pemBlock := &pem.Block{
		Type:  pemType,
		Bytes: pemData,
	}

	pemFile, err := os.Create(path)
	if err != nil {
		return err
	}

	err = pem.Encode(pemFile, pemBlock)
	if err != nil {
		return err
	}

	err = pemFile.Close()
	if err != nil {
		return err
	}

	return nil
}

// GenCert generates a new certificate
func (m *x509Manager) GenCert(name string, config Config, claim Claim, md Metadata) error {
	if err := checkName(name); err != nil {
		return err
	}

	config.Serial++
	length := config.Length
	startTime := time.Now()
	endTime := startTime.AddDate(0, 0, config.Days)

	// Generate a new public-private key pair
	publicKey, privateKey, err := genKeyPair(length)
	if err != nil {
		return err
	}

	// Read public key to compute Subject Key Identifier
	hash := sha1.New()
	pubKeyEncoder := gob.NewEncoder(hash)
	err = pubKeyEncoder.Encode(publicKey)
	if err != nil {
		return err
	}
	subjectKeyID := hash.Sum(nil)

	// Declare certificate
	rootCA := &x509.Certificate{
		SerialNumber: big.NewInt(config.Serial),

		NotBefore: startTime,
		NotAfter:  endTime,

		Subject: pkix.Name{
			CommonName:         claim.CommonName,
			Country:            claim.Country,
			Province:           claim.Province,
			Locality:           claim.Locality,
			Organization:       claim.Organization,
			OrganizationalUnit: claim.OrganizationalUnit,
			StreetAddress:      claim.StreetAddress,
			PostalCode:         claim.PostalCode,
		},

		// DNSNames:    []string{},
		// IPAddresses: []net.IP{},
		EmailAddresses: claim.EmailAddress,

		BasicConstraintsValid: true,
		IsCA: true,

		SubjectKeyId:   subjectKeyID,
		AuthorityKeyId: subjectKeyID,

		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage: []x509.ExtKeyUsage{},

		// Extensions:      []pkix.Extension{},
		// ExtraExtensions: []pkix.Extension{},
	}

	// Create the certificate
	cert, err := x509.CreateCertificate(rand.Reader, rootCA, rootCA, publicKey, privateKey)
	if err != nil {
		return err
	}

	/* Write certificate key file */
	keyFilePath := path.Join(DirRoot, name+extCAKey)
	err = writePrivateKey(privateKey, config.Password, keyFilePath)
	if err != nil {
		return err
	}

	/* Write certificate file */
	certFilePath := path.Join(DirRoot, name+extCACert)
	err = writePemFile(pemCert, cert, certFilePath)
	if err != nil {
		return err
	}

	return nil
}

// GenCSR generates a certificate signing request
func (m *x509Manager) GenCSR(name string, config Config, claim Claim, md Metadata) error {
	if err := checkName(name); err != nil {
		return err
	}

	config.Serial++
	length := config.Length
	// TODO startTime := time.Now()
	// TODO endTime := startTime.AddDate(0, 0, config.Days)

	// Generate a new public-private key pair
	_, privateKey, err := genKeyPair(length) // TODO
	if err != nil {
		return err
	}

	// Declare certificate request
	intermCSR := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:         claim.CommonName,
			Country:            claim.Country,
			Province:           claim.Province,
			Locality:           claim.Locality,
			Organization:       claim.Organization,
			OrganizationalUnit: claim.OrganizationalUnit,
			StreetAddress:      claim.StreetAddress,
			PostalCode:         claim.PostalCode,
		},

		// DNSNames:    []string{},
		// IPAddresses: []net.IP{},
		EmailAddresses: claim.EmailAddress,

		// Extensions:      []pkix.Extension{},
		// ExtraExtensions: []pkix.Extension{},
	}

	// Create the certificate request
	csr, err := x509.CreateCertificateRequest(rand.Reader, intermCSR, privateKey)
	if err != nil {
		return err
	}

	/* Write certificate key file */
	keyFilePath := path.Join(md.Dir(), name+extCAKey)
	err = writePrivateKey(privateKey, config.Password, keyFilePath)
	if err != nil {
		return err
	}

	/* Write certificate request file */
	csrFilePath := path.Join(DirCSR, name+extCACSR)
	err = writePemFile(pemCSR, csr, csrFilePath)
	if err != nil {
		return err
	}

	return nil
}

// SignCSR signs a certificate signing request using a certificate authority
func (m *x509Manager) SignCSR(nameCA, nameCSR string, policy Policy, md Metadata) error {
	return nil
}
