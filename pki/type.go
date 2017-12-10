package pki

import (
	"strings"
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

	defaultRootPolicyMatch    = ""
	defaultRootPolicySupplied = "CommonName"

	defaultIntermPolicyMatch    = ""
	defaultIntermPolicySupplied = "CommonName"
)

type (
	// State represents the type for state
	State struct {
		Root   Config `yaml:"root"`
		Interm Config `yaml:"intermediate"`
		Server Config `yaml:"server"`
		Client Config `yaml:"client"`
	}

	// Config represents the subtype for configurations
	Config struct {
		Serial   int64  `yaml:"serial"`
		Length   int    `yaml:"length"`
		Days     int    `yaml:"days"`
		Password string `yaml:"-" secret:"required,6"`
	}

	// Spec represents the type for specs
	Spec struct {
		Root         Claim  `toml:"root"`
		Interm       Claim  `toml:"intermediate"`
		Server       Claim  `toml:"server"`
		Client       Claim  `toml:"client"`
		RootPolicy   Policy `toml:"root_policy"`
		IntermPolicy Policy `toml:"intermediate_policy"`
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

	// Policy represents the subtype for a policy
	Policy struct {
		Match    []string `toml:"match"`
		Supplied []string `toml:"supplied"`
	}

	// Metadata represents the type for metadata about a certificate
	Metadata struct {
		CertType int
	}
)

// NewState creates a new state
func NewState() *State {
	return &State{
		Root: Config{
			Serial: defaultRootCASerial,
			Length: defaultRootCALength,
			Days:   defaultRootCADays,
		},
		Interm: Config{
			Serial: defaultIntermCASerial,
			Length: defaultIntermCALength,
			Days:   defaultIntermCADays,
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

// NewSpec creates a new spec
func NewSpec() *Spec {
	return &Spec{
		Root:   Claim{},
		Interm: Claim{},
		Server: Claim{},
		Client: Claim{},
		RootPolicy: Policy{
			Match:    strings.Split(defaultRootPolicyMatch, ","),
			Supplied: strings.Split(defaultRootPolicySupplied, ","),
		},
		IntermPolicy: Policy{
			Match:    strings.Split(defaultIntermPolicyMatch, ","),
			Supplied: strings.Split(defaultIntermPolicySupplied, ","),
		},
	}
}

// Dir returns the directory corresponding to cert type
func (md Metadata) Dir() string {
	switch md.CertType {
	case CertTypeRoot:
		return DirRoot
	case CertTypeInterm:
		return DirInterm
	case CertTypeServer:
		return DirServer
	case CertTypeClient:
		return DirClient
	default:
		return ""
	}
}
