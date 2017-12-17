package pki

import (
	"net"
	"path"
	"strings"
)

const (
	extKey     = ".key"
	extCert    = ".cert"
	extCSR     = ".csr"
	extCAKey   = ".ca.key"
	extCACert  = ".ca.cert"
	extCACSR   = ".ca.csr"
	extCAChain = ".ca.chain"

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

	titleRoot   = "Root Certificate Authority"
	titleInterm = "Intermediate Certificate Authority"
	titleServer = "Server Certificate Authority"
	titleClient = "Client Certificate Authority"
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
		DNSName            []string `toml:"dns_name"`
		IPAddress          []net.IP `toml:"ip_address"`
		EmailAddress       []string `toml:"email_address"`
		StreetAddress      []string `toml:"street_address"`
		PostalCode         []string `toml:"postal_code"`
	}

	// Policy represents the subtype for a policy
	Policy struct {
		Match    []string `toml:"match"`
		Supplied []string `toml:"supplied"`
	}

	// Metadata represents the type for metadata about a certificate
	Metadata struct {
		Name     string
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

// ConfigFor returns config for a certificate type
func (s *State) ConfigFor(certType int) (Config, bool) {
	switch certType {
	case CertTypeRoot:
		return s.Root, true
	case CertTypeInterm:
		return s.Interm, true
	case CertTypeServer:
		return s.Server, true
	case CertTypeClient:
		return s.Client, true
	default:
		return Config{}, false
	}
}

// ClaimFor returns claim for a certificate type
func (s *Spec) ClaimFor(certType int) (Claim, bool) {
	switch certType {
	case CertTypeRoot:
		return s.Root, true
	case CertTypeInterm:
		return s.Interm, true
	case CertTypeServer:
		return s.Server, true
	case CertTypeClient:
		return s.Client, true
	default:
		return Claim{}, false
	}
}

// PolicyFor returns policy for a certificate type
func (s *Spec) PolicyFor(certType int) (Policy, bool) {
	switch certType {
	case CertTypeRoot:
		return s.RootPolicy, true
	case CertTypeInterm:
		return s.IntermPolicy, true
	default:
		return Policy{}, false
	}
}

// Clone return a deep copy of claim
func (c Claim) Clone() Claim {
	return Claim{
		CommonName:         c.CommonName,
		Country:            c.Country,
		Province:           c.Province,
		Locality:           c.Locality,
		Organization:       c.Organization,
		OrganizationalUnit: c.OrganizationalUnit,
		DNSName:            c.DNSName,
		IPAddress:          c.IPAddress,
		EmailAddress:       c.EmailAddress,
		StreetAddress:      c.StreetAddress,
		PostalCode:         c.PostalCode,
	}
}

// Title returns a descriptive title
func (md Metadata) Title() string {
	switch md.CertType {
	case CertTypeRoot:
		return titleRoot
	case CertTypeInterm:
		return titleInterm
	case CertTypeServer:
		return titleServer
	case CertTypeClient:
		return titleClient
	default:
		return ""
	}
}

// KeyPath returns key file path
func (md Metadata) KeyPath() string {
	if md.Name == "" {
		return ""
	}

	switch md.CertType {
	case CertTypeRoot:
		return path.Join(DirRoot, md.Name+extCAKey)
	case CertTypeInterm:
		return path.Join(DirInterm, md.Name+extCAKey)
	case CertTypeServer:
		return path.Join(DirServer, md.Name+extKey)
	case CertTypeClient:
		return path.Join(DirClient, md.Name+extKey)
	default:
		return ""
	}
}

// CertPath returns cert file path
func (md Metadata) CertPath() string {
	if md.Name == "" {
		return ""
	}

	switch md.CertType {
	case CertTypeRoot:
		return path.Join(DirRoot, md.Name+extCACert)
	case CertTypeInterm:
		return path.Join(DirInterm, md.Name+extCACert)
	case CertTypeServer:
		return path.Join(DirServer, md.Name+extCert)
	case CertTypeClient:
		return path.Join(DirClient, md.Name+extCert)
	default:
		return ""
	}
}

// CSRPath returns certificate signing request file path
func (md Metadata) CSRPath() string {
	if md.Name == "" {
		return ""
	}

	switch md.CertType {
	case CertTypeInterm:
		return path.Join(DirCSR, md.Name+extCACSR)
	case CertTypeServer:
		return path.Join(DirCSR, md.Name+extCSR)
	case CertTypeClient:
		return path.Join(DirCSR, md.Name+extCSR)
	default:
		return ""
	}
}

// ChainPath returns certificate chain file path
func (md Metadata) ChainPath() string {
	if md.Name == "" {
		return ""
	}

	switch md.CertType {
	case CertTypeRoot:
		return path.Join(DirRoot, md.Name+extCACert)
	case CertTypeInterm:
		return path.Join(DirInterm, md.Name+extCAChain)
	default:
		return ""
	}
}
