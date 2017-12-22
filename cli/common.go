package cli

import (
	"os"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
)

const (
	textRootEnterConfig   = "\nCONFIGURATIONS FOR ROOT CERTIFICATE AUTHORITIES ..."
	textIntermEnterConfig = "\nCONFIGURATIONS FOR INTERMEDIATE CERTIFICATE AUTHORITIES ..."
	textServerEnterConfig = "\nCONFIGURATIONS FOR SERVER CERTIFICATES ..."
	textClientEnterConfig = "\nCONFIGURATIONS FOR CLIENT CERTIFICATES ..."

	textCommonEnterClaim = "\nCOMMON SPECIFICATIONS FOR ALL TYPES OF CERTIFICATES ..."
	textRootEnterClaim   = "\nSPECIFICATIONS FOR ROOT CERTIFICATE AUTHORITIES ..."
	textIntermEnterClaim = "\nSPECIFICATIONS FOR INTERMEDIATE CERTIFICATE AUTHORITIES ..."
	textServerEnterClaim = "\nSPECIFICATIONS FOR SERVER CERTIFICATES ..."
	textClientEnterClaim = "\nSPECIFICATIONS FOR CLIENT CERTIFICATES ..."

	textRootEnterPolicy   = "\nTRUST POLICY RULES FOR ROOT CERTIFICATE AUTHORITIES ..."
	textIntermEnterPolicy = "\nTRUST POLICY RULES FOR INTERMEDIATE CERTIFICATE AUTHORITIES ..."
)

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

func loadWorkspace(ui cli.Ui) (*pki.State, *pki.Spec, int) {
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

func saveWorkspace(state *pki.State, spec *pki.Spec, ui cli.Ui) int {
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

func resolveByName(name string) pki.Cert {
	var c pki.Cert

	c.Name, c.Type = name, pki.CertTypeRoot
	if _, err := os.Stat(c.KeyPath()); err == nil && name == rootName {
		return c
	}

	c.Name, c.Type = name, pki.CertTypeInterm
	if _, err := os.Stat(c.KeyPath()); err == nil {
		return c
	}

	c.Name, c.Type = name, pki.CertTypeServer
	if _, err := os.Stat(c.KeyPath()); err == nil {
		return c
	}

	c.Name, c.Type = name, pki.CertTypeClient
	if _, err := os.Stat(c.KeyPath()); err == nil {
		return c
	}

	return pki.Cert{}
}

func askForNewState(ui cli.Ui) *pki.State {
	root := pki.Config{}
	ui.Output(textRootEnterConfig)
	askForStruct(&root, "yaml", true, ui)

	interm := pki.Config{}
	ui.Output(textIntermEnterConfig)
	askForStruct(&interm, "yaml", true, ui)

	server := pki.Config{}
	ui.Output(textServerEnterConfig)
	askForStruct(&server, "yaml", true, ui)

	client := pki.Config{}
	ui.Output(textClientEnterConfig)
	askForStruct(&client, "yaml", true, ui)

	return &pki.State{
		Root:   root,
		Interm: interm,
		Server: server,
		Client: client,
	}
}

func askForNewSpec(ui cli.Ui) *pki.Spec {
	// Common specs
	common := pki.Claim{}
	ui.Output(textCommonEnterClaim)
	askForStruct(&common, "toml", true, ui)

	root := common.Clone()
	ui.Output(textRootEnterClaim)
	askForStruct(&root, "toml", true, ui)

	interm := common.Clone()
	ui.Output(textIntermEnterClaim)
	askForStruct(&interm, "toml", true, ui)

	server := common.Clone()
	ui.Output(textServerEnterClaim)
	askForStruct(&server, "toml", true, ui)

	client := common.Clone()
	ui.Output(textClientEnterClaim)
	askForStruct(&client, "toml", true, ui)

	rootPolicy := pki.Policy{}
	ui.Output(textRootEnterPolicy)
	askForStruct(&rootPolicy, "toml", true, ui)

	intermPolicy := pki.Policy{}
	ui.Output(textIntermEnterPolicy)
	askForStruct(&intermPolicy, "toml", true, ui)

	return &pki.Spec{
		Root:         root,
		Interm:       interm,
		Server:       server,
		Client:       client,
		RootPolicy:   rootPolicy,
		IntermPolicy: intermPolicy,
	}
}

func askForConfig(config *pki.Config, ui cli.Ui) {
	askForStruct(config, "yaml", false, ui)
}

func askForClaim(claim *pki.Claim, ui cli.Ui) {
	askForStruct(claim, "toml", false, ui)
}
