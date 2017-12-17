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

	"github.com/moorara/go-box/util"
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
		md         Metadata
		writeFiles bool
	}{
		{
			"RootNoName",
			Config{},
			Claim{},
			Metadata{
				CertType: CertTypeRoot,
			},
			false,
		},
		{
			"RootExistingName",
			Config{},
			Claim{},
			Metadata{
				Name:     "root",
				CertType: CertTypeRoot,
			},
			true,
		},
		{
			"RootNoConfig",
			Config{},
			Claim{},
			Metadata{
				Name:     "root",
				CertType: CertTypeRoot,
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
				err = ioutil.WriteFile(test.md.KeyPath(), nil, 0644)
				assert.NoError(t, err)
				err = ioutil.WriteFile(test.md.CertPath(), nil, 0644)
				assert.NoError(t, err)
			}

			manager := NewX509Manager()
			err = manager.GenCert(test.config, test.claim, test.md)
			assert.Error(t, err)

			err = util.DeleteAll("", test.md.KeyPath(), test.md.CertPath())
			assert.NoError(t, err)
		})
	}
}

func TestGenCSRError(t *testing.T) {
	tests := []struct {
		title      string
		config     Config
		claim      Claim
		md         Metadata
		writeFiles bool
	}{
		{
			"IntermNoName",
			Config{},
			Claim{},
			Metadata{
				CertType: CertTypeInterm,
			},
			false,
		},
		{
			"ServerNoName",
			Config{},
			Claim{},
			Metadata{
				CertType: CertTypeServer,
			},
			false,
		},
		{
			"ClientNoName",
			Config{},
			Claim{},
			Metadata{
				CertType: CertTypeClient,
			},
			false,
		},
		{
			"IntermExistingName",
			Config{},
			Claim{},
			Metadata{
				Name:     "interm",
				CertType: CertTypeInterm,
			},
			true,
		},
		{
			"ServerExistingName",
			Config{},
			Claim{},
			Metadata{
				Name:     "server",
				CertType: CertTypeServer,
			},
			true,
		},
		{
			"ClientExistingName",
			Config{},
			Claim{},
			Metadata{
				Name:     "client",
				CertType: CertTypeClient,
			},
			true,
		},
		{
			"IntermNoConfig",
			Config{},
			Claim{},
			Metadata{
				Name:     "interm",
				CertType: CertTypeInterm,
			},
			false,
		},
		{
			"ServerNoConfig",
			Config{},
			Claim{},
			Metadata{
				Name:     "server",
				CertType: CertTypeServer,
			},
			false,
		},
		{
			"ClientNoConfig",
			Config{},
			Claim{},
			Metadata{
				Name:     "client",
				CertType: CertTypeClient,
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
				err = ioutil.WriteFile(test.md.KeyPath(), nil, 0644)
				assert.NoError(t, err)
				err = ioutil.WriteFile(test.md.CSRPath(), nil, 0644)
				assert.NoError(t, err)
				err = ioutil.WriteFile(test.md.CertPath(), nil, 0644)
				assert.NoError(t, err)
			}

			manager := NewX509Manager()
			err = manager.GenCSR(test.config, test.claim, test.md)
			assert.Error(t, err)

			err = util.DeleteAll("", test.md.KeyPath(), test.md.CSRPath(), test.md.CertPath())
			assert.NoError(t, err)
		})
	}
}

func TestSignCSRError(t *testing.T) {
	tests := []struct {
		title     string
		configCA  Config
		mdCA      Metadata
		configCSR Config
		mdCSR     Metadata
		trust     TrustFunc
	}{
		{
			"NoCA",
			Config{},
			Metadata{},
			Config{},
			Metadata{},
			nil,
		},
		{
			"NoCSR",
			Config{
				Serial: 100,
				Length: 1024,
				Days:   3650,
			},
			Metadata{
				Name:     "root",
				CertType: CertTypeRoot,
			},
			Config{},
			Metadata{},
			nil,
		},
		{
			"CannotTrust",
			Config{
				Serial: 100,
				Length: 1024,
				Days:   3650,
			},
			Metadata{
				Name:     "root",
				CertType: CertTypeRoot,
			},
			Config{
				Serial: 1000,
				Length: 1024,
				Days:   375,
			},
			Metadata{
				Name:     "interm",
				CertType: CertTypeInterm,
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
				manager.GenCert(test.configCA, Claim{}, test.mdCA)
				assert.NoError(t, err)
				files = append(files, test.mdCA.KeyPath(), test.mdCA.CertPath())
			}

			if !reflect.DeepEqual(test.configCSR, Config{}) {
				manager.GenCSR(test.configCSR, Claim{}, test.mdCSR)
				assert.NoError(t, err)
				files = append(files, test.mdCSR.KeyPath(), test.mdCSR.CSRPath(), test.mdCSR.CertPath())
			}

			err := manager.SignCSR(test.configCA, test.mdCA, test.configCSR, test.mdCSR, test.trust)
			assert.Error(t, err)

			err = util.DeleteAll("", files...)
			assert.NoError(t, err)
		})
	}
}

func TestVerifyCertError(t *testing.T) {
	tests := []struct {
		title   string
		mdCA    Metadata
		md      Metadata
		dnsName string
	}{
		{
			"InvalidCA",
			Metadata{CertType: CertTypeServer},
			Metadata{},
			"",
		},
		{
			"CANotExist",
			Metadata{CertType: CertTypeInterm},
			Metadata{Name: "sre", CertType: CertTypeInterm},
			"",
		},
		{
			"CertNotExist",
			Metadata{Name: "root", CertType: CertTypeRoot},
			Metadata{},
			"",
		},
		{
			"CannotVerify",
			Metadata{Name: "root", CertType: CertTypeRoot},
			Metadata{Name: "rd", CertType: CertTypeInterm},
			"",
		},
		{
			"CannotVerifyDNS",
			Metadata{Name: "root", CertType: CertTypeRoot},
			Metadata{Name: "sre", CertType: CertTypeInterm},
			"example.com",
		},
	}

	mockWorkspaceWithChains(t)
	defer CleanupWorkspace()

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			manager := NewX509Manager()

			err := manager.VerifyCert(test.mdCA, test.md, test.dnsName)
			assert.Error(t, err)
		})
	}
}

func TestX509Manager(t *testing.T) {
	tests := []struct {
		title     string
		state     *State
		spec      *Spec
		mdRoot    Metadata
		mdInterm  Metadata
		mdServer  Metadata
		mdClient  Metadata
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
			Metadata{
				Name:     "root",
				CertType: CertTypeRoot,
			},
			Metadata{
				Name:     "ops",
				CertType: CertTypeInterm,
			},
			Metadata{},
			Metadata{},
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
			Metadata{
				Name:     "root",
				CertType: CertTypeRoot,
			},
			Metadata{
				Name:     "ops",
				CertType: CertTypeInterm,
			},
			Metadata{
				Name:     "milad.io",
				CertType: CertTypeServer,
			},
			Metadata{},
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
			Metadata{
				Name:     "root",
				CertType: CertTypeRoot,
			},
			Metadata{
				Name:     "ops",
				CertType: CertTypeInterm,
			},
			Metadata{
				Name:     "milad.io",
				CertType: CertTypeServer,
			},
			Metadata{
				Name:     "auth.service",
				CertType: CertTypeClient,
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
			if !reflect.DeepEqual(test.state.Root, &Config{}) && !reflect.DeepEqual(test.spec.Root, &Claim{}) && !reflect.DeepEqual(test.mdRoot, &Metadata{}) {
				err = manager.GenCert(test.state.Root, test.spec.Root, test.mdRoot)
				assert.NoError(t, err)

				parseKey(t, test.state.Root.Password, test.mdRoot.KeyPath())
				parseCert(t, test.mdRoot.CertPath())

				err = manager.VerifyCert(test.mdRoot, test.mdRoot, "")
				assert.NoError(t, err)
			}

			// Generate Intermediate CSR and sign it by Root CA
			if !reflect.DeepEqual(test.state.Interm, Config{}) && !reflect.DeepEqual(test.spec.Interm, Claim{}) && !reflect.DeepEqual(test.mdInterm, Metadata{}) {
				err = manager.GenCSR(test.state.Interm, test.spec.Interm, test.mdInterm)
				assert.NoError(t, err)

				err = manager.SignCSR(test.state.Root, test.mdRoot, test.state.Interm, test.mdInterm, PolicyTrustFunc(test.spec.RootPolicy))
				assert.NoError(t, err)

				parseKey(t, test.state.Interm.Password, test.mdInterm.KeyPath())
				parseCSR(t, test.mdInterm.CSRPath())
				parseCert(t, test.mdInterm.CertPath())
				parseChain(t, test.mdInterm.ChainPath())

				err = manager.VerifyCert(test.mdRoot, test.mdInterm, "")
				assert.NoError(t, err)
			}

			// Generate Server CSR and it by Intermediate CA
			if !reflect.DeepEqual(test.state.Server, Config{}) && !reflect.DeepEqual(test.spec.Server, Claim{}) && !reflect.DeepEqual(test.mdServer, Metadata{}) {
				err = manager.GenCSR(test.state.Server, test.spec.Server, test.mdServer)
				assert.NoError(t, err)

				err = manager.SignCSR(test.state.Interm, test.mdInterm, test.state.Server, test.mdServer, PolicyTrustFunc(test.spec.IntermPolicy))
				assert.NoError(t, err)

				parseKey(t, "", test.mdServer.KeyPath())
				parseCSR(t, test.mdServer.CSRPath())
				parseCert(t, test.mdServer.CertPath())

				err = manager.VerifyCert(test.mdInterm, test.mdServer, test.dnsServer)
				assert.NoError(t, err)
			}

			// Generate Client CSR and sign it by Intermediate CA
			if !reflect.DeepEqual(test.state.Client, Config{}) && !reflect.DeepEqual(test.spec.Client, Claim{}) && !reflect.DeepEqual(test.mdClient, Metadata{}) {
				err = manager.GenCSR(test.state.Client, test.spec.Client, test.mdClient)
				assert.NoError(t, err)

				err = manager.SignCSR(test.state.Interm, test.mdInterm, test.state.Client, test.mdClient, PolicyTrustFunc(test.spec.IntermPolicy))
				assert.NoError(t, err)

				parseKey(t, "", test.mdClient.KeyPath())
				parseCSR(t, test.mdClient.CSRPath())
				parseCert(t, test.mdClient.CertPath())

				err = manager.VerifyCert(test.mdInterm, test.mdClient, test.dnsClient)
				assert.NoError(t, err)
			}

			err = CleanupWorkspace()
			assert.NoError(t, err)
		})
	}
}
