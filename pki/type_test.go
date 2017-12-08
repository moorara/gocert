package pki

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/moorara/go-box/util"
	"github.com/stretchr/testify/assert"
)

func TestNewState(t *testing.T) {
	state := NewState()

	assert.Equal(t, defaultRootCASerial, state.Root.Serial)
	assert.Equal(t, defaultRootCALength, state.Root.Length)
	assert.Equal(t, defaultRootCADays, state.Root.Days)

	assert.Equal(t, defaultIntermCASerial, state.Interm.Serial)
	assert.Equal(t, defaultIntermCALength, state.Interm.Length)
	assert.Equal(t, defaultIntermCADays, state.Interm.Days)

	assert.Equal(t, defaultServerCertSerial, state.Server.Serial)
	assert.Equal(t, defaultServerCertLength, state.Server.Length)
	assert.Equal(t, defaultServerCertDays, state.Server.Days)

	assert.Equal(t, defaultClientCertSerial, state.Client.Serial)
	assert.Equal(t, defaultClientCertLength, state.Client.Length)
	assert.Equal(t, defaultClientCertDays, state.Client.Days)
}

func TestLoadState(t *testing.T) {
	tests := []struct {
		yaml          string
		expectError   bool
		expectedState *State
	}{
		{
			`invalid yaml`,
			true,
			&State{},
		},
		{
			``,
			false,
			&State{},
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
			&State{
				Root: ConfigCA{
					Config: Config{
						Serial: int64(10),
						Length: 4096,
						Days:   7300,
					},
				},
				Interm: ConfigCA{
					Config: Config{
						Serial: int64(100),
						Length: 4096,
						Days:   3650,
					},
				},
			},
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
			&State{
				Root: ConfigCA{
					Config: Config{
						Serial: int64(10),
						Length: 4096,
						Days:   7300,
					},
				},
				Interm: ConfigCA{
					Config: Config{
						Serial: int64(100),
						Length: 4096,
						Days:   3650,
					},
				},
				Server: Config{
					Serial: int64(1000),
					Length: 2048,
					Days:   375,
				},
				Client: Config{
					Serial: int64(10000),
					Length: 2048,
					Days:   40,
				},
			},
		},
	}

	for _, test := range tests {
		yaml := strings.Replace(test.yaml, "\t", "  ", -1)
		file, delete, err := util.WriteTempFile(yaml)
		defer delete()
		assert.NoError(t, err)

		state, err := LoadState(file)

		if test.expectError {
			assert.Error(t, err)
			assert.Nil(t, state)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedState, state)
		}
	}
}

func TestLoadStateError(t *testing.T) {
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
			&State{},
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
				Root: ConfigCA{
					Config: Config{
						Serial: 10,
						Length: 4096,
						Days:   7300,
					},
				},
				Interm: ConfigCA{
					Config: Config{
						Serial: 100,
						Length: 4096,
						Days:   3650,
					},
				},
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
				Root: ConfigCA{
					Config: Config{
						Serial: 10,
						Length: 4096,
						Days:   7300,
					},
				},
				Interm: ConfigCA{
					Config: Config{
						Serial: 100,
						Length: 4096,
						Days:   3650,
					},
				},
				Server: Config{
					Serial: 1000,
					Length: 2048,
					Days:   375,
				},
				Client: Config{
					Serial: 10000,
					Length: 2048,
					Days:   40,
				},
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

		yaml := strings.Replace(test.expectedYAML, "\t\t\t", "", -1)
		yaml = strings.Replace(yaml, "\t", "  ", -1)
		assert.Equal(t, yaml, string(data))
	}
}

func TestSaveStateError(t *testing.T) {
	err := SaveState(nil, "")
	assert.Error(t, err)
}

func TestNewSpec(t *testing.T) {
	spec := NewSpec()

	expectedClaim := Claim{}
	expectedRootPolicy := Policy{
		Match:    strings.Split(defaultRootPolicyMatch, ","),
		Supplied: strings.Split(defaultRootPolicySupplied, ","),
	}
	expectedIntermPolicy := Policy{
		Match:    strings.Split(defaultIntermPolicyMatch, ","),
		Supplied: strings.Split(defaultIntermPolicySupplied, ","),
	}

	assert.Equal(t, expectedClaim, spec.Root)
	assert.Equal(t, expectedClaim, spec.Interm)
	assert.Equal(t, expectedClaim, spec.Server)
	assert.Equal(t, expectedClaim, spec.Client)

	assert.Equal(t, expectedRootPolicy, spec.RootPolicy)
	assert.Equal(t, expectedIntermPolicy, spec.IntermPolicy)
}

func TestLoadSpec(t *testing.T) {
	tests := []struct {
		toml         string
		expectError  bool
		expectedSpec *Spec
	}{
		{
			`invalid toml`,
			true,
			&Spec{},
		},
		{
			``,
			false,
			&Spec{},
		},
		{
			`[root]
				locality = [ "Ottawa" ]
				organization = [ "Moorara" ]
			[server]
				country = [ "US" ]
				organization = [ "AWS" ]
				email_address = [ "moorara@example.com" ]
			`,
			false,
			&Spec{
				Root: Claim{
					Locality:     []string{"Ottawa"},
					Organization: []string{"Moorara"},
				},
				Interm: Claim{},
				Server: Claim{
					Country:      []string{"US"},
					Organization: []string{"AWS"},
					EmailAddress: []string{"moorara@example.com"},
				},
				Client: Claim{},
			},
		},
		{
			`
			[root]
				country = [ "CA", "US" ]
				province = [ "Ontario", "Massachusetts" ]
				locality = [ "Ottawa", "Boston" ]
				organization = [ "Moorara" ]
			[intermediate]
				country = [ "CA" ]
				province = [ "Ontario" ]
				locality = [ "Ottawa" ]
				organization = [ "Moorara" ]
				email_address = [ "moorara@example.com" ]
			[server]
				country = [ "US" ]
				province = [ "Virginia" ]
				locality = [ "Richmond" ]
				organization = [ "Moorara" ]
				email_address = [ "moorara@example.com" ]
			[client]
				country = [ "UK" ]
				locality = [ "London" ]
				organization = [ "Moorara" ]
				email_address = [ "moorara@example.com" ]
			[root_policy]
				match = ["Country", "Organization"]
				supplied = ["CommonName"]
			[intermediate_policy]
				match = ["Organization"]
				supplied = ["CommonName"]
			`,
			false,
			&Spec{
				Root: Claim{
					Country:      []string{"CA", "US"},
					Province:     []string{"Ontario", "Massachusetts"},
					Locality:     []string{"Ottawa", "Boston"},
					Organization: []string{"Moorara"},
				},
				Interm: Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Locality:     []string{"Ottawa"},
					Organization: []string{"Moorara"},
					EmailAddress: []string{"moorara@example.com"},
				},
				Server: Claim{
					Country:      []string{"US"},
					Province:     []string{"Virginia"},
					Locality:     []string{"Richmond"},
					Organization: []string{"Moorara"},
					EmailAddress: []string{"moorara@example.com"},
				},
				Client: Claim{
					Country:      []string{"UK"},
					Locality:     []string{"London"},
					Organization: []string{"Moorara"},
					EmailAddress: []string{"moorara@example.com"},
				},
				RootPolicy: Policy{
					Match:    []string{"Country", "Organization"},
					Supplied: []string{"CommonName"},
				},
				IntermPolicy: Policy{
					Match:    []string{"Organization"},
					Supplied: []string{"CommonName"},
				},
			},
		},
	}

	for _, test := range tests {
		file, delete, err := util.WriteTempFile(test.toml)
		defer delete()
		assert.NoError(t, err)

		spec, err := LoadSpec(file)

		if test.expectError {
			assert.Error(t, err)
			assert.Nil(t, spec)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedSpec, spec)
		}
	}
}

func TestSaveSpec(t *testing.T) {
	tests := []struct {
		spec         *Spec
		expectedTOML string
	}{
		{
			&Spec{},
			`[root]

			[intermediate]

			[server]

			[client]

			[root_policy]

			[intermediate_policy]
			`,
		},
		{
			&Spec{
				Root: Claim{
					Locality:     []string{"Ottawa"},
					Organization: []string{"Moorara"},
				},
				Interm: Claim{},
				Server: Claim{
					Country:      []string{"US"},
					Organization: []string{"AWS"}},
				Client: Claim{},
				RootPolicy: Policy{
					Match:    []string{"Organization"},
					Supplied: []string{"CommonName"},
				},
			},
			`[root]
				locality = ["Ottawa"]
				organization = ["Moorara"]

			[intermediate]

			[server]
				country = ["US"]
				organization = ["AWS"]

			[client]

			[root_policy]
				match = ["Organization"]
				supplied = ["CommonName"]

			[intermediate_policy]
			`,
		},
		{
			&Spec{
				Root: Claim{
					Country:      []string{"CA", "US"},
					Province:     []string{"Ontario", "Massachusetts"},
					Locality:     []string{"Ottawa", "Boston"},
					Organization: []string{"Moorara"},
				},
				Interm: Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Locality:     []string{"Ottawa"},
					Organization: []string{"Moorara"},
				},
				Server: Claim{
					Country:      []string{"US"},
					Province:     []string{"Virginia"},
					Locality:     []string{"Richmond"},
					Organization: []string{"Moorara"},
				},
				Client: Claim{
					Country:      []string{"UK"},
					Locality:     []string{"London"},
					Organization: []string{"Moorara"},
				},
				RootPolicy: Policy{
					Match:    []string{"Country", "Organization"},
					Supplied: []string{"CommonName"},
				},
				IntermPolicy: Policy{
					Match:    []string{"Organization"},
					Supplied: []string{"CommonName"},
				},
			},
			`[root]
				country = ["CA", "US"]
				province = ["Ontario", "Massachusetts"]
				locality = ["Ottawa", "Boston"]
				organization = ["Moorara"]

			[intermediate]
				country = ["CA"]
				province = ["Ontario"]
				locality = ["Ottawa"]
				organization = ["Moorara"]

			[server]
				country = ["US"]
				province = ["Virginia"]
				locality = ["Richmond"]
				organization = ["Moorara"]

			[client]
				country = ["UK"]
				locality = ["London"]
				organization = ["Moorara"]

			[root_policy]
				match = ["Country", "Organization"]
				supplied = ["CommonName"]

			[intermediate_policy]
				match = ["Organization"]
				supplied = ["CommonName"]
			`,
		},
	}

	for _, test := range tests {
		file, delete, err := util.WriteTempFile("")
		defer delete()
		assert.NoError(t, err)

		err = SaveSpec(test.spec, file)
		assert.NoError(t, err)

		data, err := ioutil.ReadFile(file)
		assert.NoError(t, err)

		toml := strings.Replace(test.expectedTOML, "\t\t\t", "", -1)
		toml = strings.Replace(toml, "\t", "  ", -1)
		assert.Equal(t, toml, string(data))
	}
}

func TestSaveSpecError(t *testing.T) {
	err := SaveSpec(nil, "")
	assert.Error(t, err)
}
