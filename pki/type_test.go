package pki

import (
	"net"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewState(t *testing.T) {
	state := NewState()

	assert.Equal(t, defaultRootCASerial, state.Root.Serial)
	assert.Equal(t, defaultRootCALength, state.Root.Length)
	assert.Equal(t, defaultRootCADays, state.Root.Days)

	assert.Equal(t, defaultIntermCASerial, state.Interm.Serial)
	assert.Equal(t, defaultIntermCALength, state.Interm.Length)
	assert.Equal(t, defaultIntermCADays, state.Interm.Days)

	assert.Equal(t, defaultServerCertSerial, state.Server.Serial)
	assert.Equal(t, defaultServerCertLength, state.Server.Length)
	assert.Equal(t, defaultServerCertDays, state.Server.Days)

	assert.Equal(t, defaultClientCertSerial, state.Client.Serial)
	assert.Equal(t, defaultClientCertLength, state.Client.Length)
	assert.Equal(t, defaultClientCertDays, state.Client.Days)
}

func TestNewSpec(t *testing.T) {
	spec := NewSpec()

	expectedClaim := Claim{}
	expectedRootPolicy := Policy{
		Match:    strings.Split(defaultRootPolicyMatch, ","),
		Supplied: strings.Split(defaultRootPolicySupplied, ","),
	}
	expectedIntermPolicy := Policy{
		Match:    strings.Split(defaultIntermPolicyMatch, ","),
		Supplied: strings.Split(defaultIntermPolicySupplied, ","),
	}

	assert.Equal(t, expectedClaim, spec.Root)
	assert.Equal(t, expectedClaim, spec.Interm)
	assert.Equal(t, expectedClaim, spec.Server)
	assert.Equal(t, expectedClaim, spec.Client)
	assert.Equal(t, expectedRootPolicy, spec.RootPolicy)
	assert.Equal(t, expectedIntermPolicy, spec.IntermPolicy)
}

func TestState(t *testing.T) {
	tests := []struct {
		state            State
		certType         int
		expectedConfig   Config
		expectedConfigOK bool
	}{
		{
			*NewState(),
			-1,
			Config{},
			false,
		},
		{
			*NewState(),
			CertTypeRoot,
			Config{
				Serial: defaultRootCASerial,
				Length: defaultRootCALength,
				Days:   defaultRootCADays,
			},
			true,
		},
		{
			*NewState(),
			CertTypeInterm,
			Config{
				Serial: defaultIntermCASerial,
				Length: defaultIntermCALength,
				Days:   defaultIntermCADays,
			},
			true,
		},
		{
			*NewState(),
			CertTypeServer,
			Config{
				Serial: defaultServerCertSerial,
				Length: defaultServerCertLength,
				Days:   defaultServerCertDays,
			},
			true,
		},
		{
			*NewState(),
			CertTypeClient,
			Config{
				Serial: defaultClientCertSerial,
				Length: defaultClientCertLength,
				Days:   defaultClientCertDays,
			},
			true,
		},
	}

	for _, test := range tests {
		config, ok := test.state.ConfigFor(test.certType)

		assert.Equal(t, test.expectedConfigOK, ok)
		assert.Equal(t, test.expectedConfig, config)
	}
}

func TestSpec(t *testing.T) {
	spec := Spec{
		Root: Claim{
			Country:      []string{"CA"},
			Province:     []string{"Ontario"},
			Locality:     []string{"Ottawa"},
			Organization: []string{"Milad"},
		},
		Interm: Claim{
			Country:            []string{"CA"},
			Province:           []string{"Ontario"},
			Locality:           []string{"Ottawa"},
			Organization:       []string{"Milad"},
			OrganizationalUnit: []string{"SRE"},
		},
		Server: Claim{
			Country:            []string{"CA"},
			Province:           []string{"Ontario"},
			Locality:           []string{"Ottawa"},
			Organization:       []string{"Milad"},
			OrganizationalUnit: []string{"R&D"},
		},
		Client: Claim{
			Country:            []string{"CA"},
			Province:           []string{"Ontario"},
			Locality:           []string{"Ottawa"},
			Organization:       []string{"Milad"},
			OrganizationalUnit: []string{"QE"},
		},
		RootPolicy: Policy{
			Match:    []string{"Organization"},
			Supplied: []string{"CommonName", "DNSName"},
		},
		IntermPolicy: Policy{
			Match:    []string{"Country", "Organization"},
			Supplied: []string{"CommonName", "DNSName", "EmailAddress"},
		},
	}

	tests := []struct {
		spec             Spec
		certType         int
		expectedClaim    Claim
		expectedClaimOK  bool
		expectedPolicy   Policy
		expectedPolicyOK bool
	}{
		{
			spec,
			-1,
			Claim{},
			false,
			Policy{},
			false,
		},
		{
			spec,
			CertTypeRoot,
			Claim{
				Country:      []string{"CA"},
				Province:     []string{"Ontario"},
				Locality:     []string{"Ottawa"},
				Organization: []string{"Milad"},
			},
			true,
			Policy{
				Match:    []string{"Organization"},
				Supplied: []string{"CommonName", "DNSName"},
			},
			true,
		},
		{
			spec,
			CertTypeInterm,
			Claim{
				Country:            []string{"CA"},
				Province:           []string{"Ontario"},
				Locality:           []string{"Ottawa"},
				Organization:       []string{"Milad"},
				OrganizationalUnit: []string{"SRE"},
			},
			true,
			Policy{
				Match:    []string{"Country", "Organization"},
				Supplied: []string{"CommonName", "DNSName", "EmailAddress"},
			},
			true,
		},
		{
			spec,
			CertTypeServer,
			Claim{
				Country:            []string{"CA"},
				Province:           []string{"Ontario"},
				Locality:           []string{"Ottawa"},
				Organization:       []string{"Milad"},
				OrganizationalUnit: []string{"R&D"},
			},
			true,
			Policy{},
			false,
		},
		{
			spec,
			CertTypeClient,
			Claim{
				Country:            []string{"CA"},
				Province:           []string{"Ontario"},
				Locality:           []string{"Ottawa"},
				Organization:       []string{"Milad"},
				OrganizationalUnit: []string{"QE"},
			},
			true,
			Policy{},
			false,
		},
	}

	for _, test := range tests {
		claim, ok := test.spec.ClaimFor(test.certType)
		assert.Equal(t, test.expectedClaimOK, ok)
		assert.Equal(t, test.expectedClaim, claim)

		policy, ok := test.spec.PolicyFor(test.certType)
		assert.Equal(t, test.expectedPolicyOK, ok)
		assert.Equal(t, test.expectedPolicy, policy)
	}
}

func TestClaim(t *testing.T) {
	tests := []struct {
		claim Claim
	}{
		{
			Claim{},
		},
		{
			Claim{
				CommonName: "Test",
			},
		},
		{
			Claim{
				CommonName:         "Test",
				Country:            []string{"CA"},
				Province:           []string{"Ontario"},
				Locality:           []string{"Ottawa"},
				Organization:       []string{"Milad"},
				OrganizationalUnit: []string{"R&D"},
			},
		},
		{
			Claim{
				CommonName:         "Test",
				Country:            []string{"CA"},
				Province:           []string{"Ontario"},
				Locality:           []string{"Ottawa"},
				Organization:       []string{"Milad"},
				OrganizationalUnit: []string{"R&D"},
				DNSName:            []string{"example.com"},
				IPAddress:          []net.IP{net.ParseIP("127.0.0.1")},
				EmailAddress:       []string{"milad@example.com"},
				StreetAddress:      []string{"Island Park Dr"},
				PostalCode:         []string{"K1Z"},
			},
		},
	}

	for _, test := range tests {
		claim := test.claim.Clone()

		assert.Equal(t, test.claim, claim)
	}
}

func TestCert(t *testing.T) {
	tests := []struct {
		c                 Cert
		expectedTitle     string
		expectedCertPath  string
		expectedKeyPath   string
		expectedCSRPath   string
		expectedChainPath string
	}{
		{
			Cert{},
			"",
			"",
			"",
			"",
			"",
		},
		{
			Cert{Name: "root"},
			"",
			"",
			"",
			"",
			"",
		},
		{
			Cert{
				Name: "root",
				Type: CertTypeRoot,
			},
			titleRoot,
			path.Join(DirRoot, "root"+extCACert),
			path.Join(DirRoot, "root"+extCAKey),
			"",
			path.Join(DirRoot, "root"+extCACert),
		},
		{
			Cert{
				Name: "ops",
				Type: CertTypeInterm,
			},
			titleInterm,
			path.Join(DirInterm, "ops"+extCACert),
			path.Join(DirInterm, "ops"+extCAKey),
			path.Join(DirCSR, "ops"+extCACSR),
			path.Join(DirInterm, "ops"+extCAChain),
		},
		{
			Cert{
				Name: "webapp",
				Type: CertTypeServer,
			},
			titleServer,
			path.Join(DirServer, "webapp"+extCert),
			path.Join(DirServer, "webapp"+extKey),
			path.Join(DirCSR, "webapp"+extCSR),
			"",
		},
		{
			Cert{
				Name: "service",
				Type: CertTypeClient,
			},
			titleClient,
			path.Join(DirClient, "service"+extCert),
			path.Join(DirClient, "service"+extKey),
			path.Join(DirCSR, "service"+extCSR),
			"",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedTitle, test.c.Title())
		assert.Equal(t, test.expectedCertPath, test.c.CertPath())
		assert.Equal(t, test.expectedKeyPath, test.c.KeyPath())
		assert.Equal(t, test.expectedCSRPath, test.c.CSRPath())
		assert.Equal(t, test.expectedChainPath, test.c.ChainPath())
	}
}
