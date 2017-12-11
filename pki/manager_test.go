package pki

import (
	"crypto/x509"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/moorara/go-box/util"
	"github.com/stretchr/testify/assert"
)

type Meta struct {
	Root   Metadata
	Interm Metadata
	Server Metadata
	Client Metadata
}

func verifyKey(t *testing.T, password, path string) {
	key, err := readPrivateKey(password, path)
	assert.NoError(t, err)
	assert.NotNil(t, key)
}

func verifyCSR(t *testing.T, path string) {
	csr, er := readCertificateRequest(path)
	assert.NoError(t, er)
	assert.NotNil(t, csr)
}

func verifyCert(t *testing.T, path string) {
	cert, err := readCertificate(path)
	assert.NoError(t, err)
	assert.NotNil(t, cert)
}

func verifyChain(t *testing.T, path string) {
	certs, err := readCertificateChain(path)
	assert.NoError(t, err)
	assert.NotNil(t, certs)
	assert.True(t, len(certs) >= 2)
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
		title string
		md    Metadata
		mdCA  Metadata
	}{
		{
			"InvalidCertMetadata",
			Metadata{},
			Metadata{},
		},
	}

	err := NewWorkspace(NewState(), NewSpec())
	assert.NoError(t, err)
	defer CleanupWorkspace()

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			manager := NewX509Manager()

			err := manager.VerifyCert(test.md, test.mdCA)
			assert.Error(t, err)

			err = CleanupWorkspace()
			assert.NoError(t, err)
		})
	}
}

func TestX509Manager(t *testing.T) {
	tests := []struct {
		title string
		state *State
		spec  *Spec
		meta  *Meta
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
					CommonName:   "Milad Ops CA",
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
			&Meta{
				Root: Metadata{
					Name:     "root",
					CertType: CertTypeRoot,
				},
				Interm: Metadata{
					Name:     "ops",
					CertType: CertTypeInterm,
				},
			},
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
					CommonName:         "Milad Ops CA",
					Country:            []string{"CA"},
					Province:           []string{"Ontario"},
					Organization:       []string{"Milad"},
					OrganizationalUnit: []string{"Ops"},
					EmailAddress:       []string{"ops@example.com"},
				},
				Server: Claim{
					CommonName:         "milad.io",
					Country:            []string{"CA"},
					Organization:       []string{"Milad"},
					OrganizationalUnit: []string{"R&D"},
					EmailAddress:       []string{"rd@example.com"},
				},
				RootPolicy: Policy{
					Match:    []string{"Country", "Organization"},
					Supplied: []string{"CommonName", "EmailAddress"},
				},
				IntermPolicy: Policy{
					Match:    []string{"Organization"},
					Supplied: []string{"CommonName", "EmailAddress"},
				},
			},
			&Meta{
				Root: Metadata{
					Name:     "root",
					CertType: CertTypeRoot,
				},
				Interm: Metadata{
					Name:     "ops",
					CertType: CertTypeInterm,
				},
				Server: Metadata{
					Name:     "milad.io",
					CertType: CertTypeServer,
				},
			},
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
					CommonName:         "Milad Ops CA",
					Country:            []string{"CA"},
					Province:           []string{"Ontario"},
					Organization:       []string{"Milad"},
					OrganizationalUnit: []string{"Ops"},
					EmailAddress:       []string{"ops@example.com"},
				},
				Server: Claim{
					CommonName:         "milad.io",
					Country:            []string{"CA"},
					Organization:       []string{"Milad"},
					OrganizationalUnit: []string{"R&D"},
					EmailAddress:       []string{"rd@example.com"},
				},
				Client: Claim{
					CommonName:         "auth.service",
					Country:            []string{"CA"},
					Organization:       []string{"Milad"},
					OrganizationalUnit: []string{"R&D"},
					EmailAddress:       []string{"rd@example.com"},
				},
				RootPolicy: Policy{
					Match:    []string{"Country", "Organization"},
					Supplied: []string{"CommonName", "EmailAddress"},
				},
				IntermPolicy: Policy{
					Match:    []string{"Organization"},
					Supplied: []string{"CommonName", "EmailAddress"},
				},
			},
			&Meta{
				Root: Metadata{
					Name:     "root",
					CertType: CertTypeRoot,
				},
				Interm: Metadata{
					Name:     "ops",
					CertType: CertTypeInterm,
				},
				Server: Metadata{
					Name:     "milad.io",
					CertType: CertTypeServer,
				},
				Client: Metadata{
					Name:     "auth.service",
					CertType: CertTypeClient,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			err := NewWorkspace(test.state, test.spec)
			assert.NoError(t, err)

			manager := NewX509Manager()

			// Generate Root CA
			if !reflect.DeepEqual(test.state.Root, &Config{}) && !reflect.DeepEqual(test.spec.Root, &Claim{}) && !reflect.DeepEqual(test.meta.Root, &Metadata{}) {
				err = manager.GenCert(test.state.Root, test.spec.Root, test.meta.Root)
				assert.NoError(t, err)

				verifyKey(t, test.state.Root.Password, test.meta.Root.KeyPath())
				verifyCert(t, test.meta.Root.CertPath())
			}

			// Generate Intermediate CSR and sign it by Root CA
			if !reflect.DeepEqual(test.state.Interm, Config{}) && !reflect.DeepEqual(test.spec.Interm, Claim{}) && !reflect.DeepEqual(test.meta.Interm, Metadata{}) {
				err = manager.GenCSR(test.state.Interm, test.spec.Interm, test.meta.Interm)
				assert.NoError(t, err)

				err = manager.SignCSR(test.state.Root, test.meta.Root, test.state.Interm, test.meta.Interm, PolicyTrustFunc(test.spec.RootPolicy))
				assert.NoError(t, err)

				verifyKey(t, test.state.Interm.Password, test.meta.Interm.KeyPath())
				verifyCSR(t, test.meta.Interm.CSRPath())
				verifyCert(t, test.meta.Interm.CertPath())
				verifyChain(t, test.meta.Interm.ChainPath())
			}

			// Generate Server CSR and it by Intermediate CA
			if !reflect.DeepEqual(test.state.Server, Config{}) && !reflect.DeepEqual(test.spec.Server, Claim{}) && !reflect.DeepEqual(test.meta.Server, Metadata{}) {
				err = manager.GenCSR(test.state.Server, test.spec.Server, test.meta.Server)
				assert.NoError(t, err)

				err = manager.SignCSR(test.state.Interm, test.meta.Interm, test.state.Server, test.meta.Server, PolicyTrustFunc(test.spec.IntermPolicy))
				assert.NoError(t, err)

				verifyKey(t, "", test.meta.Server.KeyPath())
				verifyCSR(t, test.meta.Server.CSRPath())
				verifyCert(t, test.meta.Server.CertPath())
			}

			// Generate Client CSR and sign it by Intermediate CA
			if !reflect.DeepEqual(test.state.Client, Config{}) && !reflect.DeepEqual(test.spec.Client, Claim{}) && !reflect.DeepEqual(test.meta.Client, Metadata{}) {
				err = manager.GenCSR(test.state.Client, test.spec.Client, test.meta.Client)
				assert.NoError(t, err)

				err = manager.SignCSR(test.state.Interm, test.meta.Interm, test.state.Client, test.meta.Client, PolicyTrustFunc(test.spec.IntermPolicy))
				assert.NoError(t, err)

				verifyKey(t, "", test.meta.Client.KeyPath())
				verifyCSR(t, test.meta.Client.CSRPath())
				verifyCert(t, test.meta.Client.CertPath())
			}

			err = CleanupWorkspace()
			assert.NoError(t, err)
		})
	}
}
