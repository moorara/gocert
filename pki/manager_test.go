package pki

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/moorara/go-box/util"
	"github.com/stretchr/testify/assert"
)

func mockWorkspace(state *State, spec *Spec) (func(), error) {
	items := make([]string, 0)
	deleteFunc := func() {
		for _, item := range items {
			err := os.RemoveAll(item)
			if err != nil {
				log.Printf("Failed to delete %v.", item)
			}
		}
	}

	// Mock sub-directories
	_, err := util.MkDirs("", DirRoot, DirInterm, DirServer, DirClient, DirCSR)
	items = append(items, DirRoot, DirInterm, DirServer, DirClient, DirCSR)
	if err != nil {
		return deleteFunc, err
	}

	// Write state file
	err = SaveState(state, FileState)
	items = append(items, FileState)
	if err != nil {
		return deleteFunc, err
	}

	// Write spec file
	err = SaveSpec(spec, FileSpec)
	items = append(items, FileSpec)
	if err != nil {
		return deleteFunc, err
	}

	return deleteFunc, nil
}

func TestGenRootCA(t *testing.T) {
	tests := []struct {
		name        string
		config      ConfigCA
		claim       Claim
		expectError bool
	}{
		{
			"root",
			ConfigCA{},
			Claim{},
			true,
		},
		{
			"root",
			ConfigCA{
				Config: Config{
					Serial: int64(10),
					Length: 1024,
					Days:   30,
				},
			},
			Claim{
				CommonName:   "Root CA",
				Country:      []string{"CA"},
				Organization: []string{"Moorara"},
			},
			false,
		},
		{
			"root",
			ConfigCA{
				Config: Config{
					Serial: int64(10),
					Length: 1024,
					Days:   30,
				},
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
			false,
		},
	}

	for _, test := range tests {
		state := &State{Root: test.config}
		spec := &Spec{Root: test.claim}
		cleanup, err := mockWorkspace(state, spec)
		assert.NoError(t, err)

		manager := NewX509Manager()
		err = manager.GenRootCA(test.name, test.config, test.claim)

		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

		cleanup()
	}
}

func TestGenRootCAError(t *testing.T) {
	state := NewState()
	spec := NewSpec()
	cleanup, err := mockWorkspace(state, spec)
	assert.NoError(t, err)
	defer cleanup()

	t.Run("No Name", func(t *testing.T) {
		manager := NewX509Manager()
		err := manager.GenRootCA("", ConfigCA{}, Claim{})
		assert.Error(t, err)
	})

	t.Run("Existing Name", func(t *testing.T) {
		certFile := DirRoot + "/root" + extCACert
		keyFile := DirRoot + "/root" + extCAKey
		err := ioutil.WriteFile(certFile, nil, 0644)
		assert.NoError(t, err)
		err = ioutil.WriteFile(keyFile, nil, 0644)
		assert.NoError(t, err)

		manager := NewX509Manager()
		err = manager.GenRootCA("root", ConfigCA{}, Claim{})
		assert.Error(t, err)

		err = util.DeleteAll("", certFile, keyFile)
		assert.NoError(t, err)
	})
}

func TestGenIntermCSR(t *testing.T) {
	tests := []struct {
		name        string
		config      ConfigCA
		claim       Claim
		expectError bool
	}{
		{
			"interm",
			ConfigCA{},
			Claim{},
			true,
		},
		{
			"it",
			ConfigCA{
				Config: Config{
					Serial: 100,
					Length: 1024,
					Days:   7,
				},
			},
			Claim{
				CommonName:         "IT CA",
				Country:            []string{"CA"},
				Organization:       []string{"Moorara"},
				OrganizationalUnit: []string{"IT"},
			},
			false,
		},
		{
			"ops",
			ConfigCA{
				Config: Config{
					Serial: 100,
					Length: 1024,
					Days:   7,
				},
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
			false,
		},
	}

	for _, test := range tests {
		state := &State{Interm: test.config}
		spec := &Spec{Interm: test.claim}
		cleanup, err := mockWorkspace(state, spec)
		assert.NoError(t, err)

		manager := NewX509Manager()
		err = manager.GenIntermCSR(test.name, test.config, test.claim)

		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

		cleanup()
	}
}

func TestGenIntermCSRError(t *testing.T) {
	state := NewState()
	spec := NewSpec()
	cleanup, err := mockWorkspace(state, spec)
	assert.NoError(t, err)
	defer cleanup()

	t.Run("No Name", func(t *testing.T) {
		manager := NewX509Manager()
		err = manager.GenIntermCSR("", ConfigCA{}, Claim{})
		assert.Error(t, err)
	})

	t.Run("Existing Name", func(t *testing.T) {
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
		err = manager.GenIntermCSR("interm", ConfigCA{}, Claim{})
		assert.Error(t, err)

		err = util.DeleteAll("", csrFile, certFile, keyFile)
		assert.NoError(t, err)
	})
}
