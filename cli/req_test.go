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
		md pki.Metadata
	}{
		{
			pki.Metadata{CertType: pki.CertTypeRoot},
		},
		{
			pki.Metadata{CertType: pki.CertTypeInterm},
		},
		{
			pki.Metadata{CertType: pki.CertTypeServer},
		},
		{
			pki.Metadata{CertType: pki.CertTypeClient},
		},
	}

	for _, test := range tests {
		cmd := NewReqCommand(test.md)

		assert.Equal(t, newColoredUI(), cmd.ui)
		assert.Equal(t, pki.NewX509Manager(), cmd.pki)
		assert.Equal(t, test.md, cmd.md)
	}
}

func TestReqCommand(t *testing.T) {
	tests := []struct {
		state        *pki.State
		spec         *pki.Spec
		md           pki.Metadata
		args         []string
		input        string
		expectedExit int
	}{
		{
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeInterm},
			[]string{"-name=ops"},
			``,
			0,
		},
		{
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeServer},
			[]string{"-name=ops"},
			``,
			0,
		},
		{
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeClient},
			[]string{"-name=ops"},
			``,
			0,
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
			[]string{"-name", "ops"},
			`IntermediateCA
			Ottawa,Toronto
			R&D
			`,
			0,
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
			[]string{"-name", "ops"},
			`WebApp
			`,
			0,
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
			[]string{"-name", "ops"},
			`Service
			`,
			0,
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

		assert.Equal(t, reqSynopsis, cmd.Synopsis())
		assert.Equal(t, reqHelp, cmd.Help())
		assert.Equal(t, test.expectedExit, exit)

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
		GenCSRError  error
		expectedExit int
	}{
		{
			nil,
			nil,
			pki.Metadata{},
			[]string{"-invalid"},
			``,
			nil,
			ErrorInvalidFlag,
		},
		{
			nil,
			&pki.Spec{},
			pki.Metadata{},
			[]string{},
			``,
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
			ErrorReadSpec,
		},
		{
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{},
			[]string{},
			``,
			nil,
			ErrorMetadata,
		},
		{
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeRoot},
			[]string{},
			``,
			nil,
			ErrorMetadata,
		},
		{
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeInterm},
			[]string{},
			``,
			nil,
			ErrorNoName,
		},
		{
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeServer},
			[]string{},
			``,
			nil,
			ErrorNoName,
		},
		{
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeClient},
			[]string{},
			``,
			nil,
			ErrorNoName,
		},
		{
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeInterm},
			[]string{"-name=ops"},
			``,
			errors.New("error"),
			ErrorCSR,
		},
		{
			pki.NewState(),
			pki.NewSpec(),
			pki.Metadata{CertType: pki.CertTypeServer},
			[]string{""},
			`ops
      `,
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
				GenCSRError: test.GenCSRError,
			},
			md: test.md,
		}
		exit := cmd.Run(test.args)

		assert.Equal(t, reqSynopsis, cmd.Synopsis())
		assert.Equal(t, reqHelp, cmd.Help())
		assert.Equal(t, test.expectedExit, exit)

		err = pki.CleanupWorkspace()
		assert.NoError(t, err)
	}
}
