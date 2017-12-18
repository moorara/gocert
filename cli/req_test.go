package cli

import (
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
	"github.com/stretchr/testify/assert"
)

func TestNewReqCommand(t *testing.T) {
	tests := []struct {
		md               pki.Metadata
		expectedSynopsis string
	}{
		{
			pki.Metadata{CertType: pki.CertTypeRoot},
			"Creates a new root certificate authority.",
		},
		{
			pki.Metadata{CertType: pki.CertTypeInterm},
			"Creates a new certificate signing request.",
		},
		{
			pki.Metadata{CertType: pki.CertTypeServer},
			"Creates a new certificate signing request.",
		},
		{
			pki.Metadata{CertType: pki.CertTypeClient},
			"Creates a new certificate signing request.",
		},
	}

	for _, test := range tests {
		cmd := NewReqCommand(test.md)

		assert.Equal(t, newColoredUI(), cmd.ui)
		assert.Equal(t, pki.NewX509Manager(), cmd.pki)
		assert.Equal(t, test.md, cmd.md)

		assert.Equal(t, test.expectedSynopsis, cmd.Synopsis())
		assert.NotEmpty(t, cmd.Help())
	}
}

func TestReqCommand(t *testing.T) {
	tests := []struct {
		title string
		state *pki.State
		spec  *pki.Spec
		md    pki.Metadata
		args  []string
		input string
	}{
		{
			"GenSimpleRootCA",
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeRoot},
			[]string{},
			`RootCA
			`,
		},
		{
			"GenSimpleIntermediateCA",
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeInterm},
			[]string{"-name=ops"},
			`OpsCA
			`,
		},
		{
			"GenSimpleServerCert",
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeServer},
			[]string{"-name", "webapp"},
			`WebApp
			`,
		},
		{
			"GenSimpleClientCert",
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeClient},
			[]string{},
			`myservice
			MyService
			`,
		},
		{
			"GenRootCA",
			pki.NewState(),
			&pki.Spec{
				Root: pki.Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Organization: []string{"Milad"},
				},
			},
			pki.Metadata{CertType: pki.CertTypeRoot},
			[]string{},
			`RootCA
			`,
		},
		{
			"GenIntermediateCA",
			pki.NewState(),
			&pki.Spec{
				Interm: pki.Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Organization: []string{"Milad"},
				},
			},
			pki.Metadata{CertType: pki.CertTypeInterm},
			[]string{"-name=ops"},
			`IntermediateCA
			Ottawa
			R&D
			`,
		},
		{
			"GenServerCert",
			pki.NewState(),
			&pki.Spec{
				Server: pki.Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Locality:     []string{"Toronto"},
					Organization: []string{"Milad"},
				},
			},
			pki.Metadata{CertType: pki.CertTypeServer},
			[]string{"-name", "webapp"},
			`WebApp
			`,
		},
		{
			"GenClientCert",
			pki.NewState(),
			&pki.Spec{
				Client: pki.Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Locality:     []string{"Montreal"},
					Organization: []string{"Milad"},
				},
			},
			pki.Metadata{CertType: pki.CertTypeClient},
			[]string{},
			`myservice
			MyService
			`,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			err := pki.NewWorkspace(test.state, test.spec)
			assert.NoError(t, err)

			mockUI := cli.NewMockUi()
			mockUI.InputReader = strings.NewReader(test.input)

			cmd := &ReqCommand{
				ui:  mockUI,
				pki: &mockedManager{},
				md:  test.md,
			}

			exit := cmd.Run(test.args)
			assert.Zero(t, exit)

			err = pki.CleanupWorkspace()
			assert.NoError(t, err)
		})
	}
}

func TestReqCommandError(t *testing.T) {
	tests := []struct {
		title        string
		state        *pki.State
		spec         *pki.Spec
		md           pki.Metadata
		args         []string
		input        string
		GenCertError error
		GenCSRError  error
		expectedExit int
	}{
		{
			"InvalidFlag",
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeRoot},
			[]string{"-invalid"},
			``,
			nil,
			nil,
			ErrorInvalidFlag,
		},
		{
			"NoName",
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{},
			[]string{},
			``,
			nil,
			nil,
			ErrorInvalidName,
		},
		{
			"NoState",
			nil,
			nil,
			pki.Metadata{},
			[]string{},
			`sre
			`,
			nil,
			nil,
			ErrorReadState,
		},
		{
			"NoSpec",
			pki.NewState(),
			nil,
			pki.Metadata{},
			[]string{"-name=sre"},
			``,
			nil,
			nil,
			ErrorReadSpec,
		},
		{
			"NoMetadata",
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{},
			[]string{"-name", "sre"},
			``,
			nil,
			nil,
			ErrorInvalidMetadata,
		},
		{
			"GenCertError",
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeRoot},
			[]string{},
			``,
			errors.New("error"),
			nil,
			ErrorCert,
		},
		{
			"GenCSRError",
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeInterm},
			[]string{"-name=ops"},
			``,
			nil,
			errors.New("error"),
			ErrorCSR,
		},
		{
			"GenCSRError",
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeServer},
			[]string{"-name", "webapp"},
			``,
			nil,
			errors.New("error"),
			ErrorCSR,
		},
		{
			"GenCSRError",
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeClient},
			[]string{"-name=ops"},
			`myservice
			`,
			nil,
			errors.New("error"),
			ErrorCSR,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			err := pki.NewWorkspace(test.state, test.spec)
			assert.NoError(t, err)

			mockUI := cli.NewMockUi()
			mockUI.InputReader = strings.NewReader(test.input)

			cmd := &ReqCommand{
				ui: mockUI,
				pki: &mockedManager{
					GenCertError: test.GenCertError,
					GenCSRError:  test.GenCSRError,
				},
				md: test.md,
			}

			exit := cmd.Run(test.args)
			assert.Equal(t, test.expectedExit, exit)

			err = pki.CleanupWorkspace()
			assert.NoError(t, err)
		})
	}
}
