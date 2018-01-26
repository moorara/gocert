package pki

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"io/ioutil"
	"math/big"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/moorara/goto/util"
	"github.com/stretchr/testify/assert"
)

func parseKey(t *testing.T, password, path string) {
	key, err := readPrivateKey(password, path)
	assert.NoError(t, err)
	assert.NotNil(t, key)
}

func parseCSR(t *testing.T, path string) {
	csr, er := readCertificateRequest(path)
	assert.NoError(t, er)
	assert.NotNil(t, csr)
}

func parseCert(t *testing.T, path string) {
	cert, err := readCertificate(path)
	assert.NoError(t, err)
	assert.NotNil(t, cert)
}

func parseChain(t *testing.T, path string) {
	certs, err := readCertificateChain(path)
	assert.NoError(t, err)
	assert.NotNil(t, certs)
	assert.True(t, len(certs) >= 2)
}

func mockWorkspaceWithChains(t *testing.T) {
	err := NewWorkspace(NewState(), NewSpec())
	assert.NoError(t, err)

	// Mock root CA
	rootCertFile := path.Join(DirRoot, "root"+extCACert)
	pub, priv, err := genKeyPair(testKeyLen)
	assert.NoError(t, err)
	rootCA := &x509.Certificate{
		SerialNumber: big.NewInt(10),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(0, 0, 1),
		Subject: pkix.Name{
			CommonName: "Root CA",
		},
	}
	rootPem, err := x509.CreateCertificate(rand.Reader, rootCA, rootCA, pub, priv)
	assert.NoError(t, err)
	err = writePemFile(pemTypeCert, rootPem, rootCertFile)
	assert.NoError(t, err)

	// Mock first-level intermediate CA
	sreCertFile := path.Join(DirInterm, "sre"+extCACert)
	sreChainFile := path.Join(DirInterm, "sre"+extCAChain)
	pub, priv, err = genKeyPair(testKeyLen)
	assert.NoError(t, err)
	sreCA := &x509.Certificate{
		SerialNumber: big.NewInt(100),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(0, 0, 1),
		Subject: pkix.Name{
			CommonName: "SRE CA",
		},
	}
	srePem, err := x509.CreateCertificate(rand.Reader, sreCA, rootCA, pub, priv)
	assert.NoError(t, err)
	err = writePemFile(pemTypeCert, srePem, sreCertFile)
	assert.NoError(t, err)
	err = util.ConcatFiles(sreChainFile, false, sreCertFile, rootCertFile)
	assert.NoError(t, err)

	// Mock second-level intermediate CA
	rdCertFile := path.Join(DirInterm, "rd"+extCACert)
	rdChainFile := path.Join(DirInterm, "rd"+extCAChain)
	pub, priv, err = genKeyPair(testKeyLen)
	assert.NoError(t, err)
	rdCA := &x509.Certificate{
		SerialNumber: big.NewInt(200),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(0, 0, 1),
		Subject: pkix.Name{
			CommonName: "R&D CA",
		},
	}
	rdPem, err := x509.CreateCertificate(rand.Reader, rdCA, sreCA, pub, priv)
	assert.NoError(t, err)
	err = writePemFile(pemTypeCert, rdPem, rdCertFile)
	assert.NoError(t, err)
	err = util.ConcatFiles(rdChainFile, false, rdCertFile, sreCertFile, rootCertFile)
	assert.NoError(t, err)
}

func TestGenCertError(t *testing.T) {
	tests := []struct {
		title      string
		config     Config
		claim      Claim
		c          Cert
		writeFiles bool
	}{
		{
			"RootNoName",
			Config{},
			Claim{},
			Cert{
				Type: CertTypeRoot,
			},
			false,
		},
		{
			"RootExistingName",
			Config{},
			Claim{},
			Cert{
				Name: "root",
				Type: CertTypeRoot,
			},
			true,
		},
		{
			"RootNoConfig",
			Config{},
			Claim{},
			Cert{
				Name: "root",
				Type: CertTypeRoot,
			},
			false,
		},
	}

	err := NewWorkspace(NewState(), NewSpec())
	assert.NoError(t, err)
	defer CleanupWorkspace()

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			if test.writeFiles {
				err = ioutil.WriteFile(test.c.KeyPath(), nil, 0644)
				assert.NoError(t, err)
				err = ioutil.WriteFile(test.c.CertPath(), nil, 0644)
				assert.NoError(t, err)
			}

			manager := NewX509Manager()
			err = manager.GenCert(test.config, test.claim, test.c)
			assert.Error(t, err)

			err = util.DeleteAll("", test.c.KeyPath(), test.c.CertPath())
			assert.NoError(t, err)
		})
	}
}

func TestGenCSRError(t *testing.T) {
	tests := []struct {
		title      string
		config     Config
		claim      Claim
		c          Cert
		writeFiles bool
	}{
		{
			"IntermNoName",
			Config{},
			Claim{},
			Cert{
				Type: CertTypeInterm,
			},
			false,
		},
		{
			"ServerNoName",
			Config{},
			Claim{},
			Cert{
				Type: CertTypeServer,
			},
			false,
		},
		{
			"ClientNoName",
			Config{},
			Claim{},
			Cert{
				Type: CertTypeClient,
			},
			false,
		},
		{
			"IntermExistingName",
			Config{},
			Claim{},
			Cert{
				Name: "interm",
				Type: CertTypeInterm,
			},
			true,
		},
		{
			"ServerExistingName",
			Config{},
			Claim{},
			Cert{
				Name: "server",
				Type: CertTypeServer,
			},
			true,
		},
		{
			"ClientExistingName",
			Config{},
			Claim{},
			Cert{
				Name: "client",
				Type: CertTypeClient,
			},
			true,
		},
		{
			"IntermNoConfig",
			Config{},
			Claim{},
			Cert{
				Name: "interm",
				Type: CertTypeInterm,
			},
			false,
		},
		{
			"ServerNoConfig",
			Config{},
			Claim{},
			Cert{
				Name: "server",
				Type: CertTypeServer,
			},
			false,
		},
		{
			"ClientNoConfig",
			Config{},
			Claim{},
			Cert{
				Name: "client",
				Type: CertTypeClient,
			},
			false,
		},
	}

	err := NewWorkspace(NewState(), NewSpec())
	assert.NoError(t, err)
	defer CleanupWorkspace()

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			if test.writeFiles {
				err = ioutil.WriteFile(test.c.KeyPath(), nil, 0644)
				assert.NoError(t, err)
				err = ioutil.WriteFile(test.c.CSRPath(), nil, 0644)
				assert.NoError(t, err)
				err = ioutil.WriteFile(test.c.CertPath(), nil, 0644)
				assert.NoError(t, err)
			}

			manager := NewX509Manager()
			err = manager.GenCSR(test.config, test.claim, test.c)
			assert.Error(t, err)

			err = util.DeleteAll("", test.c.KeyPath(), test.c.CSRPath(), test.c.CertPath())
			assert.NoError(t, err)
		})
	}
}

func TestSignCSRError(t *testing.T) {
	tests := []struct {
		title     string
		configCA  Config
		cCA       Cert
		configCSR Config
		cCSR      Cert
		trust     TrustFunc
	}{
		{
			"NoCA",
			Config{},
			Cert{},
			Config{},
			Cert{},
			nil,
		},
		{
			"NoCSR",
			Config{
				Serial: 100,
				Length: 1024,
				Days:   3650,
			},
			Cert{
				Name: "root",
				Type: CertTypeRoot,
			},
			Config{},
			Cert{},
			nil,
		},
		{
			"CannotTrust",
			Config{
				Serial: 100,
				Length: 1024,
				Days:   3650,
			},
			Cert{
				Name: "root",
				Type: CertTypeRoot,
			},
			Config{
				Serial: 1000,
				Length: 1024,
				Days:   375,
			},
			Cert{
				Name: "interm",
				Type: CertTypeInterm,
			},
			func(*x509.Certificate, *x509.CertificateRequest) bool {
				return false
			},
		},
	}

	err := NewWorkspace(NewState(), NewSpec())
	assert.NoError(t, err)
	defer CleanupWorkspace()

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			files := make([]string, 0)
			manager := NewX509Manager()

			if !reflect.DeepEqual(test.configCA, Config{}) {
				manager.GenCert(test.configCA, Claim{}, test.cCA)
				assert.NoError(t, err)
				files = append(files, test.cCA.KeyPath(), test.cCA.CertPath())
			}

			if !reflect.DeepEqual(test.configCSR, Config{}) {
				manager.GenCSR(test.configCSR, Claim{}, test.cCSR)
				assert.NoError(t, err)
				files = append(files, test.cCSR.KeyPath(), test.cCSR.CSRPath(), test.cCSR.CertPath())
			}

			err := manager.SignCSR(test.configCA, test.cCA, test.configCSR, test.cCSR, test.trust)
			assert.Error(t, err)

			err = util.DeleteAll("", files...)
			assert.NoError(t, err)
		})
	}
}

func TestVerifyCertError(t *testing.T) {
	tests := []struct {
		title   string
		cCA     Cert
		c       Cert
		dnsName string
	}{
		{
			"InvalidCA",
			Cert{Type: CertTypeServer},
			Cert{},
			"",
		},
		{
			"CANotExist",
			Cert{Type: CertTypeInterm},
			Cert{Name: "sre", Type: CertTypeInterm},
			"",
		},
		{
			"CertNotExist",
			Cert{Name: "root", Type: CertTypeRoot},
			Cert{},
			"",
		},
		{
			"CannotVerify",
			Cert{Name: "root", Type: CertTypeRoot},
			Cert{Name: "rd", Type: CertTypeInterm},
			"",
		},
		{
			"CannotVerifyDNS",
			Cert{Name: "root", Type: CertTypeRoot},
			Cert{Name: "sre", Type: CertTypeInterm},
			"example.com",
		},
	}

	mockWorkspaceWithChains(t)
	defer CleanupWorkspace()

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			manager := NewX509Manager()

			err := manager.VerifyCert(test.cCA, test.c, test.dnsName)
			assert.Error(t, err)
		})
	}
}

func TestX509Manager(t *testing.T) {
	tests := []struct {
		title     string
		state     *State
		spec      *Spec
		cRoot     Cert
		cInterm   Cert
		cServer   Cert
		cClient   Cert
		dnsServer string
		dnsClient string
	}{
		{
			"RootIntermediate",
			&State{
				Root: Config{
					Serial: 10,
					Length: 1024,
					Days:   7300,
				},
				Interm: Config{
					Serial: 100,
					Length: 1024,
					Days:   3650,
				},
			},
			&Spec{
				Root: Claim{
					CommonName:   "Milad Root CA",
					Organization: []string{"Milad"},
				},
				Interm: Claim{
					CommonName:   "Milad SRE CA",
					Organization: []string{"Milad"},
				},
				RootPolicy: Policy{
					Match:    []string{"Organization"},
					Supplied: []string{"CommonName"},
				},
				IntermPolicy: Policy{
					Match:    []string{"Organization"},
					Supplied: []string{"CommonName"},
				},
			},
			Cert{
				Name: "root",
				Type: CertTypeRoot,
			},
			Cert{
				Name: "ops",
				Type: CertTypeInterm,
			},
			Cert{},
			Cert{},
			"",
			"",
		},
		{
			"RootIntermediateServer",
			&State{
				Root: Config{
					Serial:   10,
					Length:   1024,
					Days:     7300,
					Password: "rootSecret",
				},
				Interm: Config{
					Serial:   100,
					Length:   1024,
					Days:     3650,
					Password: "intermSecret",
				},
				Server: Config{
					Serial: 1000,
					Length: 1024,
					Days:   375,
				},
			},
			&Spec{
				Root: Claim{
					CommonName:   "Milad Root CA",
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Organization: []string{"Milad"},
				},
				Interm: Claim{
					CommonName:         "Milad SRE CA",
					Country:            []string{"CA"},
					Province:           []string{"Ontario"},
					Organization:       []string{"Milad"},
					OrganizationalUnit: []string{"SRE"},
					EmailAddress:       []string{"sre@example.com"},
				},
				Server: Claim{
					CommonName:         "milad.io",
					Country:            []string{"CA"},
					Organization:       []string{"Milad"},
					OrganizationalUnit: []string{"R&D"},
					DNSName:            []string{"milad.io"},
				},
				RootPolicy: Policy{
					Match:    []string{"Country", "Organization"},
					Supplied: []string{"CommonName", "EmailAddress"},
				},
				IntermPolicy: Policy{
					Match:    []string{"Organization"},
					Supplied: []string{"CommonName", "DNSName"},
				},
			},
			Cert{
				Name: "root",
				Type: CertTypeRoot,
			},
			Cert{
				Name: "ops",
				Type: CertTypeInterm,
			},
			Cert{
				Name: "milad.io",
				Type: CertTypeServer,
			},
			Cert{},
			"milad.io",
			"",
		},
		{
			"RootIntermediateServerClient",
			&State{
				Root: Config{
					Serial:   10,
					Length:   1024,
					Days:     7300,
					Password: "rootSecret",
				},
				Interm: Config{
					Serial:   100,
					Length:   1024,
					Days:     3650,
					Password: "intermSecret",
				},
				Server: Config{
					Serial: 1000,
					Length: 1024,
					Days:   375,
				},
				Client: Config{
					Serial: 10000,
					Length: 1024,
					Days:   40,
				},
			},
			&Spec{
				Root: Claim{
					CommonName:   "Milad Root CA",
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Organization: []string{"Milad"},
				},
				Interm: Claim{
					CommonName:         "Milad SRE CA",
					Country:            []string{"CA"},
					Province:           []string{"Ontario"},
					Organization:       []string{"Milad"},
					OrganizationalUnit: []string{"SRE"},
					EmailAddress:       []string{"sre@example.com"},
				},
				Server: Claim{
					CommonName:         "milad.io",
					Country:            []string{"CA"},
					Organization:       []string{"Milad"},
					OrganizationalUnit: []string{"R&D"},
					DNSName:            []string{"milad.io"},
				},
				Client: Claim{
					CommonName:         "auth.service",
					Country:            []string{"CA"},
					Organization:       []string{"Milad"},
					OrganizationalUnit: []string{"R&D"},
					DNSName:            []string{"auth.milad.io"},
					EmailAddress:       []string{"rd@example.com"},
				},
				RootPolicy: Policy{
					Match:    []string{"Country", "Organization"},
					Supplied: []string{"CommonName", "EmailAddress"},
				},
				IntermPolicy: Policy{
					Match:    []string{"Organization"},
					Supplied: []string{"CommonName", "DNSName"},
				},
			},
			Cert{
				Name: "root",
				Type: CertTypeRoot,
			},
			Cert{
				Name: "ops",
				Type: CertTypeInterm,
			},
			Cert{
				Name: "milad.io",
				Type: CertTypeServer,
			},
			Cert{
				Name: "auth.service",
				Type: CertTypeClient,
			},
			"milad.io",
			"auth.milad.io",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			err := NewWorkspace(test.state, test.spec)
			assert.NoError(t, err)

			manager := NewX509Manager()

			// Generate Root CA
			if !reflect.DeepEqual(test.state.Root, &Config{}) && !reflect.DeepEqual(test.spec.Root, &Claim{}) && !reflect.DeepEqual(test.cRoot, &Cert{}) {
				err = manager.GenCert(test.state.Root, test.spec.Root, test.cRoot)
				assert.NoError(t, err)

				parseKey(t, test.state.Root.Password, test.cRoot.KeyPath())
				parseCert(t, test.cRoot.CertPath())

				err = manager.VerifyCert(test.cRoot, test.cRoot, "")
				assert.NoError(t, err)
			}

			// Generate Intermediate CSR and sign it by Root CA
			if !reflect.DeepEqual(test.state.Interm, Config{}) && !reflect.DeepEqual(test.spec.Interm, Claim{}) && !reflect.DeepEqual(test.cInterm, Cert{}) {
				err = manager.GenCSR(test.state.Interm, test.spec.Interm, test.cInterm)
				assert.NoError(t, err)

				err = manager.SignCSR(test.state.Root, test.cRoot, test.state.Interm, test.cInterm, PolicyTrustFunc(test.spec.RootPolicy))
				assert.NoError(t, err)

				parseKey(t, test.state.Interm.Password, test.cInterm.KeyPath())
				parseCSR(t, test.cInterm.CSRPath())
				parseCert(t, test.cInterm.CertPath())
				parseChain(t, test.cInterm.ChainPath())

				err = manager.VerifyCert(test.cRoot, test.cInterm, "")
				assert.NoError(t, err)
			}

			// Generate Server CSR and it by Intermediate CA
			if !reflect.DeepEqual(test.state.Server, Config{}) && !reflect.DeepEqual(test.spec.Server, Claim{}) && !reflect.DeepEqual(test.cServer, Cert{}) {
				err = manager.GenCSR(test.state.Server, test.spec.Server, test.cServer)
				assert.NoError(t, err)

				err = manager.SignCSR(test.state.Interm, test.cInterm, test.state.Server, test.cServer, PolicyTrustFunc(test.spec.IntermPolicy))
				assert.NoError(t, err)

				parseKey(t, "", test.cServer.KeyPath())
				parseCSR(t, test.cServer.CSRPath())
				parseCert(t, test.cServer.CertPath())

				err = manager.VerifyCert(test.cInterm, test.cServer, test.dnsServer)
				assert.NoError(t, err)
			}

			// Generate Client CSR and sign it by Intermediate CA
			if !reflect.DeepEqual(test.state.Client, Config{}) && !reflect.DeepEqual(test.spec.Client, Claim{}) && !reflect.DeepEqual(test.cClient, Cert{}) {
				err = manager.GenCSR(test.state.Client, test.spec.Client, test.cClient)
				assert.NoError(t, err)

				err = manager.SignCSR(test.state.Interm, test.cInterm, test.state.Client, test.cClient, PolicyTrustFunc(test.spec.IntermPolicy))
				assert.NoError(t, err)

				parseKey(t, "", test.cClient.KeyPath())
				parseCSR(t, test.cClient.CSRPath())
				parseCert(t, test.cClient.CertPath())

				err = manager.VerifyCert(test.cInterm, test.cClient, test.dnsClient)
				assert.NoError(t, err)
			}

			err = CleanupWorkspace()
			assert.NoError(t, err)
		})
	}
}
