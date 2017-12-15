package cli

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
)

const (
	signEnterNameCA    = "\nENTER NAME FOR CERTIFICATE AUTHORITY ..."
	signEnterNameCSR   = "\nENTER NAME FOR CERTIFICATE SIGNING REQUEST ..."
	signEnterConfigCA  = "\nENTER CONFIGURATIONS FOR CERTIFICATE AUTHORITY ..."
	signEnterConfigCSR = "\nENTER CONFIGURATIONS FOR NEW CERTIFICATE ..."

	signSynopsis = `Signs a certificate signing request.`
	signHelp     = `
	You can use this command to sign a certificate signing request (CSR) and create a new certificate.

	You will be asked for entering the password for certificate authorithy.
	The root certificate authorithy can only sign intermediate certificate authorities.
	Intermediate certificate authorities can then sign other intermediate certificate authorities or server/client certificates.

	Flags:
		-ca      the name of certificate authorithy
		-name    the name of certificate signing request
	`
)

// SignCommand represents the sign command for signing a csr
type SignCommand struct {
	ui  cli.Ui
	pki pki.Manager
}

// NewSignCommand creates a new command
func NewSignCommand() *SignCommand {
	return &SignCommand{
		ui:  newColoredUI(),
		pki: pki.NewX509Manager(),
	}
}

func (c *SignCommand) resolve(nameCA, nameCSR string) (configCA pki.Config, mdCA pki.Metadata, configCSR pki.Config, mdCSR pki.Metadata, policyCA pki.Policy, status int) {
	state, spec, status := loadWorkspace(c.ui)
	if status != 0 {
		return
	}

	mdCA = resolveByName(nameCA)
	mdCSR = resolveByName(nameCSR)

	if mdCA.CertType != pki.CertTypeRoot && mdCA.CertType != pki.CertTypeInterm {
		c.ui.Error("Certificate authority name is not valid.")
		status = ErrorInvalidCA
		return
	}

	if mdCSR.CertType == 0 || mdCSR.CertType == pki.CertTypeRoot {
		c.ui.Error("Certificate name is not valid.")
		status = ErrorInvalidCSR
		return
	}

	// Root CA only signs intermediate CAs, and intermediate CA cannot sign root CA
	if mdCA.CertType == pki.CertTypeRoot && mdCSR.CertType != pki.CertTypeInterm {
		c.ui.Error("Root CA can only sign an intermediate ca.")
		status = ErrorInvalidCSR
		return
	}

	// CertType fields are ensured to be valid
	configCA, _ = state.ConfigFor(mdCA.CertType)
	configCSR, _ = state.ConfigFor(mdCSR.CertType)
	policyCA, _ = spec.PolicyFor(mdCA.CertType)

	return
}

// Synopsis returns the short help text for command
func (c *SignCommand) Synopsis() string {
	return signSynopsis
}

// Help returns the long help text for command
func (c *SignCommand) Help() string {
	return signHelp
}

// Run executes the command
func (c *SignCommand) Run(args []string) int {
	var nameCA, nameCSR string

	flags := flag.NewFlagSet("sign", flag.ContinueOnError)
	flags.Usage = func() {}
	flags.StringVar(&nameCA, "ca", "", "")
	flags.StringVar(&nameCSR, "name", "", "")
	err := flags.Parse(args)
	if err != nil {
		return ErrorInvalidFlag
	}

	if nameCA == "" {
		c.ui.Output(signEnterNameCA)
		nameCA, err = c.ui.Ask(fmt.Sprintf(askTemplate, "CA Name", "string"))
		if err != nil {
			return ErrorInvalidName
		}
	}

	if nameCSR == "" {
		c.ui.Output(signEnterNameCSR)
		nameCSR, err = c.ui.Ask(fmt.Sprintf(askTemplate, "CSR Name", "string"))
		if err != nil {
			return ErrorInvalidName
		}
	}

	if nameCA == nameCSR {
		c.ui.Error("CA name and request name cannot be the same.")
		return ErrorInvalidName
	}

	configCA, mdCA, configCSR, mdCSR, policyCA, status := c.resolve(nameCA, nameCSR)
	if status != 0 {
		return status
	}

	c.ui.Output(signEnterConfigCA)
	askForConfig(&configCA, c.ui)
	c.ui.Output("")

	trustFunc := pki.PolicyTrustFunc(policyCA)
	err = c.pki.SignCSR(configCA, mdCA, configCSR, mdCSR, trustFunc)
	if err != nil {
		c.ui.Error("Failed to sign certificate request. Error: " + err.Error())
		return ErrorSign
	}

	return 0
}
