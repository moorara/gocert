package cli

import (
	"errors"
	"strings"
	"testing"

	"github.com/moorara/gocert/help"
	"github.com/moorara/gocert/pki"
	"github.com/stretchr/testify/assert"
)

func TestNewReqCommand(t *testing.T) {
	tests := []struct {
		c                pki.Cert
		expectedSynopsis string
	}{
		{
			pki.Cert{Type: pki.CertTypeRoot},
			"Creates a new root certificate authority.",
		},
		{
			pki.Cert{Type: pki.CertTypeInterm},
			"Creates a new certificate signing request.",
		},
		{
			pki.Cert{Type: pki.CertTypeServer},
			"Creates a new certificate signing request.",
		},
		{
			pki.Cert{Type: pki.CertTypeClient},
			"Creates a new certificate signing request.",
		},
	}

	for _, test := range tests {
		cmd := NewReqCommand(test.c)

		assert.Equal(t, newColoredUI(), cmd.ui)
		assert.Equal(t, pki.NewX509Manager(), cmd.pki)
		assert.Equal(t, test.c, cmd.c)

		assert.Equal(t, test.expectedSynopsis, cmd.Synopsis())
		assert.NotEmpty(t, cmd.Help())
	}
}

func TestReqCommand(t *testing.T) {
	tests := []struct {
		title string
		state *pki.State
		spec  *pki.Spec
		c     pki.Cert
		args  []string
		input string
	}{
		{
			"GenerateRootCA",
			pki.NewState(),
			pki.NewSpec(),
			pki.Cert{Type: pki.CertTypeRoot},
			[]string{},
			"password\npassword\n" +
				"RootCA\n\n\n\n\n\n\n\n\n\n\n",
		},
		{
			"GenerateRootCA",
			pki.NewState(),
			&pki.Spec{
				Root: pki.Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Organization: []string{"Milad"},
				},
			},
			pki.Cert{Type: pki.CertTypeRoot},
			[]string{},
			"password\npassword\n" +
				"RootCA\n\n\n\n\n\n\n\n\n",
		},
		{
			"GenerateIntermediateCA",
			pki.NewState(),
			pki.NewSpec(),
			pki.Cert{Type: pki.CertTypeInterm},
			[]string{},
			"sre\n" +
				"password\npassword\n" +
				"SRE CA\n\n\n\n\n\n\n\n\n\n\n",
		},
		{
			"GenerateIntermediateCA",
			pki.NewState(),
			&pki.Spec{
				Interm: pki.Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Organization: []string{"Milad"},
				},
			},
			pki.Cert{Type: pki.CertTypeInterm},
			[]string{"-name=sre"},
			"password\npassword\n" +
				"SRE CA\nOttawa\nSRE\n\n\n\n\n\n",
		},
		{
			"GenerateServerCert",
			pki.NewState(),
			pki.NewSpec(),
			pki.Cert{Type: pki.CertTypeServer},
			[]string{},
			"webapp\n" +
				"webapp.com\n\n\n\n\n\n\n\n\n\n\n",
		},
		{
			"GenerateServerCert",
			pki.NewState(),
			&pki.Spec{
				Server: pki.Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Locality:     []string{"Toronto"},
					Organization: []string{"Milad"},
				},
			},
			pki.Cert{Type: pki.CertTypeServer},
			[]string{"-name=webapp"},
			"webapp.com\nR&D\nwebapp.com\n\n\n\n\n",
		},
		{
			"GenerateClientCert",
			pki.NewState(),
			pki.NewSpec(),
			pki.Cert{Type: pki.CertTypeClient},
			[]string{},
			"myservice\n" +
				"MyService\n\n\n\n\n\n\n\n\n\n\n",
		},
		{
			"GenerateClientCert",
			pki.NewState(),
			&pki.Spec{
				Client: pki.Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Locality:     []string{"Montreal"},
					Organization: []string{"Milad"},
				},
			},
			pki.Cert{Type: pki.CertTypeClient},
			[]string{"-name=myservice"},
			"MyService\nR&D\n\n\n\n\n\n",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			err := pki.NewWorkspace(test.state, test.spec)
			assert.NoError(t, err)

			r := strings.NewReader(test.input)
			mockUI := help.NewMockUI(r)
			cmd := &ReqCommand{
				ui:  mockUI,
				pki: &help.MockManager{},
				c:   test.c,
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
		c            pki.Cert
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
			pki.Cert{Type: pki.CertTypeRoot},
			[]string{"-invalid"},
			"",
			nil,
			nil,
			ErrorInvalidFlag,
		},
		{
			"NoName",
			pki.NewState(),
			pki.NewSpec(),
			pki.Cert{},
			[]string{},
			"",
			nil,
			nil,
			ErrorInvalidName,
		},
		{
			"NoState",
			nil,
			nil,
			pki.Cert{},
			[]string{},
			"sre\n",
			nil,
			nil,
			ErrorReadState,
		},
		{
			"NoSpec",
			pki.NewState(),
			nil,
			pki.Cert{},
			[]string{"-name=sre"},
			"",
			nil,
			nil,
			ErrorReadSpec,
		},
		{
			"NoCert",
			pki.NewState(),
			pki.NewSpec(),
			pki.Cert{},
			[]string{"-name", "sre"},
			"",
			nil,
			nil,
			ErrorInvalidCert,
		},
		{
			"NoPasswordForCA",
			pki.NewState(),
			pki.NewSpec(),
			pki.Cert{Type: pki.CertTypeRoot},
			[]string{},
			"",
			errors.New("error"),
			nil,
			ErrorEnterConfig,
		},
		{
			"NoClaimEntered",
			pki.NewState(),
			pki.NewSpec(),
			pki.Cert{Type: pki.CertTypeRoot},
			[]string{},
			"password\npassword\n",
			errors.New("error"),
			nil,
			ErrorEnterClaim,
		},
		{
			"GenCertFails",
			pki.NewState(),
			pki.NewSpec(),
			pki.Cert{Type: pki.CertTypeRoot},
			[]string{},
			"password\npassword\n" +
				"\n\n\n\n\n\n\n\n\n\n\n",
			errors.New("error"),
			nil,
			ErrorCert,
		},
		{
			"GenCSRFails",
			pki.NewState(),
			pki.NewSpec(),
			pki.Cert{Type: pki.CertTypeInterm},
			[]string{"-name=sre"},
			"password\npassword\n" +
				"\n\n\n\n\n\n\n\n\n\n\n",
			nil,
			errors.New("error"),
			ErrorCSR,
		},
		{
			"GenCSRFails",
			pki.NewState(),
			pki.NewSpec(),
			pki.Cert{Type: pki.CertTypeServer},
			[]string{"-name=webapp"},
			"\n\n\n\n\n\n\n\n\n\n\n",
			nil,
			errors.New("error"),
			ErrorCSR,
		},
		{
			"GenCSRFails",
			pki.NewState(),
			pki.NewSpec(),
			pki.Cert{Type: pki.CertTypeClient},
			[]string{"-name", "myservice"},
			"\n\n\n\n\n\n\n\n\n\n\n",
			nil,
			errors.New("error"),
			ErrorCSR,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			err := pki.NewWorkspace(test.state, test.spec)
			assert.NoError(t, err)

			r := strings.NewReader(test.input)
			mockUI := help.NewMockUI(r)
			cmd := &ReqCommand{
				ui: mockUI,
				pki: &help.MockManager{
					GenCertError: test.GenCertError,
					GenCSRError:  test.GenCSRError,
				},
				c: test.c,
			}

			exit := cmd.Run(test.args)
			assert.Equal(t, test.expectedExit, exit)

			err = pki.CleanupWorkspace()
			assert.NoError(t, err)
		})
	}
}
