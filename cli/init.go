package cli

import (
	"flag"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
	"github.com/moorara/gocert/util"
)

const (
	initSuccess = "\n ✓ Workspace Setup \n"

	initSynopsis = `Initializes a new workspace with desired configs and specs.`
	initHelp     = `
	You can use this command to initialize a new workspace with desired configs and specs.

	You will be first asked for entering common specs which all of your certificates share. So, you enter them once.
	Next, you will be asked for entering more-specific specs for Root, Intermediate, Server, and Client certificates.

	Best-practice configs are provided by default.
	You can customize these configs by editing "state.yaml" file.
	`
)

// InitCommand represents the command for initialization
type InitCommand struct {
	ui cli.Ui
}

// NewInitCommand creates a new command
func NewInitCommand() *InitCommand {
	return &InitCommand{
		ui: newColoredUI(),
	}
}

// Synopsis returns the short help text for command
func (c *InitCommand) Synopsis() string {
	return initSynopsis
}

// Help returns the long help text for command
func (c *InitCommand) Help() string {
	return initHelp
}

// Run executes the command
func (c *InitCommand) Run(args []string) int {
	flags := flag.NewFlagSet("init", flag.ContinueOnError)
	flags.Usage = func() {}
	err := flags.Parse(args)
	if err != nil {
		return ErrorInvalidFlag
	}

	_, err = util.MkDirs("", pki.DirRoot, pki.DirInterm, pki.DirServer, pki.DirClient, pki.DirCSR)
	if err != nil {
		c.ui.Error("Failed to create directories. Error: " + err.Error())
		return ErrorMakeDir
	}

	state := pki.NewState()
	err = pki.SaveState(state, pki.FileState)
	if err != nil {
		c.ui.Error("Failed to save configs. Error: " + err.Error())
		return ErrorWriteState
	}

	// Ask user to enter values for spec
	spec, err := askForNewSpec(c.ui)
	if err != nil {
		return ErrorEnterSpec
	}

	err = pki.SaveSpec(spec, pki.FileSpec)
	if err != nil {
		c.ui.Error("Failed to save specs. Error: " + err.Error())
		return ErrorWriteSpec
	}

	c.ui.Info(initSuccess)

	return 0
}
