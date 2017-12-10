package pki

import (
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

func TestMetadataDir(t *testing.T) {
	tests := []struct {
		md          Metadata
		expectedDir string
	}{
		{
			Metadata{CertType: 0},
			"",
		},
		{
			Metadata{CertType: CertTypeRoot},
			DirRoot,
		},
		{
			Metadata{CertType: CertTypeInterm},
			DirInterm,
		},
		{
			Metadata{CertType: CertTypeServer},
			DirServer,
		},
		{
			Metadata{CertType: CertTypeClient},
			DirClient,
		},
	}

	for _, test := range tests {
		dir := test.md.Dir()

		assert.Equal(t, test.expectedDir, dir)
	}
}
