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

	for _, test := range tests {
		restore := ReplaceOSArgs(test.args)
		assert.Equal(t, test.args, os.Args)
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

	for _, test := range tests {
		inR, inW, restore, err := PipeStdin()
		assert.NoError(t, err)

		in, err := WriteToStdinPipe(inR, inW, test.in)
		assert.NoError(t, err)
		assert.Equal(t, test.in, in)

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

	for _, test := range tests {
		outR, outW, restore, err := PipeStdout()
		assert.NoError(t, err)

		out, err := ReadFromStdoutPipe(outR, outW, test.out)
		assert.NoError(t, err)
		assert.Equal(t, test.out, out)

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

	for _, test := range tests {
		errR, errW, restore, err := PipeStderr()
		assert.NoError(t, err)

		er, err := ReadFromStderrPipe(errR, errW, test.err)
		assert.NoError(t, err)
		assert.Equal(t, test.err, er)

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

	for _, test := range tests {
		r, w, restore, err := PipeStdoutAndStderr()
		assert.NoError(t, err)

		fmt.Fprint(os.Stdout, test.out)
		fmt.Fprint(os.Stderr, test.err)

		err = w.Close()
		assert.NoError(t, err)
		data, err := ioutil.ReadAll(r)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, string(data))

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

	for _, test := range tests {
		inR, inW, outR, outW, errR, errW, restore, err := PipeStdAll()
		assert.NoError(t, err)

		in, err := WriteToStdinPipe(inR, inW, test.in)
		assert.NoError(t, err)
		assert.Equal(t, test.in, in)

		out, err := ReadFromStdoutPipe(outR, outW, test.out)
		assert.NoError(t, err)
		assert.Equal(t, test.out, out)

		er, err := ReadFromStderrPipe(errR, errW, test.err)
		assert.NoError(t, err)
		assert.Equal(t, test.err, er)

		restore()
	}
}
