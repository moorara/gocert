package cmd

import (
	"os"

	"github.com/mitchellh/cli"
)

func newColoredUI() *cli.ColoredUi {
	return &cli.ColoredUi{
		OutputColor: cli.UiColorCyan,
		InfoColor:   cli.UiColorGreen,
		ErrorColor:  cli.UiColorRed,
		WarnColor:   cli.UiColorYellow,
		Ui: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}
}
