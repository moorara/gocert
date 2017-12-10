package pki

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/moorara/go-box/util"
	yaml "gopkg.in/yaml.v2"
)

// LoadState reads and parses state from a YAML file
func LoadState(file string) (*State, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	state := new(State)
	err = yaml.Unmarshal(data, state)
	if err != nil {
		return nil, err
	}

	return state, nil
}

// SaveState writes state to a YAML file
func SaveState(state *State, file string) error {
	if state == nil {
		return nil
	}

	data, err := yaml.Marshal(state)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// LoadSpec reads and parses spec from a TOML file
func LoadSpec(file string) (*Spec, error) {
	spec := new(Spec)
	_, err := toml.DecodeFile(file, spec)
	if err != nil {
		return nil, err
	}

	return spec, nil
}

// SaveSpec writes spec to a TOML file
func SaveSpec(spec *Spec, file string) error {
	if spec == nil {
		return nil
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}

	err = toml.NewEncoder(f).Encode(spec)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

// NewWorkspace creates a new workspace in current directory
func NewWorkspace(state *State, spec *Spec) error {
	// Make sub-directories
	_, err := util.MkDirs("", DirRoot, DirInterm, DirServer, DirClient, DirCSR)
	if err != nil {
		return err
	}

	// Write state file
	err = SaveState(state, FileState)
	if err != nil {
		return err
	}

	// Write spec file
	err = SaveSpec(spec, FileSpec)
	if err != nil {
		return err
	}

	return nil
}

// LoadWorkspace loads an existing workspace
func LoadWorkspace() (*State, *Spec, error) {
	// Load state file
	state, err := LoadState(FileState)
	if err != nil {
		return nil, nil, err
	}

	// Load spec file
	spec, err := LoadSpec(FileSpec)
	if err != nil {
		return nil, nil, err
	}

	return state, spec, nil
}

// SaveWorkspace saves changes to an existing workspace
func SaveWorkspace(state *State, spec *Spec) error {
	// Write state file
	err := SaveState(state, FileState)
	if err != nil {
		return err
	}

	// Write spec file
	err = SaveSpec(spec, FileSpec)
	if err != nil {
		return err
	}

	return nil
}

// CleanupWorkspace removes all directories and files in a workspace
func CleanupWorkspace() error {
	return util.DeleteAll(
		"",
		DirRoot,
		DirInterm,
		DirServer,
		DirClient,
		DirCSR,
		FileState,
		FileSpec,
	)
}
