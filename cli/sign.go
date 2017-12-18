package cli

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
)

const (
	signSuccess        = " ✓ Signed %s"
	signFailure        = " ✗ Failed to sign %s. Error: %s"
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

func (c *SignCommand) resolveCA(state *pki.State, spec *pki.Spec, nameCA string) (configCA pki.Config, mdCA pki.Metadata, policyCA pki.Policy, status int) {
	mdCA = resolveByName(nameCA)

	if mdCA.CertType != pki.CertTypeRoot && mdCA.CertType != pki.CertTypeInterm {
		c.ui.Error("Certificate authority name is not valid.")
		status = ErrorInvalidCA
		return
	}

	// CertType field is ensured to be valid
	configCA, _ = state.ConfigFor(mdCA.CertType)
	policyCA, _ = spec.PolicyFor(mdCA.CertType)

	return
}

func (c *SignCommand) resolveCSR(state *pki.State, spec *pki.Spec, nameCSR string) (configCSR pki.Config, mdCSR pki.Metadata, status int) {
	mdCSR = resolveByName(nameCSR)

	if mdCSR.CertType == 0 || mdCSR.CertType == pki.CertTypeRoot {
		c.ui.Error("Certificate name is not valid.")
		status = ErrorInvalidCSR
		return
	}

	// CertType field is ensured to be valid
	configCSR, _ = state.ConfigFor(mdCSR.CertType)

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
func (c *SignCommand) Run(args []string) (exit int) {
	var fCA, fName string

	flags := flag.NewFlagSet("sign", flag.ContinueOnError)
	flags.Usage = func() {}
	flags.StringVar(&fCA, "ca", "", "")
	flags.StringVar(&fName, "name", "", "")
	err := flags.Parse(args)
	if err != nil {
		return ErrorInvalidFlag
	}

	if fCA == "" {
		c.ui.Output(signEnterNameCA)
		fCA, err = c.ui.Ask(fmt.Sprintf(askTemplate, "CA Name", "string"))
		if err != nil {
			return ErrorInvalidCA
		}
	}

	if fName == "" {
		c.ui.Output(signEnterNameCSR)
		fName, err = c.ui.Ask(fmt.Sprintf(askTemplate, "CSR Name", "string list"))
		if err != nil {
			return ErrorInvalidCSR
		}
	}

	if fCA == fName {
		c.ui.Error("CA name and request name cannot be the same.")
		return ErrorInvalidName
	}

	state, spec, status := loadWorkspace(c.ui)
	if status != 0 {
		return status
	}

	configCA, mdCA, policyCA, status := c.resolveCA(state, spec, fCA)
	if status != 0 {
		return status
	}

	c.ui.Output(signEnterConfigCA)
	askForConfig(&configCA, c.ui)
	c.ui.Output("")

	trustFunc := pki.PolicyTrustFunc(policyCA)
	csrNames := strings.Split(fName, ",")

	for _, csrName := range csrNames {
		configCSR, mdCSR, status := c.resolveCSR(state, spec, csrName)
		if status != 0 {
			return status
		}

		// Root CA only signs intermediate CAs, and intermediate CA cannot sign root CA
		if mdCA.CertType == pki.CertTypeRoot && mdCSR.CertType != pki.CertTypeInterm {
			c.ui.Error("Root CA can only sign an intermediate ca.")
			return ErrorInvalidCSR
		}

		err = c.pki.SignCSR(configCA, mdCA, configCSR, mdCSR, trustFunc)
		if err != nil {
			c.ui.Error(fmt.Sprintf(signFailure, mdCSR.Name, err.Error()))
			exit = ErrorSign
		} else {
			c.ui.Info(fmt.Sprintf(signSuccess, mdCSR.Name))
		}
	}

	c.ui.Output("")

	return exit
}
