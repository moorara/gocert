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
		args  []string
		input string
	}{
		{
			[]string{},
			`root
			interm
			`,
		},
		{
			[]string{"-ca=root"},
			`interm
			`,
		},
		{
			[]string{"-ca=root", "--name", "interm"},
			``,
		},
		{
			[]string{"-ca=interm"},
			`server
			client`,
		},
		{
			[]string{"-ca=interm", "--name", "server,client"},
			``,
		},
	}

	for _, test := range tests {
		mockUI := cli.NewMockUi()
		mockUI.InputReader = strings.NewReader(test.input)

		cmd := &VerifyCommand{
			ui:  mockUI,
			pki: &mockedManager{},
		}

		exit := cmd.Run(test.args)
		assert.Zero(t, exit)
	}
}

func TestVerifyCommandError(t *testing.T) {
	tests := []struct {
		args            []string
		input           string
		VerifyCertError error
		expectedExit    int
	}{
		{
			[]string{"-invalid"},
			``,
			nil,
			ErrorInvalidFlag,
		},
		{
			[]string{},
			``,
			nil,
			ErrorInvalidName,
		},
		{
			[]string{},
			`root
			`,
			nil,
			ErrorInvalidName,
		},
		{
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
	}
}
