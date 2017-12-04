package new

import (
	"github.com/mitchellh/cli"
	"github.com/moorara/go-box/util"
	"github.com/moorara/gocert/common"
	"github.com/moorara/gocert/config"
)

const (
	// ErrorMakeDir is returned when cannot make a directory
	ErrorMakeDir = 11
	// ErrorWriteState is returned when cannot write state file
	ErrorWriteState = 12
	// ErrorWriteSpec is returned when cannot write spec file
	ErrorWriteSpec = 13

	textSynopsis = "Initializes a new workspace with desired specs and settings."
	textHelp     = `
	You can use this command to initialize a new workspace with desired specs and settings.

	You will be first asked for entering common specs which all of your certificates share. So, you enter them once.
	Next, you will be asked for entering more-specific specs for Root CA, Intermediate CA, Server, and Client certificates.
	You can enter a list by comma-separating values. If you don't want to use any of the specs, leave it empty.
	You can later change these specs by editing spec.toml file.

	Best-practice settings are provided by default.
	You can change these settings by editing state.yaml file.
	`
)

// Command represents the new command
type Command struct {
	ui cli.Ui
}

// NewCommand creates a new command
func NewCommand() *Command {
	return &Command{
		ui: common.NewColoredUI(),
	}
}

// Synopsis returns the short help text for command
func (c *Command) Synopsis() string {
	return textSynopsis
}

// Help returns the long help text for command
func (c *Command) Help() string {
	return textHelp
}

// Run executes the command
func (c *Command) Run(args []string) int {
	// Make sub-directories
	_, err := util.MkDirs("", config.DirNameRoot, config.DirNameInterm, config.DirNameServer, config.DirNameClient)
	if err != nil {
		return ErrorWriteState
	}

	// Write default state file
	state := config.NewState()
	err = config.SaveState(state, config.FileNameState)
	if err != nil {
		return ErrorWriteState
	}

	// Write default spec file
	spec := config.NewSpecWithInput(c.ui)
	err = config.SaveSpec(spec, config.FileNameSpec)
	if err != nil {
		return ErrorWriteSpec
	}

	c.ui.Output("")
	return 0
}
