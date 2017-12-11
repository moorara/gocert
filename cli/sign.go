package cli

import (
	"flag"
	"fmt"
	"os"

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

func certExist(name string, certType int) bool {
	md := pki.Metadata{
		Name:     name,
		CertType: certType,
	}

	if _, err := os.Stat(md.KeyPath()); os.IsNotExist(err) {
		return false
	}
	return true
}

func (c *SignCommand) resolve(mdCA, mdCSR *pki.Metadata) (configCA, configCSR pki.Config, policyCA pki.Policy, status int) {
	state, spec, status := loadWorkspace(c.ui)
	if status != 0 {
		return
	}

	if mdCA.Name == rootName && certExist(mdCA.Name, pki.CertTypeRoot) {
		mdCA.CertType = pki.CertTypeRoot
	} else if certExist(mdCA.Name, pki.CertTypeInterm) {
		mdCA.CertType = pki.CertTypeInterm
	}

	if mdCA.CertType == 0 {
		c.ui.Error("CA name can only be root or an intermediate ca name.")
		status = ErrorInvalidCA
		return
	}

	switch {
	case certExist(mdCSR.Name, pki.CertTypeInterm):
		mdCSR.CertType = pki.CertTypeInterm
	case certExist(mdCSR.Name, pki.CertTypeServer):
		mdCSR.CertType = pki.CertTypeServer
	case certExist(mdCSR.Name, pki.CertTypeClient):
		mdCSR.CertType = pki.CertTypeClient
	default:
		c.ui.Error("Certificate signing request not exist.")
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
	var mdCA, mdCSR pki.Metadata

	flags := flag.NewFlagSet("sign", flag.ContinueOnError)
	flags.Usage = func() {}
	flags.StringVar(&mdCA.Name, "ca", "", "")
	flags.StringVar(&mdCSR.Name, "name", "", "")
	err := flags.Parse(args)
	if err != nil {
		return ErrorInvalidFlag
	}

	if mdCA.Name == "" {
		c.ui.Output(signEnterNameCA)
		mdCA.Name, err = c.ui.Ask(fmt.Sprintf(askTemplate, "CA Name", "string"))
		if err != nil {
			return ErrorInvalidName
		}
	}

	if mdCSR.Name == "" {
		c.ui.Output(signEnterNameCSR)
		mdCSR.Name, err = c.ui.Ask(fmt.Sprintf(askTemplate, "CSR Name", "string"))
		if err != nil {
			return ErrorInvalidName
		}
	}

	if mdCA.Name == mdCSR.Name {
		c.ui.Error("CA name and request name cannot be the same.")
		return ErrorInvalidName
	}

	configCA, configCSR, policyCA, status := c.resolve(&mdCA, &mdCSR)
	if status != 0 {
		return status
	}

	c.ui.Output(signEnterConfigCA)
	askForConfig(&configCA, c.ui)
	c.ui.Output("")

	err = c.pki.SignCSR(configCA, mdCA, configCSR, mdCSR, pki.PolicyTrustFunc(policyCA))
	if err != nil {
		c.ui.Error("Failed to sign certificate request. Error: " + err.Error())
		return ErrorSign
	}

	return 0
}
