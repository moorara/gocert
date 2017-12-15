/*
 * https://tools.ietf.org/html/rfc5280
 */

package pki

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"log"
	"math/big"
	"path/filepath"
	"time"
)

type (
	// Manager provides methods for managing certificates
	Manager interface {
		GenCert(Config, Claim, Metadata) error
		GenCSR(Config, Claim, Metadata) error
		SignCSR(Config, Metadata, Config, Metadata, TrustFunc) error
		VerifyCert(Metadata, Metadata) error
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

// GenCert generates a new certificate
func (m *x509Manager) GenCert(config Config, claim Claim, md Metadata) error {
	if err := checkName(md.Name); err != nil {
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

	subjectKeyID, err := computeSubjectKeyID(publicKey)
	if err != nil {
		return err
	}

	// Declare certificate template
	cert := &x509.Certificate{
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
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageOCSPSigning},

		// Extensions:      []pkix.Extension{},
		// ExtraExtensions: []pkix.Extension{},

		// Authority Information Access
		// OCSPServer:            []string{},
		// IssuingCertificateURL: []string{},
	}

	// Create the certificate
	certData, err := x509.CreateCertificate(rand.Reader, cert, cert, publicKey, privateKey)
	if err != nil {
		return err
	}

	// Write certificate key file
	err = writePrivateKey(privateKey, config.Password, md.KeyPath())
	if err != nil {
		return err
	}

	// Write certificate file
	err = writePemFile(pemTypeCert, certData, md.CertPath())
	if err != nil {
		return err
	}

	return nil
}

// GenCSR generates a certificate signing request
func (m *x509Manager) GenCSR(config Config, claim Claim, md Metadata) error {
	if err := checkName(md.Name); err != nil {
		return err
	}

	length := config.Length

	// Generate a new public-private key pair
	_, privateKey, err := genKeyPair(length)
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
	err = writePrivateKey(privateKey, config.Password, md.KeyPath())
	if err != nil {
		return err
	}

	/* Write certificate request file */
	err = writePemFile(pemTypeCSR, csr, md.CSRPath())
	if err != nil {
		return err
	}

	return nil
}

// SignCSR signs a certificate signing request using a certificate authority
func (m *x509Manager) SignCSR(configCA Config, mdCA Metadata, configCSR Config, mdCSR Metadata, trust TrustFunc) error {
	keyCA, err := readPrivateKey(configCA.Password, mdCA.KeyPath())
	if err != nil {
		return err
	}

	certCA, err := readCertificate(mdCA.CertPath())
	if err != nil {
		return err
	}

	csr, err := readCertificateRequest(mdCSR.CSRPath())
	if err != nil {
		return err
	}

	// Check if the certificate authority can trust and sign the certificate request
	if !trust(certCA, csr) {
		return errors.New("CSR does not satisfy CA trust policy")
	}

	subjectKeyID, err := computeSubjectKeyID(csr.PublicKey)
	if err != nil {
		return err
	}

	configCSR.Serial++
	startTime := time.Now()
	endTime := startTime.AddDate(0, 0, configCSR.Days)

	// Declare certificate template
	cert := &x509.Certificate{
		Signature:          csr.Signature,
		SignatureAlgorithm: csr.SignatureAlgorithm,

		PublicKey:          csr.PublicKey,
		PublicKeyAlgorithm: csr.PublicKeyAlgorithm,

		SerialNumber: big.NewInt(configCSR.Serial),

		NotBefore: startTime,
		NotAfter:  endTime,

		Issuer:  certCA.Subject,
		Subject: csr.Subject,

		// DNSNames:    csr.DNSNames,
		// IPAddresses: csr.IPAddresses,
		EmailAddresses: csr.EmailAddresses,

		SubjectKeyId:   subjectKeyID,
		AuthorityKeyId: certCA.SubjectKeyId,

		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage: []x509.ExtKeyUsage{},

		// Extensions:      []pkix.Extension{},
		// ExtraExtensions: []pkix.Extension{},
	}

	switch mdCSR.CertType {
	case CertTypeInterm:
		cert.BasicConstraintsValid = true
		cert.IsCA = true
		cert.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign
		cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageOCSPSigning}
		// cert.OCSPServer = []string{}
		// cert.IssuingCertificateURL = []string{}
	case CertTypeServer:
		cert.BasicConstraintsValid = false
		cert.IsCA = false
		cert.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageContentCommitment | x509.KeyUsageKeyEncipherment | x509.KeyUsageDataEncipherment
		cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	case CertTypeClient:
		cert.BasicConstraintsValid = false
		cert.IsCA = false
		cert.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageContentCommitment | x509.KeyUsageKeyEncipherment | x509.KeyUsageDataEncipherment
		cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	}

	// Create the certificate
	certData, err := x509.CreateCertificate(rand.Reader, cert, certCA, csr.PublicKey, keyCA)
	if err != nil {
		return err
	}

	// Write certificate file
	err = writePemFile(pemTypeCert, certData, mdCSR.CertPath())
	if err != nil {
		return err
	}

	// Write certificate chain
	if mdCSR.CertType == CertTypeInterm {
		err = writeCertificateChain(mdCSR, mdCA)
		if err != nil {
			return err
		}
	}

	return nil
}

// VerifyCert
func (m *x509Manager) VerifyCert(md, mdCA Metadata) error {
	cert, err := readCertificate(md.CertPath())
	if err != nil {
		return err
	}

	chain, err := readCertificateChain(mdCA.ChainPath())
	if err != nil {
		return err
	}

	log.Printf("====> %+v %v \n", cert, chain)

	return nil
}
