package config

import (
	"io/ioutil"

	"github.com/mitchellh/cli"
	yaml "gopkg.in/yaml.v2"
)

const (
	textSettingsEnterRoot   = "\nSettings for root certificate authorities ..."
	textSettingsEnterInterm = "\nSettings for intermediate certificate authorities ..."
	textSettingsEnterServer = "\nSettings for server certificates ..."
	textSettingsEnterClient = "\nSettings for client certificates ..."
)

type (
	// State represents the type for state
	State struct {
		Root   SettingsCA `yaml:"root"`
		Interm SettingsCA `yaml:"intermediate"`
		Server Settings   `yaml:"server"`
		Client Settings   `yaml:"client"`
	}

	// Settings represents the subtype for settings
	Settings struct {
		Serial int64 `yaml:"serial"`
		Length int   `yaml:"length"`
		Days   int   `yaml:"days"`
	}

	// SettingsCA represents the subtype for settings
	SettingsCA struct {
		Settings `yaml:",inline"`
		Password string `yaml:"-" secret:"true"`
	}
)

// NewState creates a new state
func NewState() *State {
	return &State{
		Root: SettingsCA{
			Settings: Settings{
				Serial: defaultRootCASerial,
				Length: defaultRootCALength,
				Days:   defaultRootCADays,
			},
		},
		Interm: SettingsCA{
			Settings: Settings{
				Serial: defaultIntermCASerial,
				Length: defaultIntermCALength,
				Days:   defaultIntermCADays,
			},
		},
		Server: Settings{
			Serial: defaultServerCertSerial,
			Length: defaultServerCertLength,
			Days:   defaultServerCertDays,
		},
		Client: Settings{
			Serial: defaultClientCertSerial,
			Length: defaultClientCertLength,
			Days:   defaultClientCertDays,
		},
	}
}

// NewStateWithInput creates a new state with user inputs
func NewStateWithInput(ui cli.Ui) *State {
	root := SettingsCA{}
	ui.Output(textSettingsEnterRoot)
	fillIn(&root, "yaml", true, ui)

	interm := SettingsCA{}
	ui.Output(textSettingsEnterInterm)
	fillIn(&interm, "yaml", true, ui)

	server := Settings{}
	ui.Output(textSettingsEnterServer)
	fillIn(&server, "yaml", true, ui)

	client := Settings{}
	ui.Output(textSettingsEnterClient)
	fillIn(&client, "yaml", true, ui)

	return &State{
		Root:   root,
		Interm: interm,
		Server: server,
		Client: client,
	}
}

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

// FillIn asks for input for empty fields
func (s *Settings) FillIn(ui cli.Ui) {
	fillIn(s, "yaml", false, ui)
}

// FillIn asks for input for empty fields
func (s *SettingsCA) FillIn(ui cli.Ui) {
	fillIn(s, "yaml", false, ui)
}
