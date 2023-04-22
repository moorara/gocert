package pki

import (
	"net"
	"os"
	"strings"
	"testing"

	"github.com/moorara/gocert/util"
	"github.com/stretchr/testify/assert"
)

type (
	stateLoadTest struct {
		yaml          string
		expectError   bool
		expectedState *State
	}

	stateSaveTest struct {
		state        *State
		expectedYAML string
	}

	specLoadTest struct {
		toml         string
		expectError  bool
		expectedSpec *Spec
	}

	specSaveTest struct {
		spec         *Spec
		expectedTOML string
	}
)

var (
	loadTests = []struct {
		state stateLoadTest
		spec  specLoadTest
	}{
		{
			stateLoadTest{
				`invalid yaml`,
				true,
				nil,
			},
			specLoadTest{
				`invalid toml`,
				true,
				nil,
			},
		},
		{
			stateLoadTest{
				``,
				false,
				&State{},
			},
			specLoadTest{
				``,
				false,
				&Spec{},
			},
		},
		{
			stateLoadTest{
				`
				root:
					serial: 10
					length: 4096
					days: 7300
				intermediate:
					serial: 100
					length: 4096
					days: 3650
				`,
				false,
				&State{
					Root: Config{
						Serial: 10,
						Length: 4096,
						Days:   7300,
					},
					Interm: Config{
						Serial: 100,
						Length: 4096,
						Days:   3650,
					},
				},
			},
			specLoadTest{
				`
				[root]
					country = [ "CA" ]
					organization = [ "Moorara" ]
				[server]
					country = [ "US" ]
					organization = [ "Milad" ]
					dns_name = [ "example.com" ]
					email_address = [ "milad@example.com" ]
				`,
				false,
				&Spec{
					Root: Claim{
						Country:      []string{"CA"},
						Organization: []string{"Moorara"},
					},
					Server: Claim{
						Country:      []string{"US"},
						Organization: []string{"Milad"},
						DNSName:      []string{"example.com"},
						EmailAddress: []string{"milad@example.com"},
					},
				},
			},
		},
		{
			stateLoadTest{
				`
				root:
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
				false,
				&State{
					Root: Config{
						Serial: 10,
						Length: 4096,
						Days:   7300,
					},
					Interm: Config{
						Serial: 100,
						Length: 4096,
						Days:   3650,
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
			},
			specLoadTest{
				`
				[root]
					country = [ "CA", "US" ]
					province = [ "Ontario", "Massachusetts" ]
					locality = [ "Ottawa", "Boston" ]
					organization = [ "Moorara" ]
				[intermediate]
					country = [ "CA" ]
					province = [ "Ontario" ]
					locality = [ "Ottawa" ]
					organization = [ "Moorara" ]
					email_address = [ "moorara@example.com" ]
				[server]
					country = [ "US" ]
					province = [ "Virginia" ]
					locality = [ "Richmond" ]
					organization = [ "Moorara" ]
					dns_name = [ "example.com" ]
					ip_address = [ "127.0.0.1" ]
					email_address = [ "moorara@example.com" ]
				[client]
					country = [ "UK" ]
					locality = [ "London" ]
					organization = [ "Moorara" ]
					email_address = [ "moorara@example.com" ]
				[root_policy]
					match = ["Country", "Organization"]
					supplied = ["CommonName"]
				[intermediate_policy]
					match = ["Organization"]
					supplied = ["CommonName"]
				[metadata]
					RootSkip = ["IPAddress", "StreetAddress", "PostalCode"]
					IntermSkip = ["IPAddress", "StreetAddress", "PostalCode"]
					ServerSkip = ["StreetAddress", "PostalCode"]
					ClientSkip = ["StreetAddress", "PostalCode"]
				`,
				false,
				&Spec{
					Root: Claim{
						Country:      []string{"CA", "US"},
						Province:     []string{"Ontario", "Massachusetts"},
						Locality:     []string{"Ottawa", "Boston"},
						Organization: []string{"Moorara"},
					},
					Interm: Claim{
						Country:      []string{"CA"},
						Province:     []string{"Ontario"},
						Locality:     []string{"Ottawa"},
						Organization: []string{"Moorara"},
						EmailAddress: []string{"moorara@example.com"},
					},
					Server: Claim{
						Country:      []string{"US"},
						Province:     []string{"Virginia"},
						Locality:     []string{"Richmond"},
						Organization: []string{"Moorara"},
						DNSName:      []string{"example.com"},
						IPAddress:    []net.IP{net.ParseIP("127.0.0.1")},
						EmailAddress: []string{"moorara@example.com"},
					},
					Client: Claim{
						Country:      []string{"UK"},
						Locality:     []string{"London"},
						Organization: []string{"Moorara"},
						EmailAddress: []string{"moorara@example.com"},
					},
					RootPolicy: Policy{
						Match:    []string{"Country", "Organization"},
						Supplied: []string{"CommonName"},
					},
					IntermPolicy: Policy{
						Match:    []string{"Organization"},
						Supplied: []string{"CommonName"},
					},
					Metadata: Metadata{
						"RootSkip":   []string{"IPAddress", "StreetAddress", "PostalCode"},
						"IntermSkip": []string{"IPAddress", "StreetAddress", "PostalCode"},
						"ServerSkip": []string{"StreetAddress", "PostalCode"},
						"ClientSkip": []string{"StreetAddress", "PostalCode"},
					},
				},
			},
		},
	}

	saveTests = []struct {
		state stateSaveTest
		spec  specSaveTest
	}{
		{
			stateSaveTest{
				nil,
				``,
			},
			specSaveTest{
				nil,
				``,
			},
		},
		{
			stateSaveTest{
				&State{},
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
			},
			specSaveTest{
				&Spec{},
				`[root]

				[intermediate]

				[server]

				[client]

				[root_policy]

				[intermediate_policy]
				`,
			},
		},
		{
			stateSaveTest{
				NewState(),
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
			},
			specSaveTest{
				NewSpec(),
				`[root]

				[intermediate]

				[server]

				[client]

				[root_policy]
					match = []
					supplied = ["CommonName"]

				[intermediate_policy]
					match = []
					supplied = ["CommonName"]

				[metadata]
				`,
			},
		},
		{
			stateSaveTest{
				&State{
					Root: Config{
						Serial: 10,
						Length: 4096,
						Days:   7300,
					},
					Interm: Config{
						Serial: 100,
						Length: 4096,
						Days:   3650,
					},
				},
				`root:
					serial: 10
					length: 4096
					days: 7300
				intermediate:
					serial: 100
					length: 4096
					days: 3650
				server:
					serial: 0
					length: 0
					days: 0
				client:
					serial: 0
					length: 0
					days: 0
				`,
			},
			specSaveTest{
				&Spec{
					Root: Claim{
						Locality:     []string{"Ottawa"},
						Organization: []string{"Moorara"},
					},
					Interm: Claim{},
					Server: Claim{
						Country:      []string{"US"},
						Organization: []string{"Milad"}},
					Client: Claim{},
					RootPolicy: Policy{
						Match:    []string{"Organization"},
						Supplied: []string{"CommonName"},
					},
					Metadata: Metadata{},
				},
				`[root]
					locality = ["Ottawa"]
					organization = ["Moorara"]

				[intermediate]

				[server]
					country = ["US"]
					organization = ["Milad"]

				[client]

				[root_policy]
					match = ["Organization"]
					supplied = ["CommonName"]

				[intermediate_policy]

				[metadata]
				`,
			},
		},
		{
			stateSaveTest{
				&State{
					Root: Config{
						Serial: 10,
						Length: 4096,
						Days:   7300,
					},
					Interm: Config{
						Serial: 100,
						Length: 4096,
						Days:   3650,
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
			},
			specSaveTest{
				&Spec{
					Root: Claim{
						Country:      []string{"CA", "US"},
						Province:     []string{"Ontario", "Massachusetts"},
						Locality:     []string{"Ottawa", "Boston"},
						Organization: []string{"Moorara"},
					},
					Interm: Claim{
						Country:      []string{"CA"},
						Province:     []string{"Ontario"},
						Locality:     []string{"Ottawa"},
						Organization: []string{"Moorara"},
					},
					Server: Claim{
						Country:      []string{"US"},
						Province:     []string{"Virginia"},
						Locality:     []string{"Richmond"},
						Organization: []string{"Moorara"},
					},
					Client: Claim{
						Country:      []string{"UK"},
						Locality:     []string{"London"},
						Organization: []string{"Moorara"},
					},
					RootPolicy: Policy{
						Match:    []string{"Country", "Organization"},
						Supplied: []string{"CommonName"},
					},
					IntermPolicy: Policy{
						Match:    []string{"Organization"},
						Supplied: []string{"CommonName"},
					},
					Metadata: Metadata{
						"RootSkip":   []string{"IPAddress", "StreetAddress", "PostalCode"},
						"IntermSkip": []string{"IPAddress", "StreetAddress", "PostalCode"},
						"ServerSkip": []string{"StreetAddress", "PostalCode"},
						"ClientSkip": []string{"StreetAddress", "PostalCode"},
					},
				},
				`[root]
					country = ["CA", "US"]
					province = ["Ontario", "Massachusetts"]
					locality = ["Ottawa", "Boston"]
					organization = ["Moorara"]

				[intermediate]
					country = ["CA"]
					province = ["Ontario"]
					locality = ["Ottawa"]
					organization = ["Moorara"]

				[server]
					country = ["US"]
					province = ["Virginia"]
					locality = ["Richmond"]
					organization = ["Moorara"]

				[client]
					country = ["UK"]
					locality = ["London"]
					organization = ["Moorara"]

				[root_policy]
					match = ["Country", "Organization"]
					supplied = ["CommonName"]

				[intermediate_policy]
					match = ["Organization"]
					supplied = ["CommonName"]

				[metadata]
					ClientSkip = ["StreetAddress", "PostalCode"]
					IntermSkip = ["IPAddress", "StreetAddress", "PostalCode"]
					RootSkip = ["IPAddress", "StreetAddress", "PostalCode"]
					ServerSkip = ["StreetAddress", "PostalCode"]
				`,
			},
		},
	}
)

func verifyStateFile(t *testing.T, stateFile, expectedYAML string) {
	if expectedYAML == "" {
		return
	}

	stateData, err := os.ReadFile(stateFile)
	assert.NoError(t, err)

	expectedYAML = strings.Replace(expectedYAML, "\t\t\t\t", "", -1)
	expectedYAML = strings.Replace(expectedYAML, "\t", "  ", -1)

	assert.Equal(t, expectedYAML, string(stateData))
}

func verifySpecFile(t *testing.T, specFile, expectedTOML string) {
	if expectedTOML == "" {
		return
	}

	specData, err := os.ReadFile(specFile)
	assert.NoError(t, err)

	expectedTOML = strings.Replace(expectedTOML, "\t\t\t\t", "", -1)
	expectedTOML = strings.Replace(expectedTOML, "\t", "  ", -1)

	assert.Equal(t, expectedTOML, string(specData))
}

func TestLoadState(t *testing.T) {
	for _, test := range loadTests {
		yaml := strings.Replace(test.state.yaml, "\t", "  ", -1)
		file, delete, err := util.CreateTempFile(yaml)
		defer delete()
		assert.NoError(t, err)

		state, err := LoadState(file)

		if test.state.expectError {
			assert.Error(t, err)
			assert.Nil(t, state)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.state.expectedState, state)
		}
	}
}

func TestLoadStateError(t *testing.T) {
	spec, err := LoadState("")
	assert.Error(t, err)
	assert.Nil(t, spec)
}

func TestSaveState(t *testing.T) {
	for _, test := range saveTests {
		file, delete, err := util.CreateTempFile("")
		defer delete()
		assert.NoError(t, err)

		err = SaveState(test.state.state, file)
		assert.NoError(t, err)

		verifyStateFile(t, file, test.state.expectedYAML)
	}
}

func TestLoadSpec(t *testing.T) {
	for _, test := range loadTests {
		file, delete, err := util.CreateTempFile(test.spec.toml)
		defer delete()
		assert.NoError(t, err)

		spec, err := LoadSpec(file)

		if test.spec.expectError {
			assert.Error(t, err)
			assert.Nil(t, spec)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.spec.expectedSpec, spec)
		}
	}
}

func TestSaveSpec(t *testing.T) {
	for _, test := range saveTests {
		file, delete, err := util.CreateTempFile("")
		defer delete()
		assert.NoError(t, err)

		err = SaveSpec(test.spec.spec, file)
		assert.NoError(t, err)

		verifySpecFile(t, file, test.spec.expectedTOML)
	}
}

func TestNewWorkspace(t *testing.T) {
	for _, test := range saveTests {
		err := NewWorkspace(test.state.state, test.spec.spec)
		assert.NoError(t, err)

		verifyStateFile(t, FileState, test.state.expectedYAML)
		verifySpecFile(t, FileSpec, test.spec.expectedTOML)

		err = CleanupWorkspace()
		assert.NoError(t, err)
	}
}

func TestLoadWorkspace(t *testing.T) {
	for _, test := range loadTests {
		yaml := strings.Replace(test.state.yaml, "\t", "  ", -1)
		err := os.WriteFile(FileState, []byte(yaml), 0644)
		assert.NoError(t, err)
		err = os.WriteFile(FileSpec, []byte(test.spec.toml), 0644)
		assert.NoError(t, err)

		state, spec, err := LoadWorkspace()

		if test.state.expectError || test.spec.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.state.expectedState, state)
			assert.Equal(t, test.spec.expectedSpec, spec)
		}

		err = CleanupWorkspace()
		assert.NoError(t, err)
	}
}

func TestSaveWorkspace(t *testing.T) {
	for _, test := range saveTests {
		err := SaveWorkspace(test.state.state, test.spec.spec)
		assert.NoError(t, err)

		verifyStateFile(t, FileState, test.state.expectedYAML)
		verifySpecFile(t, FileSpec, test.spec.expectedTOML)

		err = CleanupWorkspace()
		assert.NoError(t, err)
	}
}

func TestCleanupWorkspace(t *testing.T) {
	tests := []struct {
		files []string
	}{
		{},
		{
			[]string{
				DirRoot + "/root.ca.key",
				DirRoot + "/root.ca.cert",
			},
		},
		{
			[]string{
				DirRoot + "/root.ca.key",
				DirRoot + "/root.ca.cert",
				DirInterm + "/interm.ca.key",
				DirInterm + "/interm.ca.cert",
				DirServer + "/webapp.ca.key",
				DirServer + "/webapp.ca.cert",
				DirClient + "/service.ca.key",
				DirClient + "/service.ca.cert",
				DirCSR + "/interm.ca.csr",
				DirCSR + "/webapp.ca.csr",
				DirCSR + "/service.ca.csr",
			},
		},
	}

	for _, test := range tests {
		// Mock directorys and files
		_, err := util.MkDirs("", DirRoot, DirInterm, DirServer, DirClient, DirCSR)
		assert.NoError(t, err)
		err = os.WriteFile(FileState, nil, 0644)
		assert.NoError(t, err)
		err = os.WriteFile(FileSpec, nil, 0644)
		assert.NoError(t, err)

		// Mock artifacts
		for _, file := range test.files {
			err = os.WriteFile(file, nil, 0644)
			assert.NoError(t, err)
		}

		err = CleanupWorkspace()
		assert.NoError(t, err)
	}
}
