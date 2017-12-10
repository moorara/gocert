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
	helpMockServerNew = "help text for mocked server new command"
	helpMockClientNew = "help text for mocked client new command"
)

func newMockApp(name, version string) *App {
	return &App{
		name:      name,
		version:   version,
		init:      &cli.MockCommand{RunResult: 0, HelpText: helpMockInit},
		rootNew:   &cli.MockCommand{RunResult: 0, HelpText: helpMockRootNew},
		intermNew: &cli.MockCommand{RunResult: 0, HelpText: helpMockIntermNew},
		serverNew: &cli.MockCommand{RunResult: 0, HelpText: helpMockServerNew},
		clientNew: &cli.MockCommand{RunResult: 0, HelpText: helpMockClientNew},
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

		{"cli", "0.2.1", []string{"-version"}, 0, []string{"0.2.1"}},
		{"cli", "0.2.2", []string{"--version"}, 0, []string{"0.2.2"}},

		{"cli", "0.3.1", []string{"-help"}, 0, helpRegexes},
		{"cli", "0.3.2", []string{"--help"}, 0, helpRegexes},

		{"cli", "0.4.1", []string{"init"}, 0, []string{}},
		{"cli", "0.4.2", []string{"init", "-help"}, 0, []string{helpMockInit}},
		{"cli", "0.4.3", []string{"init", "--help"}, 0, []string{helpMockInit}},

		{"cli", "0.5.1", []string{"root"}, 1, helpSubRegexes},
		{"cli", "0.5.2", []string{"root", "-help"}, 0, helpSubRegexes},
		{"cli", "0.5.3", []string{"root", "--help"}, 0, helpSubRegexes},
		{"cli", "0.5.4", []string{"root", "new", "-help"}, 0, []string{helpMockRootNew}},
		{"cli", "0.5.5", []string{"root", "new", "--help"}, 0, []string{helpMockRootNew}},

		{"cli", "0.6.1", []string{"intermediate"}, 1, helpSubRegexes},
		{"cli", "0.6.2", []string{"intermediate", "-help"}, 0, helpSubRegexes},
		{"cli", "0.6.3", []string{"intermediate", "--help"}, 0, helpSubRegexes},
		{"cli", "0.6.4", []string{"intermediate", "new", "-help"}, 0, []string{helpMockIntermNew}},
		{"cli", "0.6.5", []string{"intermediate", "new", "--help"}, 0, []string{helpMockIntermNew}},

		{"cli", "0.7.1", []string{"server"}, 1, helpSubRegexes},
		{"cli", "0.7.2", []string{"server", "-help"}, 0, helpSubRegexes},
		{"cli", "0.7.3", []string{"server", "--help"}, 0, helpSubRegexes},
		{"cli", "0.7.4", []string{"server", "new", "-help"}, 0, []string{helpMockServerNew}},
		{"cli", "0.7.5", []string{"server", "new", "--help"}, 0, []string{helpMockServerNew}},

		{"cli", "0.8.1", []string{"client"}, 1, helpSubRegexes},
		{"cli", "0.8.2", []string{"client", "-help"}, 0, helpSubRegexes},
		{"cli", "0.8.3", []string{"client", "--help"}, 0, helpSubRegexes},
		{"cli", "0.8.4", []string{"client", "new", "-help"}, 0, []string{helpMockClientNew}},
		{"cli", "0.8.5", []string{"client", "new", "--help"}, 0, []string{helpMockClientNew}},
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
