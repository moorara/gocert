package pki

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	yaml "gopkg.in/yaml.v2"
)

const (
	defaultRootCASerial = int64(10)
	defaultRootCALength = 4096
	defaultRootCADays   = 20 * 365

	defaultIntermCASerial = int64(100)
	defaultIntermCALength = 4096
	defaultIntermCADays   = 10 * 365

	defaultServerCertSerial = int64(1000)
	defaultServerCertLength = 2048
	defaultServerCertDays   = 10 + 365

	defaultClientCertSerial = int64(10000)
	defaultClientCertLength = 2048
	defaultClientCertDays   = 10 + 30
)

type (
	// State represents the type for state
	State struct {
		Root   ConfigCA `yaml:"root"`
		Interm ConfigCA `yaml:"intermediate"`
		Server Config   `yaml:"server"`
		Client Config   `yaml:"client"`
	}

	// Config represents the subtype for configurations
	Config struct {
		Serial int64 `yaml:"serial"`
		Length int   `yaml:"length"`
		Days   int   `yaml:"days"`
	}

	// ConfigCA represents the subtype for certificatea authority configurations
	ConfigCA struct {
		Config   `yaml:",inline"`
		Password string `yaml:"-" secret:"required,6"`
	}

	// Spec represents the type for specs
	Spec struct {
		Root   Claim `toml:"root"`
		Interm Claim `toml:"intermediate"`
		Server Claim `toml:"server"`
		Client Claim `toml:"client"`
	}

	// Claim represents the subtype for an identity claim
	Claim struct {
		CommonName         string   `toml:"-"`
		Country            []string `toml:"country"`
		Province           []string `toml:"province"`
		Locality           []string `toml:"locality"`
		Organization       []string `toml:"organization"`
		OrganizationalUnit []string `toml:"organizational_unit"`
		StreetAddress      []string `toml:"street_address"`
		PostalCode         []string `toml:"postal_code"`
		EmailAddress       []string `toml:"email_address"`
	}
)

// NewState creates a new state
func NewState() *State {
	return &State{
		Root: ConfigCA{
			Config: Config{
				Serial: defaultRootCASerial,
				Length: defaultRootCALength,
				Days:   defaultRootCADays,
			},
		},
		Interm: ConfigCA{
			Config: Config{
				Serial: defaultIntermCASerial,
				Length: defaultIntermCALength,
				Days:   defaultIntermCADays,
			},
		},
		Server: Config{
			Serial: defaultServerCertSerial,
			Length: defaultServerCertLength,
			Days:   defaultServerCertDays,
		},
		Client: Config{
			Serial: defaultClientCertSerial,
			Length: defaultClientCertLength,
			Days:   defaultClientCertDays,
		},
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

// NewSpec creates a new spec
func NewSpec() *Spec {
	return &Spec{
		Root:   Claim{},
		Interm: Claim{},
		Server: Claim{},
		Client: Claim{},
	}
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
