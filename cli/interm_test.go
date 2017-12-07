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

func TestIntermNewCommandError(t *testing.T) {
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
			[]string{"-req=interm"},
			ErrorReadState,
		},
		{
			pki.NewState(),
			nil,
			[]string{"-req=interm"},
			ErrorReadSpec,
		},
	}

	for _, test := range tests {
		cleanup, err := mockWorkspace(test.state, test.spec)
		assert.NoError(t, err)

		cmd := &IntermNewCommand{
			ui:  cli.NewMockUi(),
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

func TestIntermNewCommand(t *testing.T) {
	tests := []struct {
		state            *pki.State
		spec             *pki.Spec
		args             []string
		input            string
		GenIntermCAError error
		expectedExit     int
	}{
		{
			&pki.State{},
			&pki.Spec{},
			[]string{"-req=interm"},
			``,
			errors.New("error"),
			ErrorIntermCA,
		},
		{
			pki.NewState(),
			&pki.Spec{},
			[]string{"-req=interm"},
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
			[]string{"-req", "interm"},
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
				GenIntermCAError: test.GenIntermCAError,
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
