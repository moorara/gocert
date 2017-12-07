package cli

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/moorara/go-box/util"
	"github.com/moorara/gocert/pki"
	"github.com/stretchr/testify/assert"
)

func TestInitCommand(t *testing.T) {
	tests := []struct {
		args          []string
		input         string
		expectedExit  int
		expectedState string
		expectedSpec  string
	}{
		{
			[]string{},
			``,
			0,
			`root:
				serial: 10
				length: 4096
				days: 7300
			intermediate:
				serial: 100
				length: 4096
				days: 3650
			server:
				serial: 1000
				length: 2048
				days: 375
			client:
				serial: 10000
				length: 2048
				days: 40
			`,
			`[root]

			[intermediate]

			[server]

			[client]
			`,
		},
	}

	for _, test := range tests {
		mockUI := cli.NewMockUi()
		mockUI.InputReader = strings.NewReader(test.input)

		cmd := NewInitCommand()
		cmd.ui = mockUI
		exit := cmd.Run(test.args)

		assert.Equal(t, initSynopsis, cmd.Synopsis())
		assert.Equal(t, initHelp, cmd.Help())
		assert.Equal(t, test.expectedExit, exit)

		// Verify state file
		stateData, err := ioutil.ReadFile(pki.FileState)
		assert.NoError(t, err)
		stateYAML := strings.Replace(test.expectedState, "\t\t\t\t", "  ", -1)
		stateYAML = strings.Replace(stateYAML, "\t\t\t", "", -1)
		assert.Equal(t, stateYAML, string(stateData))

		// Verify spec file
		specData, err := ioutil.ReadFile(pki.FileSpec)
		assert.NoError(t, err)
		specTOML := strings.Replace(test.expectedSpec, "\t\t\t\t", "  ", -1)
		specTOML = strings.Replace(specTOML, "\t\t\t", "", -1)
		assert.Equal(t, specTOML, string(specData))

		util.DeleteAll(
			"",
			pki.DirRoot,
			pki.DirInterm,
			pki.DirServer,
			pki.DirClient,
			pki.DirCSR,
			pki.FileState,
			pki.FileSpec,
		)
	}
}
