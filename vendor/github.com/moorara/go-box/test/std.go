package test

import (
	"os"
)

// PipeStdin pipes standard input stream temporarily
func PipeStdin() (inR, inW *os.File, restore func(), err error) {
	origStdin := os.Stdin
	restore = func() {
		os.Stdin = origStdin
	}

	// Pipe stdin
	inR, inW, err = os.Pipe()
	if err != nil {
		return
	}
	os.Stdin = inR

	return
}

// PipeStdout pipes standard output stream temporarily
func PipeStdout() (outR, outW *os.File, restore func(), err error) {
	origStdout := os.Stdout
	restore = func() {
		os.Stdout = origStdout
	}

	// Pipe stdout
	outR, outW, err = os.Pipe()
	if err != nil {
		return
	}
	os.Stdout = outW

	return
}

// PipeStderr pipes standard error stream temporarily
func PipeStderr() (errR, errW *os.File, restore func(), err error) {
	origStderr := os.Stderr
	restore = func() {
		os.Stderr = origStderr
	}

	// Pipe stderr
	errR, errW, err = os.Pipe()
	if err != nil {
		return
	}
	os.Stderr = errW

	return
}

// PipeStdAll pipes all standard streams temporarily
func PipeStdAll() (inR, inW, outR, outW, errR, errW *os.File, restore func(), err error) {
	inR, inW, inRst, err := PipeStdin()
	if err != nil {
		return
	}

	outR, outW, outRst, err := PipeStdout()
	if err != nil {
		return
	}

	errR, errW, errRst, err := PipeStderr()
	if err != nil {
		return
	}

	restore = func() {
		inRst()
		outRst()
		errRst()
	}

	return
}
