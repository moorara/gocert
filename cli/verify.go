package cli

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
)

const (
	verifySuccess       = " ✓ Verified %s"
	verifyFailure       = " ✗ Failed to verify %s. Error: %s"
	verifyEnterNameCA   = "\nENTER NAME FOR CERTIFICATE AUTHORITY ..."
	verifyEnterNameCert = "\nENTER NAME FOR CERTIFICATE ..."

	verifySynopsis = `Verifies a certificate using its certificate authority.`
	verifyHelp     = `
	You can use this command to verify a certificate using its certificate authority
	This command tries to verify the specified certificate by checking the certificate trust chain.

	Flags:
		-ca      the name of certificate authorithy
		-name    the name of certificate
		-dns     if provided, determines whether the certificate can be used for the given dns name
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
func (c *VerifyCommand) Run(args []string) (exit int) {
	var fCA, fName, fDNS string

	flags := flag.NewFlagSet("verify", flag.ContinueOnError)
	flags.Usage = func() {}
	flags.StringVar(&fCA, "ca", "", "")
	flags.StringVar(&fName, "name", "", "")
	flags.StringVar(&fDNS, "dns", "", "")
	err := flags.Parse(args)
	if err != nil {
		return ErrorInvalidFlag
	}

	if fCA == "" {
		c.ui.Output(verifyEnterNameCA)
		fCA, err = c.ui.Ask(fmt.Sprintf(promptTemplate, "CA Name", "string"))
		if err != nil {
			return ErrorInvalidName
		}
	}

	if fName == "" {
		c.ui.Output(verifyEnterNameCert)
		fName, err = c.ui.Ask(fmt.Sprintf(promptTemplate, "Cert Name", "string list"))
		if err != nil {
			return ErrorInvalidName
		}
	}

	c.ui.Output("")

	cCA := resolveByName(fCA)
	certNames := strings.Split(fName, ",")

	for _, certName := range certNames {
		cCert := resolveByName(certName)

		err = c.pki.VerifyCert(cCA, cCert, fDNS)
		if err != nil {
			c.ui.Error(fmt.Sprintf(verifyFailure, cCert.Name, err.Error()))
			exit = ErrorVerify
		} else {
			c.ui.Info(fmt.Sprintf(verifySuccess, cCert.Name))
		}
	}

	c.ui.Output("")

	return exit
}
