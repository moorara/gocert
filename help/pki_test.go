package help

import (
	"errors"
	"testing"

	"github.com/moorara/gocert/pki"
	"github.com/stretchr/testify/assert"
)

func TestMockManager(t *testing.T) {
	tests := []struct {
		GenCertError    error
		GenCSRError     error
		SignCSRError    error
		VerifyCertError error
	}{
		{
			nil,
			nil,
			nil,
			nil,
		},
		{
			errors.New("GenCertError"),
			errors.New("GenCSRError"),
			errors.New("SignCSRError"),
			errors.New("VerifyCertError"),
		},
	}

	for _, test := range tests {
		manager := &MockManager{
			GenCertError:    test.GenCertError,
			GenCSRError:     test.GenCSRError,
			SignCSRError:    test.SignCSRError,
			VerifyCertError: test.VerifyCertError,
		}

		err := manager.GenCert(pki.Config{}, pki.Claim{}, pki.Cert{})
		assert.True(t, manager.GenCertCalled)
		assert.Equal(t, test.GenCertError, err)

		err = manager.GenCSR(pki.Config{}, pki.Claim{}, pki.Cert{})
		assert.True(t, manager.GenCSRCalled)
		assert.Equal(t, test.GenCSRError, err)

		err = manager.SignCSR(pki.Config{}, pki.Cert{}, pki.Config{}, pki.Cert{}, nil)
		assert.True(t, manager.SignCSRCalled)
		assert.Equal(t, test.SignCSRError, err)

		err = manager.VerifyCert(pki.Cert{}, pki.Cert{}, "")
		assert.True(t, manager.VerifyCertCalled)
		assert.Equal(t, test.VerifyCertError, err)
	}
}
