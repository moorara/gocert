package cli

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
)

const (
	reqEnterName   = "\nENTER NAME FOR ..."
	reqEnterConfig = "\nENTER CONFIGURATIONS ..."
	reqEnterClaim  = "\nENTER SPECIFICATIONS ..."

	reqSynopsis = "Creates a new certificate signing request (CSR)."
	reqHelp     = `
  You can use this command to create a new certificate signing request (CSR).
  The generated request can be later signed by a certificate authority to create the actual certificate.

  You will be asked for entering those specifications not set in spec.toml file.
  These specifications are supposed to be certificate-specific and not common across all ceritificates.
  You can enter a list by comma-separating values. If you don't want to use any of the entries, leave it empty.

  Flags:
    -name    set a name for the new certificate
  `
)

// ReqCommand represents the a req command for generating a new csr
type ReqCommand struct {
	ui  cli.Ui
	pki pki.Manager
	md  pki.Metadata
}

// NewReqCommand creates a ReqCommand
func NewReqCommand(md pki.Metadata) *ReqCommand {
	return &ReqCommand{
		ui:  newColoredUI(),
		pki: pki.NewX509Manager(),
		md:  md,
	}
}

func (c *ReqCommand) load() (config pki.Config, claim pki.Claim, status int) {
	state, spec, status := loadWorkspace(c.ui)
	if status != 0 {
		return
	}

	switch c.md.CertType {
	case pki.CertTypeInterm:
		config = state.Interm
		claim = spec.Interm
	case pki.CertTypeServer:
		config = state.Server
		claim = spec.Server
	case pki.CertTypeClient:
		config = state.Client
		claim = spec.Client
	default:
		status = ErrorMetadata
		return
	}

	// User certificates should not have a password
	if c.md.CertType == pki.CertTypeServer || c.md.CertType == pki.CertTypeClient {
		config.Password = "bypass"
	}

	return
}

// Synopsis returns the short help text for command
func (c *ReqCommand) Synopsis() string {
	return reqSynopsis
}

// Help returns the long help text for command
func (c *ReqCommand) Help() string {
	return reqHelp
}

// Run executes the command
func (c *ReqCommand) Run(args []string) int {
	var fName string

	flags := flag.NewFlagSet("req", flag.ContinueOnError)
	flags.Usage = func() {}
	flags.StringVar(&fName, "name", "", "")
	err := flags.Parse(args)
	if err != nil {
		return ErrorInvalidFlag
	}

	config, claim, status := c.load()
	if status != 0 {
		return status
	}

	if fName == "" {
		c.ui.Output(reqEnterName)
		fName, err = c.ui.Ask(fmt.Sprintf(askTemplate, "Name", "string"))
		if err != nil {
			return ErrorNoName
		}
	}

	c.ui.Output(reqEnterConfig)
	askForConfig(&config, c.ui)
	c.ui.Output(reqEnterClaim)
	askForClaim(&claim, c.ui)
	c.ui.Output("")

	err = c.pki.GenCSR(fName, config, claim, c.md)
	if err != nil {
		c.ui.Error("Failed to generate certificate signing request. Error: " + err.Error())
		return ErrorCSR
	}

	return 0
}
