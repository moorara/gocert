package pki

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/moorara/go-box/util"
	"github.com/stretchr/testify/assert"
)

func TestNewWorkspace(t *testing.T) {
	tests := []struct {
		state         *State
		spec          *Spec
		expectError   bool
		expectedState string
		expectedSpec  string
	}{
		{
			&State{},
			&Spec{},
			false,
			`root:
				serial: 0
				length: 0
				days: 0
			intermediate:
				serial: 0
				length: 0
				days: 0
			server:
				serial: 0
				length: 0
				days: 0
			client:
				serial: 0
				length: 0
				days: 0
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
			NewState(),
			&Spec{},
			false,
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
			&State{
				Root: ConfigCA{
					Config: Config{
						Serial: 10,
						Length: 4096,
						Days:   7300,
					},
				},
				Interm: ConfigCA{
					Config: Config{
						Serial: 100,
						Length: 4096,
						Days:   3650,
					},
				},
				Server: Config{
					Serial: 1000,
					Length: 2048,
					Days:   375,
				},
				Client: Config{
					Serial: 10000,
					Length: 2048,
					Days:   40,
				},
			},
			&Spec{
				Root: Claim{
					Country:      []string{"CA"},
					Organization: []string{"Milad"},
				},
				Interm: Claim{
					Country:      []string{"CA"},
					Locality:     []string{"Ottawa"},
					Organization: []string{"Milad"},
				},
				Server: Claim{
					Country:      []string{"CA"},
					Organization: []string{"Milad"},
				},
				Client: Claim{
					Country:      []string{"CA"},
					Organization: []string{"Mona"},
				},
				RootPolicy: Policy{
					Match:    []string{"Organization"},
					Supplied: []string{"CommonName"},
				},
				IntermPolicy: Policy{
					Match:    []string{},
					Supplied: []string{"CommonName"},
				},
			},
			false,
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
				locality = ["Ottawa"]
				organization = ["Milad"]

			[server]
				country = ["CA"]
				organization = ["Milad"]

			[client]
				country = ["CA"]
				organization = ["Mona"]

			[root_policy]
				match = ["Organization"]
				supplied = ["CommonName"]

			[intermediate_policy]
				match = []
				supplied = ["CommonName"]
			`,
		},
	}

	for _, test := range tests {
		err := NewWorkspace(test.state, test.spec)

		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)

			// Verify state file
			stateData, err := ioutil.ReadFile(FileState)
			assert.NoError(t, err)
			stateYAML := strings.Replace(test.expectedState, "\t\t\t", "", -1)
			stateYAML = strings.Replace(stateYAML, "\t", "  ", -1)
			assert.Equal(t, stateYAML, string(stateData))

			// Verify spec file
			specData, err := ioutil.ReadFile(FileSpec)
			assert.NoError(t, err)
			specTOML := strings.Replace(test.expectedSpec, "\t\t\t", "", -1)
			specTOML = strings.Replace(specTOML, "\t", "  ", -1)
			assert.Equal(t, specTOML, string(specData))
		}

		util.DeleteAll(
			"",
			DirRoot, DirInterm, DirServer, DirClient, DirCSR,
			FileState, FileSpec,
		)
	}
}
