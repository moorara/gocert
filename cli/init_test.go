package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/moorara/gocert/pki"
	"github.com/stretchr/testify/assert"
)

func TestNewInitCommand(t *testing.T) {
	tests := []struct {
		expectedSynopsis string
	}{
		{
			"Initializes a new workspace with desired configs and specs.",
		},
		{
			"Initializes a new workspace with desired configs and specs.",
		},
	}

	for _, test := range tests {
		cmd := NewInitCommand()

		assert.Equal(t, newColoredUI(), cmd.ui)
		assert.Equal(t, test.expectedSynopsis, cmd.Synopsis())
		assert.NotEmpty(t, cmd.Help())
	}
}

func TestInitCommand(t *testing.T) {
	tests := []struct {
		title         string
		args          []string
		input         string
		expectedState string
		expectedSpec  string
	}{
		{
			"DefaultStateSpec",
			[]string{},
			"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n" +
				"\n\n",
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
				supplied = ["CommonName"]

			[intermediate_policy]
				supplied = ["CommonName"]

			[metadata]
			`,
		},
		{
			"CustomStateSpec",
			[]string{},
			"CA\n\n\nMilad\n\n\n\n\n-\n-\n" +
				"\n\n\n-\n-\n-\n" +
				"\n\nSRE\n-\n-\n-\n" +
				"Ontario\nOttawa\nR&D\nexample.org\n127.0.0.1\n\n" +
				"Ontario\nOttawa\nQE\n\n\n\n" +
				"Organization\nCommonName,OrganizationalUnit\n" +
				"Organization\nCommonName\n",
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
				organization = ["Milad"]

			[intermediate]
				country = ["CA"]
				organization = ["Milad"]
				organizational_unit = ["SRE"]

			[server]
				country = ["CA"]
				province = ["Ontario"]
				locality = ["Ottawa"]
				organization = ["Milad"]
				organizational_unit = ["R&D"]
				dns_name = ["example.org"]
				ip_address = ["127.0.0.1"]

			[client]
				country = ["CA"]
				province = ["Ontario"]
				locality = ["Ottawa"]
				organization = ["Milad"]
				organizational_unit = ["QE"]

			[root_policy]
				match = ["Organization"]
				supplied = ["CommonName", "OrganizationalUnit"]

			[intermediate_policy]
				match = ["Organization"]
				supplied = ["CommonName"]

			[metadata]
				clientSkip = ["Claim.StreetAddress", "Claim.PostalCode"]
				intermSkip = ["Claim.StreetAddress", "Claim.PostalCode", "Claim.DNSName", "Claim.IPAddress", "Claim.EmailAddress"]
				rootSkip = ["Claim.StreetAddress", "Claim.PostalCode", "Claim.DNSName", "Claim.IPAddress", "Claim.EmailAddress"]
				serverSkip = ["Claim.StreetAddress", "Claim.PostalCode"]
			`,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			r := strings.NewReader(test.input)
			mockUI := newMockUI(r)
			cmd := &InitCommand{
				ui: mockUI,
			}

			exit := cmd.Run(test.args)
			assert.Zero(t, exit)

			// Verify state file
			stateData, err := os.ReadFile(pki.FileState)
			assert.NoError(t, err)
			stateYAML := strings.Replace(test.expectedState, "\t\t\t", "", -1)
			stateYAML = strings.Replace(stateYAML, "\t", "  ", -1)
			assert.Equal(t, stateYAML, string(stateData))

			// Verify spec file
			specData, err := os.ReadFile(pki.FileSpec)
			assert.NoError(t, err)
			specTOML := strings.Replace(test.expectedSpec, "\t\t\t", "", -1)
			specTOML = strings.Replace(specTOML, "\t", "  ", -1)
			assert.Equal(t, specTOML, string(specData))

			err = pki.CleanupWorkspace()
			assert.NoError(t, err)
		})
	}
}

func TestInitCommandError(t *testing.T) {
	tests := []struct {
		title        string
		args         []string
		input        string
		expectedExit int
	}{
		{
			"InvalidFlag",
			[]string{"-invalid"},
			``,
			ErrorInvalidFlag,
		},
		{
			"NoInputForSpec",
			[]string{""},
			``,
			ErrorEnterSpec,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			r := strings.NewReader(test.input)
			mockUI := newMockUI(r)
			cmd := &InitCommand{
				ui: mockUI,
			}

			exit := cmd.Run(test.args)
			assert.Equal(t, test.expectedExit, exit)
		})
	}
}
