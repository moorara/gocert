package pki

import (
	"net"
	"os"
	"testing"

	"github.com/moorara/gocert/util"
	"github.com/stretchr/testify/assert"
)

type (
	stateLoadTest struct {
		fixture       string
		expectError   bool
		expectedState *State
	}

	stateSaveTest struct {
		state           *State
		expectedFixture string
	}

	specLoadTest struct {
		fixture      string
		expectError  bool
		expectedSpec *Spec
	}

	specSaveTest struct {
		spec            *Spec
		expectedFixture string
	}
)

var (
	loadTests = []struct {
		state stateLoadTest
		spec  specLoadTest
	}{
		{
			stateLoadTest{
				fixture:       "./fixture/load/invalid.yaml",
				expectError:   true,
				expectedState: nil,
			},
			specLoadTest{
				fixture:      "./fixture/load/invalid.toml",
				expectError:  true,
				expectedSpec: nil,
			},
		},
		{
			stateLoadTest{
				fixture:       "./fixture/load/empty.yaml",
				expectError:   false,
				expectedState: &State{},
			},
			specLoadTest{
				fixture:      "./fixture/load/empty.toml",
				expectError:  false,
				expectedSpec: &Spec{},
			},
		},
		{
			stateLoadTest{
				fixture:     "./fixture/load/simple.yaml",
				expectError: false,
				expectedState: &State{
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
				fixture:     "./fixture/load/simple.toml",
				expectError: false,
				expectedSpec: &Spec{
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
				fixture:     "./fixture/load/complex.yaml",
				expectError: false,
				expectedState: &State{
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
				fixture:     "./fixture/load/complex.toml",
				expectError: false,
				expectedSpec: &Spec{
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
				state:           nil,
				expectedFixture: "./fixture/save/empty.yaml",
			},
			specSaveTest{
				spec:            nil,
				expectedFixture: "./fixture/save/empty.toml",
			},
		},
		{
			stateSaveTest{
				state:           &State{},
				expectedFixture: "./fixture/save/zero.yaml",
			},
			specSaveTest{
				spec:            &Spec{},
				expectedFixture: "./fixture/save/zero.toml",
			},
		},
		{
			stateSaveTest{
				state:           NewState(),
				expectedFixture: "./fixture/save/default.yaml",
			},
			specSaveTest{
				spec:            NewSpec(),
				expectedFixture: "./fixture/save/default.toml",
			},
		},
		{
			stateSaveTest{
				state: &State{
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
				expectedFixture: "./fixture/save/custom1.yaml",
			},
			specSaveTest{
				spec: &Spec{
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
				expectedFixture: "./fixture/save/custom1.toml",
			},
		},
		{
			stateSaveTest{
				state: &State{
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
				expectedFixture: "./fixture/save/custom2.yaml",
			},
			specSaveTest{
				spec: &Spec{
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
				expectedFixture: "./fixture/save/custom2.toml",
			},
		},
	}
)

func verifyStateFile(t *testing.T, expectedFixture, stateFile string) {
	expectedStateYAML, err := os.ReadFile(expectedFixture)
	assert.NoError(t, err)

	if len(expectedStateYAML) > 0 {
		stateYAML, err := os.ReadFile(stateFile)
		assert.NoError(t, err)

		assert.Equal(t, string(expectedStateYAML), string(stateYAML))
	}
}

func verifySpecFile(t *testing.T, expectedFixture, specFile string) {
	expectedSpecTOML, err := os.ReadFile(expectedFixture)
	assert.NoError(t, err)

	if len(expectedSpecTOML) > 0 {
		specTOML, err := os.ReadFile(specFile)
		assert.NoError(t, err)

		assert.Equal(t, string(expectedSpecTOML), string(specTOML))
	}
}

func TestLoadState(t *testing.T) {
	for _, test := range loadTests {
		stateYAML, err := os.ReadFile(test.state.fixture)
		assert.NoError(t, err)

		file, delete, err := util.CreateTempFile(string(stateYAML))
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

		verifyStateFile(t, test.state.expectedFixture, file)
	}
}

func TestLoadSpec(t *testing.T) {
	for _, test := range loadTests {
		spec, err := LoadSpec(test.spec.fixture)

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

		verifySpecFile(t, test.spec.expectedFixture, file)
	}
}

func TestNewWorkspace(t *testing.T) {
	for _, test := range saveTests {
		err := NewWorkspace(test.state.state, test.spec.spec)
		assert.NoError(t, err)

		verifyStateFile(t, test.state.expectedFixture, FileState)
		verifySpecFile(t, test.spec.expectedFixture, FileSpec)

		assert.NoError(t, CleanupWorkspace())
	}
}

func TestLoadWorkspace(t *testing.T) {
	for _, test := range loadTests {
		stateYAML, err := os.ReadFile(test.state.fixture)
		assert.NoError(t, err)
		err = os.WriteFile(FileState, stateYAML, 0644)
		assert.NoError(t, err)

		specTOML, err := os.ReadFile(test.spec.fixture)
		assert.NoError(t, err)
		err = os.WriteFile(FileSpec, specTOML, 0644)
		assert.NoError(t, err)

		state, spec, err := LoadWorkspace()

		if test.state.expectError || test.spec.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.state.expectedState, state)
			assert.Equal(t, test.spec.expectedSpec, spec)
		}

		assert.NoError(t, CleanupWorkspace())
	}
}

func TestSaveWorkspace(t *testing.T) {
	for _, test := range saveTests {
		err := SaveWorkspace(test.state.state, test.spec.spec)
		assert.NoError(t, err)

		verifyStateFile(t, test.state.expectedFixture, FileState)
		verifySpecFile(t, test.spec.expectedFixture, FileSpec)

		assert.NoError(t, CleanupWorkspace())
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

		assert.NoError(t, CleanupWorkspace())
	}
}
