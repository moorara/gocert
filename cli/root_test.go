package cli

import (
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
	"github.com/stretchr/testify/assert"
)

func TestNewRootNewCommand(t *testing.T) {
	cmd := NewRootNewCommand()

	assert.Equal(t, newColoredUI(), cmd.ui)
	assert.Equal(t, pki.NewX509Manager(), cmd.pki)
}

func TestRootNewCommand(t *testing.T) {
	tests := []struct {
		state          *pki.State
		spec           *pki.Spec
		args           []string
		input          string
		GenRootCAError error
		expectedExit   int
	}{
		{
			&pki.State{},
			&pki.Spec{},
			[]string{},
			``,
			errors.New("error"),
			ErrorRootCA,
		},
		{
			pki.NewState(),
			&pki.Spec{},
			[]string{},
			``,
			nil,
			0,
		},
		{
			pki.NewState(),
			&pki.Spec{
				Root: pki.Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Organization: []string{"Milad"},
				},
			},
			[]string{},
			`secret
			RootCA
			Ottawa,Toronto
			`,
			nil,
			0,
		},
	}

	for _, test := range tests {
		cleanup, err := mockWorkspace(test.state, test.spec)
		assert.NoError(t, err)

		mockUI := cli.NewMockUi()
		mockUI.InputReader = strings.NewReader(test.input)

		cmd := &RootNewCommand{
			ui: mockUI,
			pki: &mockedManager{
				GenRootCAError: test.GenRootCAError,
			},
		}
		exit := cmd.Run(test.args)

		assert.Equal(t, rootSynopsis, cmd.Synopsis())
		assert.Equal(t, rootHelp, cmd.Help())
		assert.Equal(t, test.expectedExit, exit)

		err = cleanup()
		assert.NoError(t, err)
	}
}

func TestRootNewCommandError(t *testing.T) {
	tests := []struct {
		state        *pki.State
		spec         *pki.Spec
		args         []string
		expectedExit int
	}{
		{
			nil,
			nil,
			[]string{"-invalid"},
			ErrorInvalidFlag,
		},
		{
			nil,
			&pki.Spec{},
			[]string{},
			ErrorReadState,
		},
		{
			pki.NewState(),
			nil,
			[]string{},
			ErrorReadSpec,
		},
	}

	for _, test := range tests {
		cleanup, err := mockWorkspace(test.state, test.spec)
		assert.NoError(t, err)

		cmd := &RootNewCommand{
			ui:  cli.NewMockUi(),
			pki: &mockedManager{},
		}
		exit := cmd.Run(test.args)

		assert.Equal(t, rootSynopsis, cmd.Synopsis())
		assert.Equal(t, rootHelp, cmd.Help())
		assert.Equal(t, test.expectedExit, exit)

		err = cleanup()
		assert.NoError(t, err)
	}
}
