package cli

import (
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
	"github.com/stretchr/testify/assert"
)

func TestNewIntermNewCommand(t *testing.T) {
	cmd := NewIntermNewCommand()

	assert.Equal(t, newColoredUI(), cmd.ui)
	assert.Equal(t, pki.NewX509Manager(), cmd.pki)
}

func TestIntermNewCommand(t *testing.T) {
	tests := []struct {
		state             *pki.State
		spec              *pki.Spec
		args              []string
		input             string
		GenIntermCSRError error
		expectedExit      int
	}{
		{
			&pki.State{},
			&pki.Spec{},
			[]string{"-req=it"},
			``,
			errors.New("error"),
			ErrorIntermCA,
		},
		{
			pki.NewState(),
			&pki.Spec{},
			[]string{"-req=ops"},
			``,
			nil,
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
			[]string{"-req", "ops"},
			`secret
			IntermediateCA
			Ottawa,Toronto
			R&D
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

		cmd := &IntermNewCommand{
			ui: mockUI,
			pki: &mockedManager{
				GenIntermCSRError: test.GenIntermCSRError,
			},
		}
		exit := cmd.Run(test.args)

		assert.Equal(t, intermNewSynopsis, cmd.Synopsis())
		assert.Equal(t, intermNewHelp, cmd.Help())
		assert.Equal(t, test.expectedExit, exit)

		err = cleanup()
		assert.NoError(t, err)
	}
}

func TestIntermNewCommandError(t *testing.T) {
	tests := []struct {
		state        *pki.State
		spec         *pki.Spec
		args         []string
		input        string
		expectedExit int
	}{
		{
			nil,
			nil,
			[]string{"-invalid"},
			``,
			ErrorInvalidFlag,
		},
		{
			nil,
			&pki.Spec{},
			[]string{},
			``,
			ErrorReadState,
		},
		{
			pki.NewState(),
			nil,
			[]string{},
			``,
			ErrorReadSpec,
		},
		{
			pki.NewState(),
			pki.NewSpec(),
			[]string{},
			``,
			ErrorNoName,
		},
	}

	for _, test := range tests {
		cleanup, err := mockWorkspace(test.state, test.spec)
		assert.NoError(t, err)

		mockUI := cli.NewMockUi()
		mockUI.InputReader = strings.NewReader(test.input)

		cmd := &IntermNewCommand{
			ui:  mockUI,
			pki: &mockedManager{},
		}
		exit := cmd.Run(test.args)

		assert.Equal(t, intermNewSynopsis, cmd.Synopsis())
		assert.Equal(t, intermNewHelp, cmd.Help())
		assert.Equal(t, test.expectedExit, exit)

		err = cleanup()
		assert.NoError(t, err)
	}
}
