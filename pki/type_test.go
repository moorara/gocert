package pki

import (
	"io/ioutil"
	"strings"
	"testing"

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

func TestLoadState(t *testing.T) {
	tests := []struct {
		yaml           string
		expectError    bool
		expectedRoot   ConfigCA
		expectedInterm ConfigCA
		expectedServer Config
		expectedClient Config
	}{
		{
			``,
			false,
			ConfigCA{},
			ConfigCA{},
			Config{},
			Config{},
		},
		{
			`invalid yaml`,
			true,
			ConfigCA{},
			ConfigCA{},
			Config{},
			Config{},
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
			ConfigCA{
				Config: Config{
					Serial: int64(10),
					Length: 4096,
					Days:   7300,
				},
			},
			ConfigCA{
				Config: Config{
					Serial: int64(100),
					Length: 4096,
					Days:   3650,
				},
			},
			Config{},
			Config{},
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
			ConfigCA{
				Config: Config{
					Serial: int64(10),
					Length: 4096,
					Days:   7300,
				},
			},
			ConfigCA{
				Config: Config{
					Serial: int64(100),
					Length: 4096,
					Days:   3650,
				},
			},
			Config{
				Serial: int64(1000),
				Length: 2048,
				Days:   375,
			},
			Config{
				Serial: int64(10000),
				Length: 2048,
				Days:   40,
			},
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
			&State{
				Root:   ConfigCA{},
				Interm: ConfigCA{},
				Server: Config{},
				Client: Config{},
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
						Length: 4096, Days: 3650,
					},
				},
				Server: Config{},
				Client: Config{},
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

		yaml := strings.Replace(test.expectedYAML, "\t\t\t\t", "  ", -1)
		yaml = strings.Replace(yaml, "\t\t\t", "", -1)
		assert.Equal(t, yaml, string(data))
	}
}

func TestSaveStateError(t *testing.T) {
	err := SaveState(nil, "")
	assert.Error(t, err)
}

func TestNewSpec(t *testing.T) {
	spec := NewSpec()
	rootClaim := Claim{}
	intermClaim := Claim{}
	serverClaim := Claim{}
	clientClaim := Claim{}

	assert.Equal(t, spec.Root, rootClaim)
	assert.Equal(t, spec.Interm, intermClaim)
	assert.Equal(t, spec.Server, serverClaim)
	assert.Equal(t, spec.Client, clientClaim)
}

func TestLoadSpec(t *testing.T) {
	tests := []struct {
		toml           string
		expectError    bool
		expectedRoot   Claim
		expectedInterm Claim
		expectedServer Claim
		expectedClient Claim
	}{
		{
			``,
			false,
			Claim{},
			Claim{},
			Claim{},
			Claim{},
		},
		{
			`invalid toml`,
			true,
			Claim{},
			Claim{},
			Claim{},
			Claim{},
		},
		{
			`
      [root]
				locality = [ "Ottawa" ]
				organization = [ "Moorara" ]
			[server]
				country = [ "US" ]
				organization = [ "AWS" ]
				email_address = [ "moorara@example.com" ]
			`,
			false,
			Claim{
				Locality:     []string{"Ottawa"},
				Organization: []string{"Moorara"},
			},
			Claim{},
			Claim{
				Country:      []string{"US"},
				Organization: []string{"AWS"},
				EmailAddress: []string{"moorara@example.com"},
			},
			Claim{},
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
      `,
			false,
			Claim{
				Country:      []string{"CA", "US"},
				Province:     []string{"Ontario", "Massachusetts"},
				Locality:     []string{"Ottawa", "Boston"},
				Organization: []string{"Moorara"},
			},
			Claim{
				Country:      []string{"CA"},
				Province:     []string{"Ontario"},
				Locality:     []string{"Ottawa"},
				Organization: []string{"Moorara"},
				EmailAddress: []string{"moorara@example.com"},
			},
			Claim{
				Country:      []string{"US"},
				Province:     []string{"Virginia"},
				Locality:     []string{"Richmond"},
				Organization: []string{"Moorara"},
				EmailAddress: []string{"moorara@example.com"},
			},
			Claim{
				Country:      []string{"UK"},
				Locality:     []string{"London"},
				Organization: []string{"Moorara"},
				EmailAddress: []string{"moorara@example.com"},
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
			assert.NotNil(t, spec)
			assert.Equal(t, test.expectedRoot, spec.Root)
			assert.Equal(t, test.expectedInterm, spec.Interm)
			assert.Equal(t, test.expectedServer, spec.Server)
			assert.Equal(t, test.expectedClient, spec.Client)
		}
	}
}

func TestSaveSpec(t *testing.T) {
	tests := []struct {
		spec         *Spec
		expectedTOML string
	}{
		{
			&Spec{
				Root:   Claim{},
				Interm: Claim{},
				Server: Claim{},
				Client: Claim{},
			},
			`[root]

			[intermediate]

			[server]

			[client]
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
			},
			`[root]
				locality = ["Ottawa"]
				organization = ["Moorara"]

			[intermediate]

			[server]
				country = ["US"]
				organization = ["AWS"]

			[client]
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

		toml := strings.Replace(test.expectedTOML, "\t\t\t\t", "  ", -1)
		toml = strings.Replace(toml, "\t\t\t", "", -1)
		assert.Equal(t, toml, string(data))
	}
}

func TestSaveSpecError(t *testing.T) {
	err := SaveSpec(nil, "")
	assert.Error(t, err)
}
