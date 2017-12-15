package cli

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
)

const (
	verifyEnterNameCA   = "\nENTER NAME FOR CERTIFICATE AUTHORITY ..."
	verifyEnterNameCert = "\nENTER NAME FOR CERTIFICATE ..."

	verifySynopsis = `Verifies a certificate using its certificate authority.`
	verifyHelp     = `
	You can use this command to verify a certificate using its certificate authority
	This command tries to verify the specified certificate by checking the certificate trust chain.

	Flags:
		-ca      the name of certificate authorithy
		-name    the name of certificate
	`
)

// VerifyCommand represents the verify command
type VerifyCommand struct {
	ui  cli.Ui
	pki pki.Manager
}

// NewVerifyCommand creates a new command
func NewVerifyCommand() *VerifyCommand {
	return &VerifyCommand{
		ui:  newColoredUI(),
		pki: pki.NewX509Manager(),
	}
}

// Synopsis returns the short help text for command
func (c *VerifyCommand) Synopsis() string {
	return verifySynopsis
}

// Help returns the long help text for command
func (c *VerifyCommand) Help() string {
	return verifyHelp
}

// Run executes the command
func (c *VerifyCommand) Run(args []string) int {
	var nameCA, nameCert string

	flags := flag.NewFlagSet("verify", flag.ContinueOnError)
	flags.Usage = func() {}
	flags.StringVar(&nameCA, "ca", "", "")
	flags.StringVar(&nameCert, "name", "", "")
	err := flags.Parse(args)
	if err != nil {
		return ErrorInvalidFlag
	}

	if nameCA == "" {
		c.ui.Output(verifyEnterNameCA)
		nameCA, err = c.ui.Ask(fmt.Sprintf(askTemplate, "CA Name", "string"))
		if err != nil {
			return ErrorInvalidName
		}
	}

	if nameCert == "" {
		c.ui.Output(verifyEnterNameCert)
		nameCert, err = c.ui.Ask(fmt.Sprintf(askTemplate, "Cert Name", "string"))
		if err != nil {
			return ErrorInvalidName
		}
	}

	c.ui.Output("")
	mdCA := resolveByName(nameCA)
	mdCert := resolveByName(nameCert)

	err = c.pki.VerifyCert(mdCA, mdCert)
	if err != nil {
		c.ui.Error("Failed to verify certificate. Error: " + err.Error())
		return ErrorVerify
	}

	c.ui.Info(" ✓ OK")
	c.ui.Output("")

	return 0
}
