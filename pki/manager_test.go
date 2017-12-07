package pki

import (
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
		config      ConfigCA
		claim       Claim
		expectError bool
	}{
		{
			ConfigCA{},
			Claim{},
			true,
		},
		{
			ConfigCA{
				Config: Config{
					Serial: int64(10),
					Length: 1024,
					Days:   7,
				},
			},
			Claim{
				Country:      []string{"CA"},
				Organization: []string{"Moorara"},
			},
			false,
		},
		{
			ConfigCA{
				Config: Config{
					Serial: int64(10),
					Length: 1024,
					Days:   7,
				},
				Password: "secret",
			},
			Claim{
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
		err = manager.GenRootCA(test.config, test.claim)

		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

		cleanup()
	}
}
