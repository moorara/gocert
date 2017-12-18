package cli

import (
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
	"github.com/stretchr/testify/assert"
)

func TestNewVerifyCommand(t *testing.T) {
	tests := []struct {
		expectedSynopsis string
	}{
		{
			"Verifies a certificate using its certificate authority.",
		},
		{
			"Verifies a certificate using its certificate authority.",
		},
	}

	for _, test := range tests {
		cmd := NewVerifyCommand()

		assert.Equal(t, newColoredUI(), cmd.ui)
		assert.Equal(t, pki.NewX509Manager(), cmd.pki)

		assert.Equal(t, test.expectedSynopsis, cmd.Synopsis())
		assert.NotEmpty(t, cmd.Help())
	}
}

func TestVerifyCommand(t *testing.T) {
	tests := []struct {
		title string
		args  []string
		input string
	}{
		{
			"RootVerifiesIntermediate",
			[]string{},
			`root
			interm
			`,
		},
		{
			"RootVerifiesIntermediate",
			[]string{"-ca=root"},
			`interm
			`,
		},
		{
			"RootVerifiesIntermediate",
			[]string{"-ca=root", "--name", "interm"},
			``,
		},
		{
			"IntermediateVerifiesServerClient",
			[]string{"-ca=interm"},
			`server
			client`,
		},
		{
			"IntermediateVerifiesServerClient",
			[]string{"-ca=interm", "--name", "server,client"},
			``,
		},
		{
			"IntermediateVerifiesServerWithDNS",
			[]string{"-ca=interm", "-name=server", "-dns=example.com"},
			``,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			mockUI := cli.NewMockUi()
			mockUI.InputReader = strings.NewReader(test.input)

			cmd := &VerifyCommand{
				ui:  mockUI,
				pki: &mockedManager{},
			}

			exit := cmd.Run(test.args)
			assert.Zero(t, exit)
		})
	}
}

func TestVerifyCommandError(t *testing.T) {
	tests := []struct {
		title           string
		args            []string
		input           string
		VerifyCertError error
		expectedExit    int
	}{
		{
			"InvalidFlag",
			[]string{"-invalid"},
			``,
			nil,
			ErrorInvalidFlag,
		},
		{
			"NoCAName",
			[]string{},
			``,
			nil,
			ErrorInvalidName,
		},
		{
			"NoCertName",
			[]string{},
			`root
			`,
			nil,
			ErrorInvalidName,
		},
		{
			"VerifyCertError",
			[]string{"-ca=root", "-name=interm"},
			``,
			errors.New("error"),
			ErrorVerify,
		},
	}

	err := pki.NewWorkspace(nil, nil)
	defer pki.CleanupWorkspace()
	assert.NoError(t, err)

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			mockUI := cli.NewMockUi()
			mockUI.InputReader = strings.NewReader(test.input)

			cmd := &VerifyCommand{
				ui: mockUI,
				pki: &mockedManager{
					VerifyCertError: test.VerifyCertError,
				},
			}

			exit := cmd.Run(test.args)
			assert.Equal(t, test.expectedExit, exit)
		})
	}
}
