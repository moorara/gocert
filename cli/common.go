package cli

import (
	"os"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
)

const (
	utilRootEnterConfig   = "\nConfigurations for root certificate authorities ..."
	utilIntermEnterConfig = "\nConfigurations for intermediate certificate authorities ..."
	utilServerEnterConfig = "\nConfigurations for server certificates ..."
	utilClientEnterConfig = "\nConfigurations for client certificates ..."

	utilCommonEnterClaim = "\nCommon specifications for all types of certificates ..."
	utilRootEnterClaim   = "\nSpecifications for root certificate authorities ..."
	utilIntermEnterClaim = "\nSpecifications for intermediate certificate authorities ..."
	utilServerEnterClaim = "\nSpecifications for server certificates ..."
	utilClientEnterClaim = "\nSpecifications for client certificates ..."
)

func copyClaim(c pki.Claim) pki.Claim {
	return pki.Claim{
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

func newColoredUI() *cli.ColoredUi {
	return &cli.ColoredUi{
		OutputColor: cli.UiColorNone,
		InfoColor:   cli.UiColorGreen,
		ErrorColor:  cli.UiColorRed,
		WarnColor:   cli.UiColorYellow,
		Ui: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}
}

// LoadWorkspace reads and parses state and spec from their corresponding files
func LoadWorkspace(ui cli.Ui) (*pki.State, *pki.Spec, int) {
	state, err := pki.LoadState(pki.FileState)
	if err != nil {
		ui.Error("Failed to read state from " + pki.FileState)
		return nil, nil, ErrorReadState
	}

	spec, err := pki.LoadSpec(pki.FileSpec)
	if err != nil {
		ui.Error("Failed to read spec from " + pki.FileSpec)
		return nil, nil, ErrorReadSpec
	}

	return state, spec, 0
}

// SaveWorkspace writes state and spec to their corresponding files
func SaveWorkspace(state *pki.State, spec *pki.Spec, ui cli.Ui) int {
	err := pki.SaveState(state, pki.FileState)
	if err != nil {
		ui.Error("Failed to write state to " + pki.FileState)
		return ErrorWriteState
	}

	err = pki.SaveSpec(spec, pki.FileSpec)
	if err != nil {
		ui.Error("Failed to read spec to " + pki.FileSpec)
		return ErrorWriteSpec
	}

	return 0
}

// AskForNewState creates a new state by asking for data
func AskForNewState(ui cli.Ui) *pki.State {
	root := pki.ConfigCA{}
	ui.Output(utilRootEnterConfig)
	askForData(&root, "yaml", true, ui)

	interm := pki.ConfigCA{}
	ui.Output(utilIntermEnterConfig)
	askForData(&interm, "yaml", true, ui)

	server := pki.Config{}
	ui.Output(utilServerEnterConfig)
	askForData(&server, "yaml", true, ui)

	client := pki.Config{}
	ui.Output(utilClientEnterConfig)
	askForData(&client, "yaml", true, ui)

	return &pki.State{
		Root:   root,
		Interm: interm,
		Server: server,
		Client: client,
	}
}

// AskForNewSpec creates a new spec by asking for data
func AskForNewSpec(ui cli.Ui) *pki.Spec {
	// Common specs
	common := pki.Claim{}
	ui.Output(utilCommonEnterClaim)
	askForData(&common, "toml", true, ui)

	root := copyClaim(common)
	ui.Output(utilRootEnterClaim)
	askForData(&root, "toml", true, ui)

	interm := copyClaim(common)
	ui.Output(utilIntermEnterClaim)
	askForData(&interm, "toml", true, ui)

	server := copyClaim(common)
	ui.Output(utilServerEnterClaim)
	askForData(&server, "toml", true, ui)

	client := copyClaim(common)
	ui.Output(utilClientEnterClaim)
	askForData(&client, "toml", true, ui)

	return &pki.Spec{
		Root:   root,
		Interm: interm,
		Server: server,
		Client: client,
	}
}

// AskForConfig fills in an existing Config by asking for data
func AskForConfig(config *pki.Config, ui cli.Ui) {
	askForData(config, "yaml", false, ui)
}

// AskForConfigCA fills in an existing ConfigCA by asking for data
func AskForConfigCA(configCA *pki.ConfigCA, ui cli.Ui) {
	askForData(configCA, "yaml", false, ui)
}

// AskForClaim fills in an existing Claim by asking for data
func AskForClaim(claim *pki.Claim, ui cli.Ui) {
	askForData(claim, "toml", false, ui)
}
