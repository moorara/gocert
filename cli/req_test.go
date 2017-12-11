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
		state *pki.State
		spec  *pki.Spec
		md    pki.Metadata
		args  []string
		input string
	}{
		{
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeRoot},
			[]string{},
			`RootCA
			`,
		},
		{
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeInterm},
			[]string{"-name=ops"},
			`OpsCA
			`,
		},
		{
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeServer},
			[]string{"-name", "webapp"},
			`WebApp
			`,
		},
		{
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeClient},
			[]string{},
			`myservice
			MyService
			`,
		},
		{
			pki.NewState(),
			&pki.Spec{
				Interm: pki.Claim{
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
	}
}

func TestReqCommandError(t *testing.T) {
	tests := []struct {
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
			nil,
			nil,
			pki.Metadata{},
			[]string{},
			``,
			nil,
			nil,
			ErrorReadState,
		},
		{
			pki.NewState(),
			nil,
			pki.Metadata{},
			[]string{},
			``,
			nil,
			nil,
			ErrorReadSpec,
		},
		{
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{},
			[]string{},
			``,
			nil,
			nil,
			ErrorInvalidMetadata,
		},
		{
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
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeInterm},
			[]string{},
			``,
			nil,
			nil,
			ErrorInvalidName,
		},
		{
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeServer},
			[]string{},
			``,
			nil,
			nil,
			ErrorInvalidName,
		},
		{
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeClient},
			[]string{},
			``,
			nil,
			nil,
			ErrorInvalidName,
		},
		{
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeRoot},
			[]string{},
			``,
			errors.New("error"),
			errors.New("error"),
			ErrorCert,
		},
		{
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeInterm},
			[]string{"-name=ops"},
			``,
			errors.New("error"),
			errors.New("error"),
			ErrorCSR,
		},
		{
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeServer},
			[]string{"-name", "webapp"},
			``,
			errors.New("error"),
			errors.New("error"),
			ErrorCSR,
		},
		{
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeClient},
			[]string{"-name=ops"},
			`myservice
			`,
			errors.New("error"),
			errors.New("error"),
			ErrorCSR,
		},
	}

	for _, test := range tests {
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
	}
}
