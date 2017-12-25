package help

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/mitchellh/cli"
)

// MockUI implements the same interface as github.com/mitchellh/cli.MockUi
type MockUI struct {
	cli.MockUi
	reader *bufio.Reader
}

// NewMockUI creates a new MockUI
func NewMockUI(r io.Reader) *MockUI {
	ui := cli.NewMockUi()
	ui.InputReader = r
	reader := bufio.NewReader(r)
	return &MockUI{*ui, reader}
}

// Ask lets you read empty strings and strings with spaces
func (u *MockUI) Ask(query string) (string, error) {
	fmt.Fprint(u.MockUi.OutputWriter, query)
	line, err := u.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	result := strings.Trim(line, "\t\r\n")

	return result, nil
}

// AskSecret lets you read empty strings and strings with spaces
func (u *MockUI) AskSecret(query string) (string, error) {
	return u.Ask(query)
}
