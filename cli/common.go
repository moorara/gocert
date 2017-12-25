package cli

import (
	"os"

	"github.com/mitchellh/cli"
	"github.com/moorara/go-box/util"
	"github.com/moorara/gocert/help"
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

func askForNewState(ui cli.Ui) (*pki.State, error) {
	root := pki.Config{}
	ui.Output(textRootEnterConfig)
	err := help.AskForStruct(&root, "yaml", true, ui)
	if err != nil {
		return nil, err
	}

	interm := pki.Config{}
	ui.Output(textIntermEnterConfig)
	err = help.AskForStruct(&interm, "yaml", true, ui)
	if err != nil {
		return nil, err
	}

	server := pki.Config{}
	ui.Output(textServerEnterConfig)
	err = help.AskForStruct(&server, "yaml", true, ui)
	if err != nil {
		return nil, err
	}

	client := pki.Config{}
	ui.Output(textClientEnterConfig)
	err = help.AskForStruct(&client, "yaml", true, ui)
	if err != nil {
		return nil, err
	}

	state := &pki.State{
		Root:   root,
		Interm: interm,
		Server: server,
		Client: client,
	}

	return state, nil
}

func askForNewSpec(ui cli.Ui) (*pki.Spec, error) {
	// Common specs
	common := pki.Claim{}
	ui.Output(textCommonEnterClaim)
	err := help.AskForStruct(&common, "toml", true, ui)
	if err != nil && !util.IsStringIn(err.Error(), "EOF", "unexpected newline ") {
		return nil, err
	}

	root := common.Clone()
	ui.Output(textRootEnterClaim)
	err = help.AskForStruct(&root, "toml", true, ui)
	if err != nil {
		return nil, err
	}

	interm := common.Clone()
	ui.Output(textIntermEnterClaim)
	err = help.AskForStruct(&interm, "toml", true, ui)
	if err != nil {
		return nil, err
	}

	server := common.Clone()
	ui.Output(textServerEnterClaim)
	err = help.AskForStruct(&server, "toml", true, ui)
	if err != nil {
		return nil, err
	}

	client := common.Clone()
	ui.Output(textClientEnterClaim)
	err = help.AskForStruct(&client, "toml", true, ui)
	if err != nil {
		return nil, err
	}

	rootPolicy := pki.Policy{}
	ui.Output(textRootEnterPolicy)
	err = help.AskForStruct(&rootPolicy, "toml", true, ui)
	if err != nil {
		return nil, err
	}

	intermPolicy := pki.Policy{}
	ui.Output(textIntermEnterPolicy)
	err = help.AskForStruct(&intermPolicy, "toml", true, ui)
	if err != nil {
		return nil, err
	}

	spec := &pki.Spec{
		Root:         root,
		Interm:       interm,
		Server:       server,
		Client:       client,
		RootPolicy:   rootPolicy,
		IntermPolicy: intermPolicy,
	}

	return spec, nil
}

func askForConfig(config *pki.Config, ui cli.Ui) error {
	return help.AskForStruct(config, "yaml", false, ui)
}

func askForClaim(claim *pki.Claim, ui cli.Ui) error {
	return help.AskForStruct(claim, "toml", false, ui)
}
