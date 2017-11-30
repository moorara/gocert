package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/moorara/go-box/test"
	"github.com/stretchr/testify/assert"
)

func TestRunApp(t *testing.T) {
	tests := []struct {
		args            []string
		expectedExit    int
		expectedRegexes []string
	}{
		{
			[]string{},
			127,
			[]string{},
		},
		{
			[]string{"invalid"},
			127,
			[]string{},
		},
		{
			[]string{"-version"},
			0,
			[]string{},
		},
		{
			[]string{"--version"},
			0,
			[]string{},
		},
		{
			[]string{"-help"},
			0,
			[]string{},
		},
		{
			[]string{"--help"},
			0,
			[]string{},
		},
		{
			[]string{"cert"},
			0,
			[]string{},
		},
		{
			[]string{"cert", "-root"},
			0,
			[]string{},
		},
	}

	// Null out standard io streams temporarily
	_, _, _, _, _, _, restore, err := test.PipeStdAll()
	defer restore()
	assert.NoError(t, err)

	for _, tc := range tests {
		status := runApp(tc.args)

		assert.Equal(t, tc.expectedExit, status)
	}
}

func TestMain(t *testing.T) {
	/* main calls os.Exit(), so we need to deal with it! */

	if os.Getenv("TEST_SUCCESS") == "1" {
		os.Args = []string{"gotls", "--version"}
		main()
	}

	if os.Getenv("TEST_FAIL") == "1" {
		os.Args = []string{"gotls"}
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
