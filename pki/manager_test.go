package pki

import (
	"io/ioutil"
	"testing"

	"github.com/moorara/go-box/util"
	"github.com/stretchr/testify/assert"
)

func TestGenCert(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		claim       Claim
		md          Metadata
		expectError bool
	}{
		{
			"root",
			Config{},
			Claim{},
			Metadata{},
			true,
		},
		{
			"root",
			Config{
				Serial: int64(10),
				Length: 1024,
				Days:   30,
			},
			Claim{
				CommonName:   "Root CA",
				Country:      []string{"CA"},
				Organization: []string{"Moorara"},
			},
			Metadata{
				CertType: CertTypeRoot,
			},
			false,
		},
		{
			"root",
			Config{
				Serial:   int64(10),
				Length:   1024,
				Days:     30,
				Password: "secret",
			},
			Claim{
				CommonName:   "Root CA",
				Country:      []string{"CA"},
				Province:     []string{"Ontario"},
				Locality:     []string{"Ottawa"},
				Organization: []string{"Moorara"},
				EmailAddress: []string{"moorara@example.com"},
			},
			Metadata{
				CertType: CertTypeRoot,
			},
			false,
		},
	}

	for _, test := range tests {
		state := &State{Root: test.config}
		spec := &Spec{Root: test.claim}
		err := NewWorkspace(state, spec)
		assert.NoError(t, err)

		manager := NewX509Manager()
		err = manager.GenCert(test.name, test.config, test.claim, test.md)

		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			// TODO: verify
		}

		CleanupWorkspace()
	}
}

func TestGenCertError(t *testing.T) {
	state := NewState()
	spec := NewSpec()
	err := NewWorkspace(state, spec)
	assert.NoError(t, err)
	defer CleanupWorkspace()

	t.Run("NoName", func(t *testing.T) {
		manager := NewX509Manager()
		err := manager.GenCert("", Config{}, Claim{}, Metadata{})
		assert.Error(t, err)
	})

	t.Run("ExistingName", func(t *testing.T) {
		certFile := DirRoot + "/root" + extCACert
		keyFile := DirRoot + "/root" + extCAKey
		err := ioutil.WriteFile(certFile, nil, 0644)
		assert.NoError(t, err)
		err = ioutil.WriteFile(keyFile, nil, 0644)
		assert.NoError(t, err)

		manager := NewX509Manager()
		err = manager.GenCert("root", Config{}, Claim{}, Metadata{})
		assert.Error(t, err)

		err = util.DeleteAll("", certFile, keyFile)
		assert.NoError(t, err)
	})
}

func TestGenCSR(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		claim       Claim
		md          Metadata
		expectError bool
	}{
		{
			"interm",
			Config{},
			Claim{},
			Metadata{},
			true,
		},
		{
			"it",
			Config{
				Serial: 100,
				Length: 1024,
				Days:   7,
			},
			Claim{
				CommonName:         "IT CA",
				Country:            []string{"CA"},
				Organization:       []string{"Moorara"},
				OrganizationalUnit: []string{"IT"},
			},
			Metadata{
				CertType: CertTypeInterm,
			},
			false,
		},
		{
			"ops",
			Config{
				Serial:   100,
				Length:   1024,
				Days:     7,
				Password: "secret",
			},
			Claim{
				CommonName:         "Ops CA",
				Country:            []string{"CA"},
				Province:           []string{"Ontario"},
				Locality:           []string{"Ottawa"},
				Organization:       []string{"Moorara"},
				OrganizationalUnit: []string{"Ops"},
				EmailAddress:       []string{"moorara@example.com"},
			},
			Metadata{
				CertType: CertTypeInterm,
			},
			false,
		},
	}

	for _, test := range tests {
		state := &State{Interm: test.config}
		spec := &Spec{Interm: test.claim}
		err := NewWorkspace(state, spec)
		assert.NoError(t, err)

		manager := NewX509Manager()
		err = manager.GenCSR(test.name, test.config, test.claim, test.md)

		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			// TODO: verify
		}

		CleanupWorkspace()
	}
}

func TestGenCSRError(t *testing.T) {
	state := NewState()
	spec := NewSpec()
	err := NewWorkspace(state, spec)
	assert.NoError(t, err)
	defer CleanupWorkspace()

	t.Run("NoName", func(t *testing.T) {
		manager := NewX509Manager()
		err = manager.GenCSR("", Config{}, Claim{}, Metadata{})
		assert.Error(t, err)
	})

	t.Run("ExistingName", func(t *testing.T) {
		csrFile := DirInterm + "/interm" + extCACSR
		certFile := DirInterm + "/interm" + extCACert
		keyFile := DirInterm + "/interm" + extCAKey
		err := ioutil.WriteFile(csrFile, nil, 0644)
		assert.NoError(t, err)
		err = ioutil.WriteFile(certFile, nil, 0644)
		assert.NoError(t, err)
		err = ioutil.WriteFile(keyFile, nil, 0644)
		assert.NoError(t, err)

		manager := NewX509Manager()
		err = manager.GenCSR("interm", Config{}, Claim{}, Metadata{})
		assert.Error(t, err)

		err = util.DeleteAll("", csrFile, certFile, keyFile)
		assert.NoError(t, err)
	})
}
