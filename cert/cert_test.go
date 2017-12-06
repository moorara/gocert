package cert

import (
	"log"
	"os"
	"testing"

	"github.com/moorara/go-box/util"
	"github.com/moorara/gocert/config"
	"github.com/stretchr/testify/assert"
)

func mockWorkspace(state *config.State, spec *config.Spec) (func(), error) {
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
	_, err := util.MkDirs("", config.DirNameRoot, config.DirNameInterm, config.DirNameServer, config.DirNameClient)
	items = append(items, config.DirNameRoot, config.DirNameInterm, config.DirNameServer, config.DirNameClient)
	if err != nil {
		return deleteFunc, err
	}

	// Write state file
	err = config.SaveState(state, config.FileNameState)
	items = append(items, config.FileNameState)
	if err != nil {
		return deleteFunc, err
	}

	// Write spec file
	err = config.SaveSpec(spec, config.FileNameSpec)
	items = append(items, config.FileNameSpec)
	if err != nil {
		return deleteFunc, err
	}

	return deleteFunc, nil
}

func TestGenRootCA(t *testing.T) {
	tests := []struct {
		settings    config.SettingsCA
		claim       config.Claim
		expectError bool
	}{
		{
			config.SettingsCA{},
			config.Claim{},
			true,
		},
		{
			config.SettingsCA{
				Settings: config.Settings{
					Serial: int64(10),
					Length: 1024,
					Days:   7,
				},
			},
			config.Claim{
				Country:      []string{"CA"},
				Organization: []string{"Moorara"},
			},
			false,
		},
		{
			config.SettingsCA{
				Settings: config.Settings{
					Serial: int64(10),
					Length: 1024,
					Days:   7,
				},
				Password: "secret",
			},
			config.Claim{
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
		state := &config.State{Root: test.settings}
		spec := &config.Spec{Root: test.claim}
		cleanup, err := mockWorkspace(state, spec)
		assert.NoError(t, err)

		manager := NewX509Manager()
		err = manager.GenRootCA(test.settings, test.claim)

		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

		cleanup()
	}
}
