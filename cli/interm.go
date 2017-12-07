package cli

import (
	"flag"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
)

const (
	intermediateEnterConfig = "\nConfigurations for intermediate certificate authority ..."
	intermediateEnterClaim  = "\nSpecifications for intermediate certificate authority ..."

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

	flags := flag.NewFlagSet("interm", flag.ContinueOnError)
	flags.Usage = func() {}
	flags.StringVar(&fReq, "req", "", "")
	err := flags.Parse(args)
	if err != nil {
		return ErrorInvalidFlag
	}

	state, spec, status := LoadWorkspace(c.ui)
	if status != 0 {
		return status
	}

	c.ui.Output(intermediateEnterConfig)
	AskForConfigCA(&state.Interm, c.ui)

	c.ui.Output(intermediateEnterClaim)
	AskForClaim(&spec.Interm, c.ui)

	err = c.pki.GenIntermCA(state.Interm, spec.Interm)
	if err != nil {
		c.ui.Error("Failed to generate intermediate ca. Error: " + err.Error())
		return ErrorIntermCA
	}

	c.ui.Output("")
	return 0
}
