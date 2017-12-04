package config

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/cli"
)

const (
	textClaimEnterCommon = "\nCommon specs for all types of certificates ..."
	textClaimEnterRoot   = "\nSpecs for root certificate authorities ..."
	textClaimEnterInterm = "\nSpecs for intermediate certificate authorities ..."
	textClaimEnterServer = "\nSpecs for server certificates ..."
	textClaimEnterClient = "\nSpecs for client certificates ..."
)

type (
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

func copyClaim(c Claim) Claim {
	return Claim{
		CommonName:         c.CommonName,
		Country:            c.Country,
		Province:           c.Province,
		Locality:           c.Locality,
		Organization:       c.Organization,
		OrganizationalUnit: c.OrganizationalUnit,
		StreetAddress:      c.StreetAddress,
		PostalCode:         c.PostalCode,
		EmailAddress:       c.EmailAddress,
	}
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

// NewSpecWithInput creates a new spec with user inputs
func NewSpecWithInput(ui cli.Ui) *Spec {
	// Common specs
	common := Claim{}
	ui.Output(textClaimEnterCommon)
	fillIn(&common, "toml", false, ui)

	root := copyClaim(common)
	ui.Output(textClaimEnterRoot)
	fillIn(&root, "toml", false, ui)

	interm := copyClaim(common)
	ui.Output(textClaimEnterInterm)
	fillIn(&interm, "toml", false, ui)

	server := copyClaim(common)
	ui.Output(textClaimEnterServer)
	fillIn(&server, "toml", false, ui)

	client := copyClaim(common)
	ui.Output(textClaimEnterClient)
	fillIn(&client, "toml", false, ui)

	return &Spec{
		Root:   root,
		Interm: interm,
		Server: server,
		Client: client,
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
