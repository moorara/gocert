package cli

import (
	"io/ioutil"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/moorara/go-box/util"
	"github.com/stretchr/testify/assert"
)

var (
	helpRegexes = []string{
		`Usage: \w+ \[--version\] \[--help\] <command> \[<args>\]`,
		`Available commands are:`,
	}

	helpSubRegexes = []string{
		`This command is accessed by using one of the subcommands below.`,
		`Subcommands:`,
	}

	helpMockInit      = "help text for mocked init command"
	helpMockRootNew   = "help text for mocked root new command"
	helpMockIntermNew = "help text for mocked intermediate new command"
)

func newMockApp(name, version string) *App {
	return &App{
		name:      name,
		version:   version,
		init:      &cli.MockCommand{RunResult: 0, HelpText: helpMockInit},
		rootNew:   &cli.MockCommand{RunResult: 0, HelpText: helpMockRootNew},
		intermNew: &cli.MockCommand{RunResult: 0, HelpText: helpMockIntermNew},
	}
}

func TestNewApp(t *testing.T) {
	tests := []struct {
		name    string
		version string
	}{
		{"gotls", "0.1.0-1"},
		{"gocert", "0.2.0-100"},
	}

	for _, test := range tests {
		app := NewApp(test.name, test.version)

		assert.NotNil(t, app.name)
		assert.NotNil(t, app.version)
		assert.NotNil(t, app.init)
		assert.NotNil(t, app.rootNew)
		assert.NotNil(t, app.intermNew)
	}
}

func TestAppRun(t *testing.T) {
	tests := []struct {
		name, version   string
		args            []string
		expectedExit    int
		expectedRegexes []string
	}{
		{"cli", "0.1.1", []string{}, 127, helpRegexes},
		{"cli", "0.1.2", []string{"invalid"}, 127, helpRegexes},

		{"cli", "0.1.3", []string{"-version"}, 0, []string{"0.1.3"}},
		{"cli", "0.1.4", []string{"--version"}, 0, []string{"0.1.4"}},

		{"cli", "0.1.5", []string{"-help"}, 0, helpRegexes},
		{"cli", "0.1.6", []string{"--help"}, 0, helpRegexes},

		{"cli", "0.1.7", []string{"init"}, 0, []string{}},
		{"cli", "0.1.8", []string{"init", "-help"}, 0, []string{helpMockInit}},
		{"cli", "0.1.9", []string{"init", "--help"}, 0, []string{helpMockInit}},

		{"cli", "0.1.10", []string{"root"}, 1, helpSubRegexes},
		{"cli", "0.1.11", []string{"root", "-help"}, 0, helpSubRegexes},
		{"cli", "0.1.12", []string{"root", "--help"}, 0, helpSubRegexes},
		{"cli", "0.1.11", []string{"root", "new", "-help"}, 0, []string{helpMockRootNew}},
		{"cli", "0.1.12", []string{"root", "new", "--help"}, 0, []string{helpMockRootNew}},

		{"cli", "0.1.13", []string{"intermediate"}, 1, helpSubRegexes},
		{"cli", "0.1.14", []string{"intermediate", "-help"}, 0, helpSubRegexes},
		{"cli", "0.1.15", []string{"intermediate", "--help"}, 0, helpSubRegexes},
		{"cli", "0.1.14", []string{"intermediate", "new", "-help"}, 0, []string{helpMockIntermNew}},
		{"cli", "0.1.15", []string{"intermediate", "new", "--help"}, 0, []string{helpMockIntermNew}},
	}

	for _, test := range tests {
		r, w, restore, err := util.PipeStdoutAndStderr()
		assert.NoError(t, err)

		app := newMockApp(test.name, test.version)
		status := app.Run(test.args)

		err = w.Close()
		assert.NoError(t, err)
		data, err := ioutil.ReadAll(r)
		assert.NoError(t, err)
		output := string(data)

		assert.Equal(t, test.expectedExit, status)
		for _, rx := range test.expectedRegexes {
			assert.Regexp(t, rx, output)
		}

		restore()
	}
}
