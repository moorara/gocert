package cli

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
	"github.com/stretchr/testify/assert"
)

func TestNewInitCommand(t *testing.T) {
	cmd := NewInitCommand()

	assert.Equal(t, newColoredUI(), cmd.ui)
}

func TestInitCommand(t *testing.T) {
	tests := []struct {
		title         string
		args          []string
		input         string
		expectedExit  int
		expectedState string
		expectedSpec  string
	}{
		{
			"DefaultStateSpec",
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

			[root_policy]

			[intermediate_policy]
			`,
		},
		{
			"CustomStateSpec",
			[]string{},
			`CA
			Ontario
			Ottawa
			Milad












			Ops
			example.com




			R&D
			example.org




			SRE





			Organization
			CommonName,OrganizationalUnit
			Organization
			CommonName
			`,
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
				country = ["CA"]
				province = ["Ontario"]
				locality = ["Ottawa"]
				organization = ["Milad"]

			[intermediate]
				country = ["CA"]
				province = ["Ontario"]
				locality = ["Ottawa"]
				organization = ["Milad"]
				organizational_unit = ["Ops"]
				dns_name = ["example.com"]

			[server]
				country = ["CA"]
				province = ["Ontario"]
				locality = ["Ottawa"]
				organization = ["Milad"]
				organizational_unit = ["R&D"]
				dns_name = ["example.org"]

			[client]
				country = ["CA"]
				province = ["Ontario"]
				locality = ["Ottawa"]
				organization = ["Milad"]
				organizational_unit = ["SRE"]

			[root_policy]
				match = ["Organization"]
				supplied = ["CommonName", "OrganizationalUnit"]

			[intermediate_policy]
				match = ["Organization"]
				supplied = ["CommonName"]
			`,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			mockUI := cli.NewMockUi()
			mockUI.InputReader = strings.NewReader(test.input)

			cmd := &InitCommand{
				ui: mockUI,
			}
			exit := cmd.Run(test.args)

			assert.Equal(t, initSynopsis, cmd.Synopsis())
			assert.Equal(t, initHelp, cmd.Help())
			assert.Equal(t, test.expectedExit, exit)

			// Verify state file
			stateData, err := ioutil.ReadFile(pki.FileState)
			assert.NoError(t, err)
			stateYAML := strings.Replace(test.expectedState, "\t\t\t", "", -1)
			stateYAML = strings.Replace(stateYAML, "\t", "  ", -1)
			assert.Equal(t, stateYAML, string(stateData))

			// Verify spec file
			specData, err := ioutil.ReadFile(pki.FileSpec)
			assert.NoError(t, err)
			specTOML := strings.Replace(test.expectedSpec, "\t\t\t", "", -1)
			specTOML = strings.Replace(specTOML, "\t", "  ", -1)
			assert.Equal(t, specTOML, string(specData))

			err = pki.CleanupWorkspace()
			assert.NoError(t, err)
		})
	}
}
