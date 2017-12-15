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

	helpMockInit   = "help text for mocked init command"
	helpMockRoot   = "help text for mocked root command"
	helpMockInterm = "help text for mocked intermediate command"
	helpMockServer = "help text for mocked server command"
	helpMockClient = "help text for mocked client command"
	helpMockSign   = "help text for mocked sign command"
	helpMockVerify = "help text for mocked verify command"
)

func newMockApp(name, version string) *App {
	return &App{
		name:    name,
		version: version,
		init:    &cli.MockCommand{RunResult: 0, HelpText: helpMockInit},
		root:    &cli.MockCommand{RunResult: 0, HelpText: helpMockRoot},
		interm:  &cli.MockCommand{RunResult: 0, HelpText: helpMockInterm},
		server:  &cli.MockCommand{RunResult: 0, HelpText: helpMockServer},
		client:  &cli.MockCommand{RunResult: 0, HelpText: helpMockClient},
		sign:    &cli.MockCommand{RunResult: 0, HelpText: helpMockSign},
		verify:  &cli.MockCommand{RunResult: 0, HelpText: helpMockVerify},
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
		assert.NotNil(t, app.root)
		assert.NotNil(t, app.interm)
		assert.NotNil(t, app.server)
		assert.NotNil(t, app.client)
		assert.NotNil(t, app.sign)
		assert.NotNil(t, app.verify)
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

		{"cli", "0.5.1", []string{"root"}, 0, nil},
		{"cli", "0.5.2", []string{"root", "-help"}, 0, []string{helpMockRoot}},
		{"cli", "0.5.3", []string{"root", "--help"}, 0, []string{helpMockRoot}},

		{"cli", "0.6.1", []string{"intermediate"}, 0, nil},
		{"cli", "0.6.2", []string{"intermediate", "-help"}, 0, []string{helpMockInterm}},
		{"cli", "0.6.3", []string{"intermediate", "--help"}, 0, []string{helpMockInterm}},

		{"cli", "0.7.1", []string{"server"}, 0, nil},
		{"cli", "0.7.2", []string{"server", "-help"}, 0, []string{helpMockServer}},
		{"cli", "0.7.3", []string{"server", "--help"}, 0, []string{helpMockServer}},

		{"cli", "0.8.1", []string{"client"}, 0, nil},
		{"cli", "0.8.2", []string{"client", "-help"}, 0, []string{helpMockClient}},
		{"cli", "0.8.3", []string{"client", "--help"}, 0, []string{helpMockClient}},

		{"cli", "0.9.1", []string{"sign"}, 0, nil},
		{"cli", "0.9.2", []string{"sign", "-help"}, 0, []string{helpMockSign}},
		{"cli", "0.9.3", []string{"sign", "--help"}, 0, []string{helpMockSign}},

		{"cli", "0.10.1", []string{"verify"}, 0, nil},
		{"cli", "0.10.2", []string{"verify", "-help"}, 0, []string{helpMockVerify}},
		{"cli", "0.10.3", []string{"verify", "--help"}, 0, []string{helpMockVerify}},
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
