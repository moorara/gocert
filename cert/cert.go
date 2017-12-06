/*
 * https://tools.ietf.org/html/rfc5280
 */

package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/gob"
	"encoding/pem"
	"math/big"
	"os"
	"path"
	"time"

	"github.com/moorara/gocert/config"
)

const (
	fileRootCACert = "root.ca.cert.pem"
	fileRootCAKey  = "root.ca.key.pem"
)

type (
	// Manager provides methods for generating certificates
	Manager interface {
		GenRootCA(config.SettingsCA, config.Claim) error
		GenIntermCA(config.SettingsCA, config.Claim) error
		GenServerCert(config.Settings, config.Claim) error
		GenClientCert(config.Settings, config.Claim) error
	}

	// X509Manager provides methods for generating x509 certificates
	X509Manager struct {
	}
)

// NewX509Manager creates a new X509Manager
func NewX509Manager() *X509Manager {
	return &X509Manager{}
}

// genKeys generates a new public-private key pair
func genKeys(length int) (*rsa.PublicKey, *rsa.PrivateKey, error) {
	private, err := rsa.GenerateKey(rand.Reader, length)
	if err != nil {
		return nil, nil, err
	}
	public := &private.PublicKey

	return public, private, nil
}

func getPrivatePem(private *rsa.PrivateKey, password string) (keyPem *pem.Block, err error) {
	keyType := "RSA PRIVATE KEY"
	keyData := x509.MarshalPKCS1PrivateKey(private)

	// Encrypt private key if a password set
	if password == "" {
		keyPem = &pem.Block{
			Type:  keyType,
			Bytes: keyData,
		}
	} else {
		keyPem, err = x509.EncryptPEMBlock(rand.Reader, keyType, keyData, []byte(password), x509.PEMCipherAES256)
		if err != nil {
			return nil, err
		}
	}

	return keyPem, nil
}

// GenRootCA generates root certificate authority
func (m *X509Manager) GenRootCA(settings config.SettingsCA, claim config.Claim) error {
	settings.Serial++
	length := settings.Length
	startTime := time.Now()
	endTime := startTime.AddDate(0, 0, settings.Days)

	// Generate a new public-private key pair
	publicKey, privateKey, err := genKeys(length)
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
		SerialNumber: big.NewInt(settings.Serial),

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

	/* Write certificate file */

	certFilePath := path.Join(config.DirNameRoot, fileRootCACert)
	certFile, err := os.Create(certFilePath)
	if err != nil {
		return err
	}

	certPem := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	}

	err = pem.Encode(certFile, certPem)
	if err != nil {
		return err
	}

	err = certFile.Close()
	if err != nil {
		return err
	}

	/* Write certificate key file */

	keyFilePath := path.Join(config.DirNameRoot, fileRootCAKey)
	keyFile, err := os.OpenFile(keyFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	keyPem, err := getPrivatePem(privateKey, settings.Password)
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

// GenIntermCA generates an intermediate certificate authority
func (m *X509Manager) GenIntermCA(settings config.SettingsCA, claim config.Claim) error {
	return nil
}

// GenServerCert generates a Server certificate
func (m *X509Manager) GenServerCert(settings config.Settings, claim config.Claim) error {
	return nil
}

// GenClientCert generates a client certificate
func (m *X509Manager) GenClientCert(settings config.Settings, claim config.Claim) error {
	return nil
}
