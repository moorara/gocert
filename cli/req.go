package cli

import (
	"bytes"
	"flag"
	"fmt"
	"strings"
	"text/template"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
)

const (
	reqMessageRoot  = "\n ✓ Created %s\n"
	reqMessageOther = "\n ✓ Requested %s\n"
	reqEnterName    = "\nENTER NAME FOR %s ..."
	reqEnterConfig  = "\nENTER CONFIGURATIONS FOR %s ..."
	reqEnterClaim   = "\nENTER SPECIFICATIONS FOR %s ..."

	reqSynopsis = `Creates a new {{if eq .CertType 1}}root certificate authority{{- else}}certificate signing request{{- end}}.`
	reqHelp     = `
	{{if eq .CertType 1}}
	You can use this command to create a new root certificate authority (CA).
	The generated CA can be used for signing more intermediate certificate authorities.
	{{- else}}
	You can use this command to create a new certificate signing request (CSR).
	The generated request can be later signed by a certificate authority to create the actual certificate.
	{{- end}}

	{{if eq .CertType 1}}
	The name of root certificate authority will be "root" by default.
	{{- else}}
	You need to choose a name for the new certificate and its signing request.
	{{- end}}
	You will be asked for entering those specifications not set in "spec.toml" file.
	These specifications are supposed to be certificate-specific and not common across all ceritificates.
	You can enter a list by comma-separating values. If you don't want to use any of the entries, leave it empty.
	{{if ne .CertType 1}}

	Flags:
		-name    set a name for the new certificate
	{{- end}}
	`
)

// ReqCommand represents the command for generating a new csr
type ReqCommand struct {
	ui  cli.Ui
	pki pki.Manager
	md  pki.Metadata
}

// NewReqCommand creates a new command
func NewReqCommand(md pki.Metadata) *ReqCommand {
	return &ReqCommand{
		ui:  newColoredUI(),
		pki: pki.NewX509Manager(),
		md:  md,
	}
}

func (c *ReqCommand) output(text string) {
	text = fmt.Sprintf(text, strings.ToUpper(c.md.Title()))
	c.ui.Output(text)
}

// Synopsis returns the short help text for command
func (c *ReqCommand) Synopsis() string {
	var buf bytes.Buffer
	t := template.Must(template.New("synopsis").Parse(reqSynopsis))
	t.Execute(&buf, c.md) // In case of error, empty string will be returned
	return buf.String()
}

// Help returns the long help text for command
func (c *ReqCommand) Help() string {
	var buf bytes.Buffer
	t := template.Must(template.New("help").Parse(reqHelp))
	t.Execute(&buf, c.md) // In case of error, empty string will be returned
	return buf.String()
}

// Run executes the command
func (c *ReqCommand) Run(args []string) int {
	flags := flag.NewFlagSet("req", flag.ContinueOnError)
	flags.Usage = func() {}
	flags.StringVar(&c.md.Name, "name", "", "")
	err := flags.Parse(args)
	if err != nil {
		return ErrorInvalidFlag
	}

	// There should be only one root ca with a default name
	if c.md.CertType == pki.CertTypeRoot {
		c.md.Name = rootName
	}

	if c.md.Name == "" {
		c.output(reqEnterName)
		c.md.Name, err = c.ui.Ask(fmt.Sprintf(askTemplate, "Name", "string"))
		if err != nil {
			return ErrorInvalidName
		}
	}

	state, spec, status := loadWorkspace(c.ui)
	if status != 0 {
		return status
	}

	config, ok1 := state.ConfigFor(c.md.CertType)
	claim, ok2 := spec.ClaimFor(c.md.CertType)
	if !ok1 || !ok2 {
		return ErrorInvalidMetadata
	}

	// User certificates should not have a password
	if c.md.CertType == pki.CertTypeServer || c.md.CertType == pki.CertTypeClient {
		config.Password = "bypass"
	}

	c.output(reqEnterConfig)
	askForConfig(&config, c.ui)
	config.Password = ""
	c.output(reqEnterClaim)
	askForClaim(&claim, c.ui)

	if c.md.CertType == pki.CertTypeRoot {
		err = c.pki.GenCert(config, claim, c.md)
		if err != nil {
			c.ui.Error("Failed to generate root ca. Error: " + err.Error())
			return ErrorCert
		}
		c.ui.Info(fmt.Sprintf(reqMessageRoot, c.md.Name))
	} else {
		err = c.pki.GenCSR(config, claim, c.md)
		if err != nil {
			c.ui.Error("Failed to generate certificate signing request. Error: " + err.Error())
			return ErrorCSR
		}
		c.ui.Info(fmt.Sprintf(reqMessageOther, c.md.Name))
	}

	return 0
}
