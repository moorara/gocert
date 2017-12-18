package cli

import (
	"errors"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
	"github.com/stretchr/testify/assert"
)

func writeSignMocks(t *testing.T, mocks []pki.Metadata) {
	for _, md := range mocks {
		err := ioutil.WriteFile(md.KeyPath(), []byte(""), 0644)
		assert.NoError(t, err)
		err = ioutil.WriteFile(md.CertPath(), []byte(""), 0644)
		assert.NoError(t, err)
		if md.CSRPath() != "" {
			err = ioutil.WriteFile(md.CSRPath(), []byte(""), 0644)
			assert.NoError(t, err)
		}
	}
}

func TestNewSignCommand(t *testing.T) {
	tests := []struct {
		expectedSynopsis string
	}{
		{
			"Signs a certificate signing request.",
		},
		{
			"Signs a certificate signing request.",
		},
	}

	for _, test := range tests {
		cmd := NewSignCommand()

		assert.Equal(t, newColoredUI(), cmd.ui)
		assert.Equal(t, pki.NewX509Manager(), cmd.pki)

		assert.Equal(t, test.expectedSynopsis, cmd.Synopsis())
		assert.NotEmpty(t, cmd.Help())
	}
}

func TestSignCommand(t *testing.T) {
	tests := []struct {
		title string
		state *pki.State
		spec  *pki.Spec
		mocks []pki.Metadata
		args  []string
		input string
	}{
		{
			"RootSignsIntermediate",
			pki.NewState(),
			pki.NewSpec(),
			[]pki.Metadata{
				pki.Metadata{Name: rootName, CertType: pki.CertTypeRoot},
				pki.Metadata{Name: "ops", CertType: pki.CertTypeInterm},
			},
			[]string{"-ca=root", "-name=ops"},
			``,
		},
		{
			"IntermediateSignsServer",
			pki.NewState(),
			pki.NewSpec(),
			[]pki.Metadata{
				pki.Metadata{Name: "ops", CertType: pki.CertTypeInterm},
				pki.Metadata{Name: "server", CertType: pki.CertTypeServer},
			},
			[]string{"-ca=ops"},
			`server
			`,
		},
		{
			"IntermediateSignsClient",
			pki.NewState(),
			pki.NewSpec(),
			[]pki.Metadata{
				pki.Metadata{Name: "ops", CertType: pki.CertTypeInterm},
				pki.Metadata{Name: "client", CertType: pki.CertTypeClient},
			},
			[]string{},
			`ops
			client
			`,
		},
		{
			"IntermediateSignsServerClient",
			pki.NewState(),
			pki.NewSpec(),
			[]pki.Metadata{
				pki.Metadata{Name: "ops", CertType: pki.CertTypeInterm},
				pki.Metadata{Name: "server", CertType: pki.CertTypeServer},
				pki.Metadata{Name: "client", CertType: pki.CertTypeClient},
			},
			[]string{"-ca=ops", "-name=server,client"},
			``,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			err := pki.NewWorkspace(test.state, test.spec)
			assert.NoError(t, err)

			writeSignMocks(t, test.mocks)

			mockUI := cli.NewMockUi()
			mockUI.InputReader = strings.NewReader(test.input)

			cmd := &SignCommand{
				ui:  mockUI,
				pki: &mockedManager{},
			}

			exit := cmd.Run(test.args)
			assert.Zero(t, exit)

			err = pki.CleanupWorkspace()
			assert.NoError(t, err)
		})
	}
}

func TestSignCommandError(t *testing.T) {
	tests := []struct {
		title        string
		state        *pki.State
		spec         *pki.Spec
		mocks        []pki.Metadata
		args         []string
		input        string
		SignCSRError error
		expectedExit int
	}{
		{
			"InvalidFlag",
			nil,
			nil,
			nil,
			[]string{"-invalid"},
			``,
			nil,
			ErrorInvalidFlag,
		},
		{
			"NoCAName",
			nil,
			nil,
			nil,
			[]string{},
			``,
			nil,
			ErrorInvalidCA,
		},
		{
			"NoCertName",
			nil,
			nil,
			nil,
			[]string{},
			`root
			`,
			nil,
			ErrorInvalidCSR,
		},
		{
			"SameCACertName",
			nil,
			nil,
			nil,
			[]string{"-name=interm"},
			`interm
			`,
			nil,
			ErrorInvalidName,
		},
		{
			"NoState",
			nil,
			nil,
			nil,
			[]string{},
			`root
			interm`,
			nil,
			ErrorReadState,
		},
		{
			"NoSpec",
			pki.NewState(),
			nil,
			nil,
			[]string{"-ca=root", "-name=interm"},
			``,
			nil,
			ErrorReadSpec,
		},
		{
			"CANotExist",
			pki.NewState(),
			pki.NewSpec(),
			nil,
			[]string{"-ca=root", "-name=interm"},
			``,
			nil,
			ErrorInvalidCA,
		},
		{
			"CANotExist",
			pki.NewState(),
			pki.NewSpec(),
			nil,
			[]string{"-ca=ops", "-name=server"},
			``,
			nil,
			ErrorInvalidCA,
		},
		{
			"CertNotExist",
			pki.NewState(),
			pki.NewSpec(),
			[]pki.Metadata{
				pki.Metadata{Name: rootName, CertType: pki.CertTypeRoot},
			},
			[]string{"-ca=root", "-name=ops"},
			``,
			nil,
			ErrorInvalidCSR,
		},
		{
			"CertNotExist",
			pki.NewState(),
			pki.NewSpec(),
			[]pki.Metadata{
				pki.Metadata{Name: "ops", CertType: pki.CertTypeInterm},
			},
			[]string{"-ca=ops", "-name=server"},
			``,
			nil,
			ErrorInvalidCSR,
		},
		{
			"RootCannotSignIntermediate",
			pki.NewState(),
			pki.NewSpec(),
			[]pki.Metadata{
				pki.Metadata{Name: rootName, CertType: pki.CertTypeRoot},
				pki.Metadata{Name: "server", CertType: pki.CertTypeServer},
			},
			[]string{"-ca=root", "-name=server"},
			``,
			nil,
			ErrorInvalidCSR,
		},
		{
			"IntermediateCannotSignRoot",
			pki.NewState(),
			pki.NewSpec(),
			[]pki.Metadata{
				pki.Metadata{Name: "interm", CertType: pki.CertTypeInterm},
				pki.Metadata{Name: rootName, CertType: pki.CertTypeRoot},
			},
			[]string{"-ca=interm", "-name=root"},
			``,
			nil,
			ErrorInvalidCSR,
		},
		{
			"SignCSRError",
			pki.NewState(),
			pki.NewSpec(),
			[]pki.Metadata{
				pki.Metadata{Name: rootName, CertType: pki.CertTypeRoot},
				pki.Metadata{Name: "ops", CertType: pki.CertTypeInterm},
			},
			[]string{"-ca=root", "-name=ops"},
			``,
			errors.New("error"),
			ErrorSign,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			err := pki.NewWorkspace(test.state, test.spec)
			assert.NoError(t, err)

			writeSignMocks(t, test.mocks)

			mockUI := cli.NewMockUi()
			mockUI.InputReader = strings.NewReader(test.input)

			cmd := &SignCommand{
				ui: mockUI,
				pki: &mockedManager{
					SignCSRError: test.SignCSRError,
				},
			}

			exit := cmd.Run(test.args)
			assert.Equal(t, test.expectedExit, exit)

			err = pki.CleanupWorkspace()
			assert.NoError(t, err)
		})
	}
}
