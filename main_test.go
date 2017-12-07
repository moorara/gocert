package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
