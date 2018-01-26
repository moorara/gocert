package util

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// ReplaceOSArgs replaces current args with custom args and returns a restore function
func ReplaceOSArgs(args []string) func() {
	origArgs := make([]string, len(os.Args))
	copy(origArgs, os.Args)
	os.Args = args

	return func() {
		os.Args = origArgs
	}
}

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

// WriteToStdinPipe writes a string to stdin through a pipe
func WriteToStdinPipe(inR, inW *os.File, str string) (string, error) {
	var in string

	_, err := inW.WriteString(str + "\n")
	if err != nil {
		return "", err
	}

	err = inW.Close()
	if err != nil {
		return "", err
	}

	_, err = fmt.Fscan(os.Stdin, &in)
	if err != nil {
		return "", err
	}

	return in, nil
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

// ReadFromStdoutPipe reads a written string to stdout through a pipe
func ReadFromStdoutPipe(outR, outW *os.File, str string) (string, error) {
	var buf bytes.Buffer

	fmt.Fprint(os.Stdout, str)
	err := outW.Close()
	if err != nil {
		return "", err
	}

	_, err = io.Copy(&buf, outR)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
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

// ReadFromStderrPipe reads a written string to stderr through a pipe
func ReadFromStderrPipe(errR, errW *os.File, str string) (string, error) {
	var buf bytes.Buffer

	fmt.Fprint(os.Stderr, str)
	err := errW.Close()
	if err != nil {
		return "", err
	}

	_, err = io.Copy(&buf, errR)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// PipeStdoutAndStderr pipes standard out and error streams together temporarily
func PipeStdoutAndStderr() (r, w *os.File, restore func(), err error) {
	origStdout := os.Stdout
	origStderr := os.Stderr
	restore = func() {
		os.Stdout = origStdout
		os.Stderr = origStderr
	}

	// Pipe stdout
	r, w, err = os.Pipe()
	if err != nil {
		return
	}
	os.Stdout = w
	os.Stderr = w

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
