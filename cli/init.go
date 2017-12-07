package cli

import (
	"github.com/mitchellh/cli"
	"github.com/moorara/go-box/util"
	"github.com/moorara/gocert/pki"
)

const (
	initSynopsis = "Initializes a new workspace with desired configurations and specifications."
	initHelp     = `
	You can use this command to initialize a new workspace with desired configurations and specifications.

	You will be first asked for entering common specs which all of your certificates share. So, you enter them once.
	Next, you will be asked for entering more-specific specs for Root CA, Intermediate CA, Server, and Client certificates.
	You can enter a list by comma-separating values. If you don't want to use any of the specs, leave it empty.
	You can later change these specs by editing spec.toml file.

	Best-practice configurations are provided by default.
	You can change these configurations by editing state.yaml file.
	`
)

// InitCommand represents the init command
type InitCommand struct {
	ui cli.Ui
}

// NewInitCommand creates an init command
func NewInitCommand() *InitCommand {
	return &InitCommand{
		ui: newColoredUI(),
	}
}

// Synopsis returns the short help text for init command
func (c *InitCommand) Synopsis() string {
	return initSynopsis
}

// Help returns the long help text for init command
func (c *InitCommand) Help() string {
	return initHelp
}

// Run executes the init command
func (c *InitCommand) Run(args []string) int {
	// Make sub-directories
	_, err := util.MkDirs("", pki.DirRoot, pki.DirInterm, pki.DirServer, pki.DirClient, pki.DirCSR)
	if err != nil {
		return ErrorWriteState
	}

	// Write default state file
	state := pki.NewState()
	err = pki.SaveState(state, pki.FileState)
	if err != nil {
		return ErrorWriteState
	}

	// Write default spec file
	spec := AskForNewSpec(c.ui)
	err = pki.SaveSpec(spec, pki.FileSpec)
	if err != nil {
		return ErrorWriteSpec
	}

	c.ui.Output("")
	return 0
}
