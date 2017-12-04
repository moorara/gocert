package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/moorara/go-box/util"
	"github.com/moorara/gocert/version"
	"github.com/stretchr/testify/assert"
)

var (
	helpRegexes = []string{
		`Usage: gocert \[--version\] \[--help\] <command> \[<args>\]`,
		`Available commands are:`,
	}
	versionRegexes = []string{
		`\d+\.\d+\.\d+-\d+\+[0-9a-f]{7}`,
		`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z\+\d{4}`,
		`go\d\.\d(\.\d+)?`,
	}

	helpMockedNew = "help text for mocked new"
)

func mockVersion() {
	version.Version = "0.1.0-27"
	version.Revision = "abcdeff"
	version.Branch = "test"
	version.BuildTime = "2017-12-03T20:00:17Z+0000"
}

func mockCommands() {
	cmdNew = &cli.MockCommand{RunResult: 0, HelpText: helpMockedNew}
}

func TestRunApp(t *testing.T) {
	tests := []struct {
		args            []string
		expectedExit    int
		expectedRegexes []string
	}{
		{[]string{}, 127, helpRegexes},
		{[]string{"invalid"}, 127, helpRegexes},

		{[]string{"-version"}, 0, versionRegexes},
		{[]string{"--version"}, 0, versionRegexes},

		{[]string{"-help"}, 0, helpRegexes},
		{[]string{"--help"}, 0, helpRegexes},

		{[]string{"new"}, 0, []string{}},
		{[]string{"new", "-help"}, 0, []string{helpMockedNew}},
		{[]string{"new", "--help"}, 0, []string{helpMockedNew}},
	}

	for _, test := range tests {
		r, w, restore, err := util.PipeStdoutAndStderr()
		assert.NoError(t, err)

		mockVersion()
		mockCommands()
		status := runApp(test.args)

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

func TestMain(t *testing.T) {
	/* main calls os.Exit(), so we need to deal with it! */

	if os.Getenv("TEST_SUCCESS") == "1" {
		os.Args = []string{"gocert", "-version"}
		main()
	}

	if os.Getenv("TEST_FAIL") == "1" {
		os.Args = []string{"gocert"}
		main()
	}

	name := os.Args[0]
	args := []string{"-test.run=TestMain"}

	cmd := exec.Command(name, args...)
	cmd.Env = append(os.Environ(), "TEST_SUCCESS=1")
	err := cmd.Run()
	assert.NoError(t, err)

	cmd = exec.Command(name, args...)
	cmd.Env = append(os.Environ(), "TEST_FAIL=1")
	err = cmd.Run()
	e, ok := err.(*exec.ExitError)
	assert.True(t, ok)
	assert.False(t, e.Success())
}
