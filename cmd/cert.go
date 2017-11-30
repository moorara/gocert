package cmd

import (
	"fmt"

	"github.com/mitchellh/cli"
)

// Cert ...
type Cert struct {
	ui cli.Ui
}

// NewCert ...
func NewCert() *Cert {
	return &Cert{
		ui: newColoredUI(),
	}
}

// Help ...
func (c *Cert) Help() string {
	return "cert command help"
}

// Run ...
func (c *Cert) Run(args []string) int {
	fmt.Printf("args: %+v\n", args)
	c.ui.Error("Error!")
	c.ui.Info("Info!")
	c.ui.Output("Output")
	c.ui.Warn("Warn!")
	return 0
}

// Synopsis ...
func (c *Cert) Synopsis() string {
	return "cert command short help"
}
