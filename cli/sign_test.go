package cli

import (
	"errors"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/moorara/gocert/help"
	"github.com/moorara/gocert/pki"
	"github.com/stretchr/testify/assert"
)

func writeSignMocks(t *testing.T, mocks []pki.Cert) {
	for _, c := range mocks {
		err := ioutil.WriteFile(c.KeyPath(), []byte(""), 0644)
		assert.NoError(t, err)
		err = ioutil.WriteFile(c.CertPath(), []byte(""), 0644)
		assert.NoError(t, err)
		if c.CSRPath() != "" {
			err = ioutil.WriteFile(c.CSRPath(), []byte(""), 0644)
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
		mocks []pki.Cert
		args  []string
		input string
	}{
		{
			"RootSignsIntermediate",
			pki.NewState(),
			pki.NewSpec(),
			[]pki.Cert{
				pki.Cert{Name: rootName, Type: pki.CertTypeRoot},
				pki.Cert{Name: "ops", Type: pki.CertTypeInterm},
			},
			[]string{"-ca=root", "-name=ops"},
			`password
			password
			`,
		},
		{
			"IntermediateSignsServer",
			pki.NewState(),
			pki.NewSpec(),
			[]pki.Cert{
				pki.Cert{Name: "ops", Type: pki.CertTypeInterm},
				pki.Cert{Name: "server", Type: pki.CertTypeServer},
			},
			[]string{"-ca=ops"},
			`server
			password
			password
			`,
		},
		{
			"IntermediateSignsClient",
			pki.NewState(),
			pki.NewSpec(),
			[]pki.Cert{
				pki.Cert{Name: "ops", Type: pki.CertTypeInterm},
				pki.Cert{Name: "client", Type: pki.CertTypeClient},
			},
			[]string{},
			`ops
			client
			password
			password
			`,
		},
		{
			"IntermediateSignsServerClient",
			pki.NewState(),
			pki.NewSpec(),
			[]pki.Cert{
				pki.Cert{Name: "ops", Type: pki.CertTypeInterm},
				pki.Cert{Name: "server", Type: pki.CertTypeServer},
				pki.Cert{Name: "client", Type: pki.CertTypeClient},
			},
			[]string{"-ca=ops", "-name=server,client"},
			`password
			password
			`,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			err := pki.NewWorkspace(test.state, test.spec)
			assert.NoError(t, err)

			writeSignMocks(t, test.mocks)

			r := strings.NewReader(test.input)
			mockUI := help.NewMockUI(r)
			cmd := &SignCommand{
				ui:  mockUI,
				pki: &help.MockManager{},
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
		mocks        []pki.Cert
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
			interm
			`,
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
			"EnterNoPassword",
			pki.NewState(),
			pki.NewSpec(),
			[]pki.Cert{
				pki.Cert{Name: rootName, Type: pki.CertTypeRoot},
				pki.Cert{Name: "server", Type: pki.CertTypeServer},
			},
			[]string{"-ca=root", "-name=server"},
			``,
			nil,
			ErrorEnterConfig,
		},
		{
			"EnterInvalidPassword",
			pki.NewState(),
			pki.NewSpec(),
			[]pki.Cert{
				pki.Cert{Name: rootName, Type: pki.CertTypeRoot},
				pki.Cert{Name: "server", Type: pki.CertTypeServer},
			},
			[]string{"-ca=root", "-name=server"},
			`pass
			pass
			`,
			nil,
			ErrorEnterConfig,
		},
		{
			"EnterNotMatchingPassword",
			pki.NewState(),
			pki.NewSpec(),
			[]pki.Cert{
				pki.Cert{Name: rootName, Type: pki.CertTypeRoot},
				pki.Cert{Name: "server", Type: pki.CertTypeServer},
			},
			[]string{"-ca=root", "-name=server"},
			`secret
			password
			`,
			nil,
			ErrorEnterConfig,
		},
		{
			"CertNotExist",
			pki.NewState(),
			pki.NewSpec(),
			[]pki.Cert{
				pki.Cert{Name: rootName, Type: pki.CertTypeRoot},
			},
			[]string{"-ca=root", "-name=ops"},
			`password
			password
			`,
			nil,
			ErrorInvalidCSR,
		},
		{
			"CertNotExist",
			pki.NewState(),
			pki.NewSpec(),
			[]pki.Cert{
				pki.Cert{Name: "ops", Type: pki.CertTypeInterm},
			},
			[]string{"-ca=ops", "-name=server"},
			`password
			password
			`,
			nil,
			ErrorInvalidCSR,
		},
		{
			"RootCannotSignIntermediate",
			pki.NewState(),
			pki.NewSpec(),
			[]pki.Cert{
				pki.Cert{Name: rootName, Type: pki.CertTypeRoot},
				pki.Cert{Name: "server", Type: pki.CertTypeServer},
			},
			[]string{"-ca=root", "-name=server"},
			`password
			password
			`,
			nil,
			ErrorInvalidCSR,
		},
		{
			"IntermediateCannotSignRoot",
			pki.NewState(),
			pki.NewSpec(),
			[]pki.Cert{
				pki.Cert{Name: "interm", Type: pki.CertTypeInterm},
				pki.Cert{Name: rootName, Type: pki.CertTypeRoot},
			},
			[]string{"-ca=interm", "-name=root"},
			`password
			password
			`,
			nil,
			ErrorInvalidCSR,
		},
		{
			"SignCSRFails",
			pki.NewState(),
			pki.NewSpec(),
			[]pki.Cert{
				pki.Cert{Name: rootName, Type: pki.CertTypeRoot},
				pki.Cert{Name: "ops", Type: pki.CertTypeInterm},
			},
			[]string{"-ca=root", "-name=ops"},
			`password
			password
			`,
			errors.New("error"),
			ErrorSign,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			err := pki.NewWorkspace(test.state, test.spec)
			assert.NoError(t, err)

			writeSignMocks(t, test.mocks)

			r := strings.NewReader(test.input)
			mockUI := help.NewMockUI(r)
			cmd := &SignCommand{
				ui: mockUI,
				pki: &help.MockManager{
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
