package config

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/moorara/go-box/util"
	"github.com/stretchr/testify/assert"
)

func TestNewState(t *testing.T) {
	st := NewState()

	assert.Equal(t, st.Root.Serial, defaultRootCASerial)
	assert.Equal(t, st.Root.Length, defaultRootCALength)
	assert.Equal(t, st.Root.Days, defaultRootCADays)

	assert.Equal(t, st.Interm.Serial, defaultIntermCASerial)
	assert.Equal(t, st.Interm.Length, defaultIntermCALength)
	assert.Equal(t, st.Interm.Days, defaultIntermCADays)

	assert.Equal(t, st.Server.Serial, defaultServerCertSerial)
	assert.Equal(t, st.Server.Length, defaultServerCertLength)
	assert.Equal(t, st.Server.Days, defaultServerCertDays)

	assert.Equal(t, st.Client.Serial, defaultClientCertSerial)
	assert.Equal(t, st.Client.Length, defaultClientCertLength)
	assert.Equal(t, st.Client.Days, defaultClientCertDays)
}

func TestNewStateWithInput(t *testing.T) {
	tests := []struct {
		input         string
		expectedState State
	}{
		{
			``,
			State{
				Root:   Settings{},
				Interm: Settings{},
				Server: Settings{},
				Client: Settings{},
			},
		},
		{
			`10
			4098
			7300
			100
			4098
			3650
			1000
			2048
			375
			10000
			2048
			40`,
			State{
				Root:   Settings{Serial: int64(10), Length: 4098, Days: 7300},
				Interm: Settings{Serial: int64(100), Length: 4098, Days: 3650},
				Server: Settings{Serial: int64(1000), Length: 2048, Days: 375},
				Client: Settings{Serial: int64(10000), Length: 2048, Days: 40},
			},
		},
	}

	for _, test := range tests {
		mockUI := cli.NewMockUi()
		mockUI.InputReader = strings.NewReader(test.input)
		state := NewStateWithInput(mockUI)

		assert.Equal(t, test.expectedState, *state)
	}
}

func TestLoadState(t *testing.T) {
	tests := []struct {
		yaml           string
		expectError    bool
		expectedRoot   Settings
		expectedInterm Settings
		expectedServer Settings
		expectedClient Settings
	}{
		{
			``,
			false,
			Settings{},
			Settings{},
			Settings{},
			Settings{},
		},
		{
			`invalid yaml`,
			true,
			Settings{},
			Settings{},
			Settings{},
			Settings{},
		},
		{
			`
			root:
				serial: 10
				length: 4096
				days: 7300
			intermediate:
				serial: 100
				length: 4096
				days: 3650
			`,
			false,
			Settings{Serial: int64(10), Length: 4096, Days: 7300},
			Settings{Serial: int64(100), Length: 4096, Days: 3650},
			Settings{},
			Settings{},
		},
		{
			`
			root:
				serial: 10
				length: 4096
				days: 7300
			intermediate:
				serial: 100
				length: 4096
				days: 3650
			server:
				serial: 1000
				length: 2048
				days: 375
			client:
				serial: 10000
				length: 2048
				days: 40
			`,
			false,
			Settings{Serial: int64(10), Length: 4096, Days: 7300},
			Settings{Serial: int64(100), Length: 4096, Days: 3650},
			Settings{Serial: int64(1000), Length: 2048, Days: 375},
			Settings{Serial: int64(10000), Length: 2048, Days: 40},
		},
	}

	for _, test := range tests {
		yaml := strings.Replace(test.yaml, "\t", "  ", -1)
		file, delete, err := util.WriteTempFile(yaml)
		defer delete()
		assert.NoError(t, err)

		spec, err := LoadState(file)

		if test.expectError {
			assert.Error(t, err)
			assert.Nil(t, spec)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, spec)
			assert.Equal(t, test.expectedRoot, spec.Root)
			assert.Equal(t, test.expectedInterm, spec.Interm)
			assert.Equal(t, test.expectedServer, spec.Server)
			assert.Equal(t, test.expectedClient, spec.Client)
		}
	}
}

func TestLoadStateWithError(t *testing.T) {
	spec, err := LoadState("")
	assert.Error(t, err)
	assert.Nil(t, spec)
}

func TestSaveState(t *testing.T) {
	tests := []struct {
		state        *State
		expectedYAML string
	}{
		{
			&State{
				Root:   Settings{},
				Interm: Settings{},
				Server: Settings{},
				Client: Settings{},
			},
			`root:
				serial: 0
				length: 0
				days: 0
			intermediate:
				serial: 0
				length: 0
				days: 0
			server:
				serial: 0
				length: 0
				days: 0
			client:
				serial: 0
				length: 0
				days: 0
			`,
		},
		{
			&State{
				Root:   Settings{Serial: 10, Length: 4096, Days: 7300},
				Interm: Settings{Serial: 100, Length: 4096, Days: 3650},
				Server: Settings{},
				Client: Settings{},
			},
			`root:
				serial: 10
				length: 4096
				days: 7300
			intermediate:
				serial: 100
				length: 4096
				days: 3650
			server:
				serial: 0
				length: 0
				days: 0
			client:
				serial: 0
				length: 0
				days: 0
			`,
		},
		{
			&State{
				Root:   Settings{Serial: 10, Length: 4096, Days: 7300},
				Interm: Settings{Serial: 100, Length: 4096, Days: 3650},
				Server: Settings{Serial: 1000, Length: 2048, Days: 375},
				Client: Settings{Serial: 10000, Length: 2048, Days: 40},
			},
			`root:
				serial: 10
				length: 4096
				days: 7300
			intermediate:
				serial: 100
				length: 4096
				days: 3650
			server:
				serial: 1000
				length: 2048
				days: 375
			client:
				serial: 10000
				length: 2048
				days: 40
			`,
		},
	}

	for _, test := range tests {
		file, delete, err := util.WriteTempFile("")
		defer delete()
		assert.NoError(t, err)

		err = SaveState(test.state, file)
		assert.NoError(t, err)

		data, err := ioutil.ReadFile(file)
		assert.NoError(t, err)

		yaml := strings.Replace(test.expectedYAML, "\t\t\t\t", "  ", -1)
		yaml = strings.Replace(yaml, "\t\t\t", "", -1)
		assert.Equal(t, yaml, string(data))
	}
}

func TestSaveStateWithError(t *testing.T) {
	err := SaveState(nil, "")
	assert.Error(t, err)
}

func TestSettingsFillIn(t *testing.T) {
	tests := []struct {
		settings         Settings
		input            string
		expectedSettings Settings
	}{
		{
			Settings{},
			``,
			Settings{},
		},
		{
			Settings{
				Length: 2048,
			},
			`1000
			375
			`,
			Settings{
				Serial: int64(1000),
				Length: 2048,
				Days:   375,
			},
		},
	}

	for _, test := range tests {
		mockUI := cli.NewMockUi()
		mockUI.InputReader = strings.NewReader(test.input)
		test.settings.FillIn(mockUI)

		assert.Equal(t, test.expectedSettings, test.settings)
	}
}
