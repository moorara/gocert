package config

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/moorara/go-box/util"
	"github.com/stretchr/testify/assert"
)

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

func TestNewSpecWithInput(t *testing.T) {
	tests := []struct {
		input        string
		expectedSpec Spec
	}{
		{
			``,
			Spec{
				Root:   Claim{},
				Interm: Claim{},
				Server: Claim{},
				Client: Claim{},
			},
		},
		{
			`CA
			Ontario

			Milad
			`,
			Spec{
				Root: Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Organization: []string{"Milad"},
				},
				Interm: Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Organization: []string{"Milad"},
				},
				Server: Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Organization: []string{"Milad"},
				},
				Client: Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Organization: []string{"Milad"},
				},
			},
		},
		{
			`CA
			Ontario

			Milad









			Ottawa
			R&D



			Toronto,Montreal




			Ottawa
			`,
			Spec{
				Root: Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Organization: []string{"Milad"},
				},
				Interm: Claim{
					Country:            []string{"CA"},
					Province:           []string{"Ontario"},
					Locality:           []string{"Ottawa"},
					Organization:       []string{"Milad"},
					OrganizationalUnit: []string{"R&D"},
				},
				Server: Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Locality:     []string{"Toronto", "Montreal"},
					Organization: []string{"Milad"},
				},
				Client: Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Locality:     []string{"Ottawa"},
					Organization: []string{"Milad"},
				},
			},
		},
	}

	for _, test := range tests {
		mockUI := cli.NewMockUi()
		mockUI.InputReader = strings.NewReader(test.input)
		spec := NewSpecWithInput(mockUI)

		assert.Equal(t, test.expectedSpec, *spec)
	}
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
