package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceOSArgs(t *testing.T) {
	tests := []struct {
		args []string
	}{
		{[]string{}},
		{[]string{""}},
		{[]string{"go"}},
		{[]string{"go", "test"}},
	}

	origArgs := make([]string, len(os.Args))
	copy(origArgs, os.Args)

	for _, tc := range tests {
		restore := ReplaceOSArgs(tc.args)
		assert.Equal(t, tc.args, os.Args)
		restore()
	}

	assert.Equal(t, origArgs, os.Args)
}

func TestPipeStdin(t *testing.T) {
	tests := []struct {
		in string
	}{
		{"input"},
		{"milad"},
	}

	for _, tc := range tests {
		inR, inW, restore, err := PipeStdin()
		assert.NoError(t, err)

		in, err := WriteToStdinPipe(inR, inW, tc.in)
		assert.NoError(t, err)
		assert.Equal(t, tc.in, in)

		restore()
	}
}

func TestPipeStdout(t *testing.T) {
	tests := []struct {
		out string
	}{
		{"output"},
		{"moorara"},
	}

	for _, tc := range tests {
		outR, outW, restore, err := PipeStdout()
		assert.NoError(t, err)

		out, err := ReadFromStdoutPipe(outR, outW, tc.out)
		assert.NoError(t, err)
		assert.Equal(t, tc.out, out)

		restore()
	}
}

func TestPipeStderr(t *testing.T) {
	tests := []struct {
		err string
	}{
		{"error"},
		{"bad descriptor"},
	}

	for _, tc := range tests {
		errR, errW, restore, err := PipeStderr()
		assert.NoError(t, err)

		er, err := ReadFromStderrPipe(errR, errW, tc.err)
		assert.NoError(t, err)
		assert.Equal(t, tc.err, er)

		restore()
	}
}

func TestPipeStdoutAndStderr(t *testing.T) {
	tests := []struct {
		out      string
		err      string
		expected string
	}{
		{"out", "err", "outerr"},
		{"output", "error", "outputerror"},
	}

	for _, tc := range tests {
		r, w, restore, err := PipeStdoutAndStderr()
		assert.NoError(t, err)

		fmt.Fprint(os.Stdout, tc.out)
		fmt.Fprint(os.Stderr, tc.err)

		err = w.Close()
		assert.NoError(t, err)
		data, err := ioutil.ReadAll(r)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, string(data))

		restore()
	}
}

func TestPipeStdAll(t *testing.T) {
	tests := []struct {
		in, out, err string
	}{
		{"input", "output", "error"},
		{"milad", "moorara", "bad descriptor"},
	}

	for _, tc := range tests {
		inR, inW, outR, outW, errR, errW, restore, err := PipeStdAll()
		assert.NoError(t, err)

		in, err := WriteToStdinPipe(inR, inW, tc.in)
		assert.NoError(t, err)
		assert.Equal(t, tc.in, in)

		out, err := ReadFromStdoutPipe(outR, outW, tc.out)
		assert.NoError(t, err)
		assert.Equal(t, tc.out, out)

		er, err := ReadFromStderrPipe(errR, errW, tc.err)
		assert.NoError(t, err)
		assert.Equal(t, tc.err, er)

		restore()
	}
}
