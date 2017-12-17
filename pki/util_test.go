package pki

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"path"
	"testing"

	"github.com/moorara/go-box/util"
	"github.com/stretchr/testify/assert"
)

const (
	testKeyLen = 1024
)

func mockWorkspaceWithCA(t *testing.T) {
	err := NewWorkspace(NewState(), NewSpec())
	assert.NoError(t, err)

	// Mock root CA
	pub, priv, err := genKeyPair(testKeyLen)
	assert.NoError(t, err)
	rootCA := &x509.Certificate{
		SerialNumber: big.NewInt(10),
		Subject: pkix.Name{
			CommonName: "Root CA",
		},
	}
	pemData, err := x509.CreateCertificate(rand.Reader, rootCA, rootCA, pub, priv)
	assert.NoError(t, err)
	err = writePemFile(pemTypeCert, pemData, path.Join(DirRoot, "root"+extCACert))
	assert.NoError(t, err)

	// Mock an intermediate CA
	pub, priv, err = genKeyPair(testKeyLen)
	assert.NoError(t, err)
	opsCA := &x509.Certificate{
		SerialNumber: big.NewInt(100),
		Subject: pkix.Name{
			CommonName: "Ops CA",
		},
	}
	pemData, err = x509.CreateCertificate(rand.Reader, opsCA, rootCA, pub, priv)
	assert.NoError(t, err)
	err = writePemFile(pemTypeCert, pemData, path.Join(DirInterm, "ops"+extCACert))
	assert.NoError(t, err)

	// Mock first-level intermediate CA
	pub, priv, err = genKeyPair(testKeyLen)
	assert.NoError(t, err)
	sreCA := &x509.Certificate{
		SerialNumber: big.NewInt(200),
		Subject: pkix.Name{
			CommonName: "SRE CA",
		},
	}
	pemData, err = x509.CreateCertificate(rand.Reader, sreCA, rootCA, pub, priv)
	assert.NoError(t, err)
	err = writePemFile(pemTypeCert, pemData, path.Join(DirInterm, "sre"+extCACert))
	assert.NoError(t, err)

	// Mock second-level intermediate CA
	pub, priv, err = genKeyPair(testKeyLen)
	assert.NoError(t, err)
	rdCA := &x509.Certificate{
		SerialNumber: big.NewInt(300),
		Subject: pkix.Name{
			CommonName: "R&D CA",
		},
	}
	pemData, err = x509.CreateCertificate(rand.Reader, rdCA, sreCA, pub, priv)
	assert.NoError(t, err)
	err = writePemFile(pemTypeCert, pemData, path.Join(DirInterm, "rd"+extCACert))
	assert.NoError(t, err)
}

func TestGenKeyPair(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	tests := []struct {
		length      int
		expectError bool
	}{
		{0, true},
		{1024, false},
		{2048, false},
		{4096, false},
	}

	for _, test := range tests {
		pub, priv, err := genKeyPair(test.length)

		if test.expectError {
			assert.Error(t, err)
			assert.Nil(t, pub)
			assert.Nil(t, priv)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, pub)
			assert.NotNil(t, priv)
		}
	}
}

func TestComputeSubjectKeyID(t *testing.T) {
	pub, _, err := genKeyPair(testKeyLen)
	assert.NoError(t, err)

	tests := []struct {
		pubKey      interface{}
		expectError bool
	}{
		{nil, true},
		{pub, false},
	}

	for _, test := range tests {
		id, err := computeSubjectKeyID(test.pubKey)

		if test.expectError {
			assert.Error(t, err)
			assert.Nil(t, id)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, id)
		}
	}
}

func TestWriteReadPrivateKey(t *testing.T) {
	_, priv, err := genKeyPair(testKeyLen)
	assert.NoError(t, err)

	tests := []*struct {
		privKey    *rsa.PrivateKey
		writePW    string
		readPW     string
		setPath    bool
		path       string
		writeError bool
		readError  bool
	}{
		{
			priv,
			"", "",
			false, "",
			true, true,
		},
		{
			priv,
			"", "",
			true, "",
			false, false,
		},
		{
			priv,
			"secret", "secret",
			true, "",
			false, false,
		},
		{
			priv,
			"secret", "",
			true, "",
			false, true,
		},
		{
			priv,
			"secret", "different",
			true, "",
			false, true,
		},
	}

	// Prepare temporary files
	for _, test := range tests {
		if test.setPath {
			path, cleanup, err := util.CreateTempFile("")
			defer cleanup()
			assert.NoError(t, err)
			test.path = path
		}
	}

	t.Run("TestWritePrivateKey", func(t *testing.T) {
		for _, test := range tests {
			err := writePrivateKey(test.privKey, test.writePW, test.path)

			if test.writeError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		}
	})

	t.Run("TestReadPrivateKey", func(t *testing.T) {
		for _, test := range tests {
			privKey, err := readPrivateKey(test.readPW, test.path)

			if test.readError {
				assert.Error(t, err)
				assert.Nil(t, privKey)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.privKey, privKey)
			}
		}
	})
}

func TestWritePemFileReadCertificate(t *testing.T) {
	pub, priv, err := genKeyPair(testKeyLen)
	assert.NoError(t, err)

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(0),
		Subject: pkix.Name{
			CommonName: "Test Cert",
		},
	}
	certData, err := x509.CreateCertificate(rand.Reader, cert, cert, pub, priv)
	assert.NoError(t, err)

	tests := []*struct {
		pemType    string
		certData   []byte
		setPath    bool
		path       string
		writeError bool
		readError  bool
	}{
		{
			"", nil,
			false, "",
			true, true,
		},
		{
			pemTypeCert, certData,
			true, "",
			false, false,
		},
		{
			pemTypeCert, certData[1:],
			true, "",
			false, true,
		},
	}

	// Prepare temporary files
	for _, test := range tests {
		if test.setPath {
			path, cleanup, err := util.CreateTempFile("")
			defer cleanup()
			assert.NoError(t, err)
			test.path = path
		}
	}

	t.Run("TestWritePemFile", func(t *testing.T) {
		for _, test := range tests {
			err := writePemFile(test.pemType, test.certData, test.path)

			if test.writeError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		}
	})

	t.Run("TestReadCertificate", func(t *testing.T) {
		for _, test := range tests {
			cert, err := readCertificate(test.path)

			if test.readError {
				assert.Error(t, err)
				assert.Nil(t, cert)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cert)
			}
		}
	})
}

func TestWritePemFileReadCertificateRequest(t *testing.T) {
	_, priv, err := genKeyPair(testKeyLen)
	assert.NoError(t, err)

	csr := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: "Test CSR",
		},
	}
	csrData, err := x509.CreateCertificateRequest(rand.Reader, csr, priv)
	assert.NoError(t, err)

	tests := []*struct {
		pemType    string
		csrData    []byte
		setPath    bool
		path       string
		writeError bool
		readError  bool
	}{
		{
			"", nil,
			false, "",
			true, true,
		},
		{
			pemTypeCert, csrData,
			true, "",
			false, false,
		},
		{
			pemTypeCert, csrData[1:],
			true, "",
			false, true,
		},
	}

	// Prepare temporary files
	for _, test := range tests {
		if test.setPath {
			path, cleanup, err := util.CreateTempFile("")
			defer cleanup()
			assert.NoError(t, err)
			test.path = path
		}
	}

	t.Run("TestWritePemFile", func(t *testing.T) {
		for _, test := range tests {
			err := writePemFile(test.pemType, test.csrData, test.path)

			if test.writeError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		}
	})

	t.Run("TestReadCertificateRequest", func(t *testing.T) {
		for _, test := range tests {
			csr, err := readCertificateRequest(test.path)

			if test.readError {
				assert.Error(t, err)
				assert.Nil(t, csr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, csr)
			}
		}
	})
}

func TestWriteReadCertificateChain(t *testing.T) {
	mockWorkspaceWithCA(t)
	defer CleanupWorkspace()

	tests := []struct {
		title       string
		md          Metadata
		mdCA        Metadata
		expectError bool
	}{
		{
			"InvalidCertTypeRoot",
			Metadata{CertType: CertTypeRoot},
			Metadata{},
			true,
		},
		{
			"InvalidCertTypeServer",
			Metadata{CertType: CertTypeServer},
			Metadata{},
			true,
		},
		{
			"InvalidCertTypeClient",
			Metadata{CertType: CertTypeClient},
			Metadata{},
			true,
		},
		{
			"InvalidCATypeServer",
			Metadata{CertType: CertTypeInterm},
			Metadata{CertType: CertTypeServer},
			true,
		},
		{
			"InvalidCATypeClient",
			Metadata{CertType: CertTypeInterm},
			Metadata{CertType: CertTypeClient},
			true,
		},
		{
			"InvalidCertName",
			Metadata{CertType: CertTypeInterm},
			Metadata{CertType: CertTypeRoot},
			true,
		},
		{
			"CertNotExist",
			Metadata{Name: "interm", CertType: CertTypeInterm},
			Metadata{CertType: CertTypeRoot},
			true,
		},
		{
			"CANotExist",
			Metadata{Name: "ops", CertType: CertTypeInterm},
			Metadata{Name: "bad", CertType: CertTypeRoot},
			true,
		},
		{
			"RootInterm",
			Metadata{Name: "sre", CertType: CertTypeInterm},
			Metadata{Name: "root", CertType: CertTypeRoot},
			false,
		},
		{
			"IntermInterm",
			Metadata{Name: "rd", CertType: CertTypeInterm},
			Metadata{Name: "sre", CertType: CertTypeInterm},
			false,
		},
	}

	t.Run("TestWriteCertificateChain", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.title, func(t *testing.T) {
				err := writeCertificateChain(test.md, test.mdCA)

				if test.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("TestReadCertificateChain", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.title, func(t *testing.T) {
				certs, err := readCertificateChain(test.md.ChainPath())

				if test.expectError {
					assert.Error(t, err)
					assert.Nil(t, certs)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, certs)
				}
			})
		}
	})
}