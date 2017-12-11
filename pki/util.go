package pki

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/gob"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

const (
	pemKey  = "RSA PRIVATE KEY"
	pemCert = "CERTIFICATE"
	pemCSR  = "CERTIFICATE REQUEST"
)

func genKeyPair(length int) (*rsa.PublicKey, *rsa.PrivateKey, error) {
	private, err := rsa.GenerateKey(rand.Reader, length)
	if err != nil {
		return nil, nil, err
	}
	public := &private.PublicKey

	return public, private, nil
}

func computeSubjectKeyID(pubKey interface{}) ([]byte, error) {
	hash := sha1.New()
	pubKeyEncoder := gob.NewEncoder(hash)
	err := pubKeyEncoder.Encode(pubKey)
	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
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

func readPrivateKey(password, path string) (*rsa.PrivateKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	keyPem, _ := pem.Decode(data)
	if keyPem == nil {
		return nil, errors.New("decoding private key failed")
	}

	pemBytes := keyPem.Bytes

	// Decrypt private key if a password set
	if x509.IsEncryptedPEMBlock(keyPem) {
		if password != "" {
			pemBytes, err = x509.DecryptPEMBlock(keyPem, []byte(password))
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("password required but not set")
		}
	}

	private, err := x509.ParsePKCS1PrivateKey(pemBytes)
	if err != nil {
		return nil, err
	}

	return private, nil
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

func readCertificate(path string) (*x509.Certificate, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	certPem, _ := pem.Decode(data)
	if certPem == nil {
		return nil, errors.New("decoding certificate failed")
	}

	cert, err := x509.ParseCertificate(certPem.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func readCertificateRequest(path string) (*x509.CertificateRequest, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	csrPem, _ := pem.Decode(data)
	if csrPem == nil {
		return nil, errors.New("decoding certificate request failed")
	}

	csr, err := x509.ParseCertificateRequest(csrPem.Bytes)
	if err != nil {
		return nil, err
	}

	err = csr.CheckSignature()
	if err != nil {
		return nil, err
	}

	return csr, nil
}

func writeCertificateChain(md, mdCA Metadata) error {
	// Only an intermediate ca needs a certificate chain
	if md.CertType != CertTypeInterm {
		return errors.New("Only intermediate CAs have certificate chain")
	}

	// CA can only be root or another intermediate
	if mdCA.CertType != CertTypeRoot && mdCA.CertType != CertTypeInterm {
		return errors.New("CA can only be root ca or another intermediate CA")
	}

	chainBuf := new(bytes.Buffer)

	/* First, write certificate to chain */

	certFile, err := os.Open(md.CertPath())
	if err != nil {
		return err
	}

	_, err = io.Copy(chainBuf, certFile)
	if err != nil {
		return nil
	}

	err = certFile.Close()
	if err != nil {
		return nil
	}

	/* Next, wrtite the root ca certifiate or ca chain */

	caFile, err := os.Open(mdCA.ChainPath())
	if err != nil {
		return err
	}

	_, err = io.Copy(chainBuf, caFile)
	if err != nil {
		return err
	}

	err = caFile.Close()
	if err != nil {
		return err
	}

	/* Finally, write certificate chain file */

	chainFile, err := os.Create(md.ChainPath())
	if err != nil {
		return err
	}

	_, err = io.Copy(chainFile, chainBuf)
	if err != nil {
		return err
	}

	err = chainFile.Close()
	if err != nil {
		return err
	}

	return nil
}

func readCertificateChain(path string) ([]*x509.Certificate, error) {
	certs := make([]*x509.Certificate, 0)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	for len(data) > 0 {
		certPem, rest := pem.Decode(data)
		if certPem == nil {
			return nil, errors.New("decoding certificate failed")
		}

		cert, err := x509.ParseCertificate(certPem.Bytes)
		if err != nil {
			return nil, err
		}

		certs = append(certs, cert)
		data = rest
	}

	return certs, nil
}
