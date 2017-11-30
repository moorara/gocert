package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipeStdin(t *testing.T) {
	tests := []struct {
		in string
	}{
		{"input"},
		{"milad"},
	}

	for _, test := range tests {
		inR, inW, restore, err := PipeStdin()
		defer restore()
		assert.NoError(t, err)

		in, err := WriteToStdinPipe(inR, inW, test.in)
		assert.NoError(t, err)
		assert.Equal(t, test.in, in)
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
		defer restore()
		assert.NoError(t, err)

		out, err := ReadFromStdoutPipe(outR, outW, test.out)
		assert.NoError(t, err)
		assert.Equal(t, test.out, out)
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
		defer restore()
		assert.NoError(t, err)

		er, err := ReadFromStderrPipe(errR, errW, test.err)
		assert.NoError(t, err)
		assert.Equal(t, test.err, er)
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
		defer restore()
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
	}
}
