package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
)

type mockUI struct {
	*cli.MockUi
	reader *bufio.Reader
}

func newMockUI(r io.Reader) *mockUI {
	ui := cli.NewMockUi()
	ui.InputReader = r
	reader := bufio.NewReader(r)
	return &mockUI{ui, reader}
}

func (u *mockUI) Ask(query string) (string, error) {
	_, _ = fmt.Fprint(u.MockUi.OutputWriter, query)

	line, err := u.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	result := strings.Trim(line, "\t\r\n")

	return result, nil
}

func (u *mockUI) AskSecret(query string) (string, error) {
	return u.Ask(query)
}

type mockManager struct {
	GenCertError    error
	GenCSRError     error
	SignCSRError    error
	VerifyCertError error

	GenCertCalled    bool
	GenCSRCalled     bool
	SignCSRCalled    bool
	VerifyCertCalled bool
}

func (m *mockManager) GenCert(pki.Config, pki.Claim, pki.Cert) error {
	m.GenCertCalled = true
	return m.GenCertError
}

func (m *mockManager) GenCSR(pki.Config, pki.Claim, pki.Cert) error {
	m.GenCSRCalled = true
	return m.GenCSRError
}

func (m *mockManager) SignCSR(pki.Config, pki.Cert, pki.Config, pki.Cert, pki.TrustFunc) error {
	m.SignCSRCalled = true
	return m.SignCSRError
}

func (m *mockManager) VerifyCert(pki.Cert, pki.Cert, string) error {
	m.VerifyCertCalled = true
	return m.VerifyCertError
}
