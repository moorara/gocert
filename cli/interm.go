package cli

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
)

const (
	intermEnterName   = "\nNAME FOR INTERMEDIATE CERTIFICATE AUTHORITY ..."
	intermEnterConfig = "\nCONFIGURATIONS FOR INTERMEDIATE CERTIFICATE AUTHORITY ..."
	intermEnterClaim  = "\nSPECIFICATIONS FOR INTERMEDIATE CERTIFICATE AUTHORITY ..."

	intermNewSynopsis = ""
	intermNewHelp     = ``
)

// IntermNewCommand represents the "intermediate new" command
type IntermNewCommand struct {
	ui  cli.Ui
	pki pki.Manager
}

// NewIntermNewCommand creates an IntermNewCommand
func NewIntermNewCommand() *IntermNewCommand {
	return &IntermNewCommand{
		ui:  newColoredUI(),
		pki: pki.NewX509Manager(),
	}
}

// Synopsis returns the short help text for intermediate command
func (c *IntermNewCommand) Synopsis() string {
	return intermNewSynopsis
}

// Help returns the long help text for intermediate command
func (c *IntermNewCommand) Help() string {
	return intermNewHelp
}

// Run executes the intermediate command
func (c *IntermNewCommand) Run(args []string) int {
	var fReq string

	flags := flag.NewFlagSet("intermediate new", flag.ContinueOnError)
	flags.Usage = func() {}
	flags.StringVar(&fReq, "req", "", "")
	err := flags.Parse(args)
	if err != nil {
		return ErrorInvalidFlag
	}

	state, spec, status := loadWorkspace(c.ui)
	if status != 0 {
		return status
	}

	if fReq == "" {
		c.ui.Output(intermEnterName)
		fReq, err = c.ui.Ask(fmt.Sprintf(askTemplate, "Name", "string"))
		if err != nil {
			return ErrorNoName
		}
	}

	c.ui.Output(intermEnterConfig)
	askForConfigCA(&state.Interm, c.ui)

	c.ui.Output(intermEnterClaim)
	askForClaim(&spec.Interm, c.ui)

	c.ui.Output("")

	// TODO: deal with fReq, req name, and existing files!
	err = c.pki.GenIntermCSR(fReq, state.Interm, spec.Interm)
	if err != nil {
		c.ui.Error("Failed to generate intermediate ca. Error: " + err.Error())
		return ErrorIntermCA
	}

	return 0
}
