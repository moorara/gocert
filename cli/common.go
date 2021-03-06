package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
	"github.com/moorara/gocert/util"
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

	textEnterConfig = "\nENTER CONFIGURATIONS FOR %s ..."
	textEnterClaim  = "\nENTER SPECIFICATIONS FOR %s ..."

	textEnterStateTips = ``
	textEnterSpecTips  = `
	You can enter a list by comma-separating values.
	If you don't want to use any of the specs, leave it empty.
	If you don't want to be asked about a spec every time, enter "-" to skip it.
	You can later change these specs by editing "spec.toml" file.`
	textEnterPolicyTips = `
	You can specify the signing policy for certificate authorities.
	Enter the name of each spec you want be matched/supplied as appeared in specs.`
	textEnterConfigTips = `
	Using passwords for certificate authorities is mandatory.
	The password length should be at least 6 characters.`
	textEnterClaimTips = `
	You can enter a list by comma-separating values.
	If you don't want to use any of the specs, leave it empty.`
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
	ui.Info(textEnterStateTips)

	root := pki.Config{}
	ui.Output(textRootEnterConfig)
	err := util.AskForStruct(&root, "yaml", true, nil, ui)
	if err != nil {
		return nil, err
	}

	interm := pki.Config{}
	ui.Output(textIntermEnterConfig)
	err = util.AskForStruct(&interm, "yaml", true, nil, ui)
	if err != nil {
		return nil, err
	}

	server := pki.Config{}
	ui.Output(textServerEnterConfig)
	err = util.AskForStruct(&server, "yaml", true, nil, ui)
	if err != nil {
		return nil, err
	}

	client := pki.Config{}
	ui.Output(textClientEnterConfig)
	err = util.AskForStruct(&client, "yaml", true, nil, ui)
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
	ui.Info(textEnterSpecTips)

	// Common specs
	common := pki.Claim{}
	commonSkip := []string{}
	ui.Output(textCommonEnterClaim)
	err := util.AskForStruct(&common, "toml", true, &commonSkip, ui)
	if err != nil {
		return nil, err
	}

	root := common.Clone()
	rootSkip := make([]string, len(commonSkip))
	copy(rootSkip, commonSkip)
	ui.Output(textRootEnterClaim)
	err = util.AskForStruct(&root, "toml", true, &rootSkip, ui)
	if err != nil {
		return nil, err
	}

	interm := common.Clone()
	intermSkip := make([]string, len(commonSkip))
	copy(intermSkip, commonSkip)
	ui.Output(textIntermEnterClaim)
	err = util.AskForStruct(&interm, "toml", true, &intermSkip, ui)
	if err != nil {
		return nil, err
	}

	server := common.Clone()
	serverSkip := make([]string, len(commonSkip))
	copy(serverSkip, commonSkip)
	ui.Output(textServerEnterClaim)
	err = util.AskForStruct(&server, "toml", true, &serverSkip, ui)
	if err != nil {
		return nil, err
	}

	client := common.Clone()
	clientSkip := make([]string, len(commonSkip))
	copy(clientSkip, commonSkip)
	ui.Output(textClientEnterClaim)
	err = util.AskForStruct(&client, "toml", true, &clientSkip, ui)
	if err != nil {
		return nil, err
	}

	ui.Info(textEnterPolicyTips)

	rootPolicy := pki.Policy{}
	ui.Output(textRootEnterPolicy)
	err = util.AskForStruct(&rootPolicy, "toml", true, nil, ui)
	if err != nil {
		return nil, err
	}

	intermPolicy := pki.Policy{}
	ui.Output(textIntermEnterPolicy)
	err = util.AskForStruct(&intermPolicy, "toml", true, nil, ui)
	if err != nil {
		return nil, err
	}

	metadata := pki.Metadata{}
	if len(rootSkip) > 0 {
		metadata[mdRootSkip] = rootSkip
	}
	if len(intermSkip) > 0 {
		metadata[mdIntermSkip] = intermSkip
	}
	if len(serverSkip) > 0 {
		metadata[mdServerSkip] = serverSkip
	}
	if len(clientSkip) > 0 {
		metadata[mdClientSkip] = clientSkip
	}

	spec := &pki.Spec{
		Root:         root,
		Interm:       interm,
		Server:       server,
		Client:       client,
		RootPolicy:   rootPolicy,
		IntermPolicy: intermPolicy,
		Metadata:     metadata,
	}

	return spec, nil
}

func askForConfig(config *pki.Config, c pki.Cert, skipList *[]string, ui cli.Ui) error {
	// Only Root and Intermediate CAs are asked for password
	if c.Type == pki.CertTypeRoot || c.Type == pki.CertTypeInterm {
		ui.Info(textEnterConfigTips)
	}

	// User certificates should not have a password
	if c.Type == pki.CertTypeServer || c.Type == pki.CertTypeClient {
		config.Password = "bypass"
		defer func() {
			config.Password = ""
		}()
	}

	title := strings.ToUpper(c.Title())
	text := fmt.Sprintf(textEnterConfig, title)
	ui.Output(text)

	return util.AskForStruct(config, "yaml", false, skipList, ui)
}

func askForClaim(claim *pki.Claim, c pki.Cert, skipList *[]string, ui cli.Ui) error {
	ui.Info(textEnterClaimTips)

	title := strings.ToUpper(c.Title())
	text := fmt.Sprintf(textEnterConfig, title)
	ui.Output(text)

	return util.AskForStruct(claim, "toml", false, skipList, ui)
}
