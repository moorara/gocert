package test

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// WriteToStdinPipe writes a string to stdin through a pipe
func WriteToStdinPipe(r, w *os.File, str string) (string, error) {
	var in string

	_, err := w.WriteString(str + "\n")
	if err != nil {
		return "", err
	}

	err = w.Close()
	if err != nil {
		return "", err
	}

	_, err = fmt.Fscan(os.Stdin, &in)
	if err != nil {
		return "", err
	}

	return in, nil
}

// ReadFromStdoutPipe reads a written string to stdout through a pipe
func ReadFromStdoutPipe(r, w *os.File, str string) (string, error) {
	var buf bytes.Buffer

	fmt.Fprint(os.Stdout, str)
	err := w.Close()
	if err != nil {
		return "", err
	}

	_, err = io.Copy(&buf, r)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ReadFromStderrPipe reads a written string to stderr through a pipe
func ReadFromStderrPipe(r, w *os.File, str string) (string, error) {
	var buf bytes.Buffer

	fmt.Fprint(os.Stderr, str)
	err := w.Close()
	if err != nil {
		return "", err
	}

	_, err = io.Copy(&buf, r)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
