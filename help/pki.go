package help

import (
	"github.com/moorara/gocert/pki"
)

// MockManager mocks pki.Manager
type MockManager struct {
	GenCertError    error
	GenCSRError     error
	SignCSRError    error
	VerifyCertError error

	GenCertCalled    bool
	GenCSRCalled     bool
	SignCSRCalled    bool
	VerifyCertCalled bool
}

// GenCert is mock
func (m *MockManager) GenCert(pki.Config, pki.Claim, pki.Cert) error {
	m.GenCertCalled = true
	return m.GenCertError
}

// GenCSR is mock
func (m *MockManager) GenCSR(pki.Config, pki.Claim, pki.Cert) error {
	m.GenCSRCalled = true
	return m.GenCSRError
}

// SignCSR is mock
func (m *MockManager) SignCSR(pki.Config, pki.Cert, pki.Config, pki.Cert, pki.TrustFunc) error {
	m.SignCSRCalled = true
	return m.SignCSRError
}

// VerifyCert is mock
func (m *MockManager) VerifyCert(pki.Cert, pki.Cert, string) error {
	m.VerifyCertCalled = true
	return m.VerifyCertError
}
