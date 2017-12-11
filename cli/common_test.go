package cli

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/moorara/go-box/util"
	"github.com/moorara/gocert/pki"
	"github.com/stretchr/testify/assert"
)

type mockedManager struct {
	GenCertError error
	GenCSRError  error
	SignCSRError error

	GenCertCalled bool
	GenCSRCalled  bool
	SignCSRCalled bool
}

func (m *mockedManager) GenCert(pki.Config, pki.Claim, pki.Metadata) error {
	m.GenCertCalled = true
	return m.GenCertError
}

func (m *mockedManager) GenCSR(pki.Config, pki.Claim, pki.Metadata) error {
	m.GenCSRCalled = true
	return m.GenCSRError
}

func (m *mockedManager) SignCSR(pki.Config, pki.Metadata, pki.Config, pki.Metadata, pki.TrustFunc) error {
	m.SignCSRCalled = true
	return m.SignCSRError
}

func TestNewColoredUi(t *testing.T) {
	tests := []struct {
		in, out, er string
	}{
		{"in", "out", "err"},
		{"input", "output", "error"},
	}

	for _, test := range tests {
		inR, inW, outR, outW, errR, errW, restore, err := util.PipeStdAll()
		defer restore()
		assert.NoError(t, err)

		ui := newColoredUI()

		assert.Equal(t, cli.UiColorNone, ui.OutputColor)
		assert.Equal(t, cli.UiColorGreen, ui.InfoColor)
		assert.Equal(t, cli.UiColorRed, ui.ErrorColor)
		assert.Equal(t, cli.UiColorYellow, ui.WarnColor)

		in, err := util.WriteToStdinPipe(inR, inW, test.in)
		assert.NoError(t, err)
		assert.Equal(t, test.in, in)

		out, err := util.ReadFromStdoutPipe(outR, outW, test.out)
		assert.NoError(t, err)
		assert.Equal(t, test.out, out)

		er, err := util.ReadFromStderrPipe(errR, errW, test.er)
		assert.NoError(t, err)
		assert.Equal(t, test.er, er)
	}
}

func TestLoadWorkspace(t *testing.T) {
	tests := []struct {
		stateYAML      string
		specTOML       string
		expectedStatus int
		expectedState  *pki.State
		expectedSpec   *pki.Spec
	}{
		{
			``,
			``,
			0,
			&pki.State{},
			&pki.Spec{},
		},
		{
			`invalid yaml`,
			``,
			ErrorReadState,
			&pki.State{},
			&pki.Spec{},
		},
		{
			``,
			`invalid toml`,
			ErrorReadSpec,
			&pki.State{},
			&pki.Spec{},
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
			0,
			&pki.State{
				Root: pki.Config{
					Serial: int64(10),
					Length: 4096,
					Days:   7300,
				},
				Interm: pki.Config{
					Serial: int64(100),
					Length: 4096,
					Days:   3650,
				},
				Server: pki.Config{
					Serial: int64(1000),
					Length: 2048,
					Days:   375,
				},
				Client: pki.Config{
					Serial: int64(10000),
					Length: 2048,
					Days:   40,
				},
			},
			&pki.Spec{
				Root: pki.Claim{
					Country:      []string{"CA", "US"},
					Province:     []string{"Ontario", "Massachusetts"},
					Locality:     []string{"Ottawa", "Boston"},
					Organization: []string{"Moorara"},
				},
				Interm: pki.Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Locality:     []string{"Ottawa"},
					Organization: []string{"Moorara"},
					EmailAddress: []string{"moorara@example.com"},
				},
				Server: pki.Claim{
					Country:      []string{"US"},
					Province:     []string{"Virginia"},
					Locality:     []string{"Richmond"},
					Organization: []string{"Moorara"},
					EmailAddress: []string{"moorara@example.com"},
				},
				Client: pki.Claim{
					Country:      []string{"UK"},
					Locality:     []string{"London"},
					Organization: []string{"Moorara"},
					EmailAddress: []string{"moorara@example.com"},
				},
				RootPolicy: pki.Policy{
					Match:    []string{"Country", "Organization"},
					Supplied: []string{"CommonName"},
				},
				IntermPolicy: pki.Policy{
					Match:    []string{"Organization"},
					Supplied: []string{"CommonName"},
				},
			},
		},
	}

	for _, test := range tests {
		stateYAML := strings.Replace(test.stateYAML, "\t", "  ", -1)
		err := ioutil.WriteFile(pki.FileState, []byte(stateYAML), 0644)
		assert.NoError(t, err)
		err = ioutil.WriteFile(pki.FileSpec, []byte(test.specTOML), 0644)
		assert.NoError(t, err)

		mockUI := cli.NewMockUi()
		state, spec, status := loadWorkspace(mockUI)

		assert.Equal(t, test.expectedStatus, status)
		if test.expectedStatus == 0 {
			assert.Equal(t, test.expectedState, state)
			assert.Equal(t, test.expectedSpec, spec)
		}

		err = util.DeleteAll("", pki.FileState, pki.FileSpec)
		assert.NoError(t, err)
	}
}

func TestSaveWorkspace(t *testing.T) {
	tests := []struct {
		state             *pki.State
		spec              *pki.Spec
		expectedStatus    int
		expectedStateYAML string
		expectedSpecTOML  string
	}{
		{
			&pki.State{},
			&pki.Spec{},
			0,
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
			`[root]

			[intermediate]

			[server]

			[client]

			[root_policy]

			[intermediate_policy]
			`,
		},
		{
			&pki.State{
				Root: pki.Config{
					Serial: 10,
					Length: 4096,
					Days:   7300,
				},
				Interm: pki.Config{
					Serial: 100,
					Length: 4096,
					Days:   3650,
				},
				Server: pki.Config{
					Serial: 1000,
					Length: 2048,
					Days:   375,
				},
				Client: pki.Config{
					Serial: 10000,
					Length: 2048,
					Days:   40,
				},
			},
			&pki.Spec{
				Root: pki.Claim{
					Country:      []string{"CA", "US"},
					Province:     []string{"Ontario", "Massachusetts"},
					Locality:     []string{"Ottawa", "Boston"},
					Organization: []string{"Moorara"},
				},
				Interm: pki.Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Locality:     []string{"Ottawa"},
					Organization: []string{"Moorara"},
				},
				Server: pki.Claim{
					Country:      []string{"US"},
					Province:     []string{"Virginia"},
					Locality:     []string{"Richmond"},
					Organization: []string{"Moorara"},
				},
				Client: pki.Claim{
					Country:      []string{"UK"},
					Locality:     []string{"London"},
					Organization: []string{"Moorara"},
				},
				RootPolicy: pki.Policy{
					Match:    []string{"Country", "Organization"},
					Supplied: []string{"CommonName"},
				},
				IntermPolicy: pki.Policy{
					Match:    []string{"Organization"},
					Supplied: []string{"CommonName"},
				},
			},
			0,
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
		mockUI := cli.NewMockUi()
		status := saveWorkspace(test.state, test.spec, mockUI)

		assert.Equal(t, test.expectedStatus, status)
		if test.expectedStatus == 0 {
			stateYAML, err := ioutil.ReadFile(pki.FileState)
			assert.NoError(t, err)
			expectedStateYAML := strings.Replace(test.expectedStateYAML, "\t\t\t\t", "  ", -1)
			expectedStateYAML = strings.Replace(expectedStateYAML, "\t\t\t", "", -1)
			assert.Equal(t, expectedStateYAML, string(stateYAML))

			specTOML, err := ioutil.ReadFile(pki.FileSpec)
			assert.NoError(t, err)
			expectedSpecTOML := strings.Replace(test.expectedSpecTOML, "\t\t\t\t", "  ", -1)
			expectedSpecTOML = strings.Replace(expectedSpecTOML, "\t\t\t", "", -1)
			assert.Equal(t, expectedSpecTOML, string(specTOML))
		}

		err := util.DeleteAll("", pki.FileState, pki.FileSpec)
		assert.NoError(t, err)
	}
}

func TestAskForNewState(t *testing.T) {
	tests := []struct {
		input         string
		expectedState pki.State
	}{
		{
			``,
			pki.State{
				Root:   pki.Config{},
				Interm: pki.Config{},
				Server: pki.Config{},
				Client: pki.Config{},
			},
		},
		{
			`10
			4096
			7300
			100
			4096
			3650
			1000
			2048
			375
			10000
			2048
			40
			`,
			pki.State{
				Root: pki.Config{
					Serial: int64(10),
					Length: 4096,
					Days:   7300,
				},
				Interm: pki.Config{
					Serial: int64(100),
					Length: 4096,
					Days:   3650,
				},
				Server: pki.Config{
					Serial: int64(1000),
					Length: 2048, Days: 375,
				},
				Client: pki.Config{
					Serial: int64(10000),
					Length: 2048,
					Days:   40,
				},
			},
		},
	}

	for _, test := range tests {
		mockUI := cli.NewMockUi()
		mockUI.InputReader = strings.NewReader(test.input)
		state := askForNewState(mockUI)

		assert.Equal(t, test.expectedState, *state)
	}
}

func TestAskForNewSpec(t *testing.T) {
	tests := []struct {
		input        string
		expectedSpec pki.Spec
	}{
		{
			``,
			pki.Spec{},
		},
		{
			`CA
			Ontario

			Milad
			`,
			pki.Spec{
				Root: pki.Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Organization: []string{"Milad"},
				},
				Interm: pki.Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Organization: []string{"Milad"},
				},
				Server: pki.Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Organization: []string{"Milad"},
				},
				Client: pki.Claim{
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




			Country,Organization
      CommonName
      Organization
      CommonName
			`,
			pki.Spec{
				Root: pki.Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Organization: []string{"Milad"},
				},
				Interm: pki.Claim{
					Country:            []string{"CA"},
					Province:           []string{"Ontario"},
					Locality:           []string{"Ottawa"},
					Organization:       []string{"Milad"},
					OrganizationalUnit: []string{"R&D"},
				},
				Server: pki.Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Locality:     []string{"Toronto", "Montreal"},
					Organization: []string{"Milad"},
				},
				Client: pki.Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Locality:     []string{"Ottawa"},
					Organization: []string{"Milad"},
				},
				RootPolicy: pki.Policy{
					Match:    []string{"Country", "Organization"},
					Supplied: []string{"CommonName"},
				},
				IntermPolicy: pki.Policy{
					Match:    []string{"Organization"},
					Supplied: []string{"CommonName"},
				},
			},
		},
	}

	for _, test := range tests {
		mockUI := cli.NewMockUi()
		mockUI.InputReader = strings.NewReader(test.input)
		spec := askForNewSpec(mockUI)

		assert.Equal(t, test.expectedSpec, *spec)
	}
}

func TestAskForConfig(t *testing.T) {
	tests := []struct {
		config         pki.Config
		input          string
		expectedConfig pki.Config
	}{
		{
			pki.Config{},
			``,
			pki.Config{},
		},
		{
			pki.Config{
				Length: 2048,
			},
			`1000
			375
			`,
			pki.Config{
				Serial: int64(1000),
				Length: 2048,
				Days:   375,
			},
		},
	}

	for _, test := range tests {
		mockUI := cli.NewMockUi()
		mockUI.InputReader = strings.NewReader(test.input)
		askForConfig(&test.config, mockUI)

		assert.Equal(t, test.expectedConfig, test.config)
	}
}

func TestAskForConfigCA(t *testing.T) {
	tests := []struct {
		config           pki.Config
		input            string
		expectedConfigCA pki.Config
	}{
		{
			pki.Config{},
			``,
			pki.Config{},
		},
		{
			pki.Config{},
			`100
			4096
			3650
			secret
			secret
			`,
			pki.Config{
				Serial:   int64(100),
				Length:   4096,
				Days:     3650,
				Password: "secret",
			},
		},
	}

	for _, test := range tests {
		mockUI := cli.NewMockUi()
		mockUI.InputReader = strings.NewReader(test.input)
		askForConfig(&test.config, mockUI)

		assert.Equal(t, test.expectedConfigCA, test.config)
	}
}

func TestAskForClaim(t *testing.T) {
	tests := []struct {
		claim         pki.Claim
		input         string
		expectedClaim pki.Claim
	}{
		{
			pki.Claim{},
			``,
			pki.Claim{},
		},
		{
			pki.Claim{
				Country:      []string{"CA"},
				Province:     []string{"Ontario"},
				Organization: []string{"Milad"},
			},
			`IntermediateCA
			Ottawa
			`,
			pki.Claim{
				CommonName:   "IntermediateCA",
				Country:      []string{"CA"},
				Province:     []string{"Ontario"},
				Locality:     []string{"Ottawa"},
				Organization: []string{"Milad"},
			},
		},
	}

	for _, test := range tests {
		mockUI := cli.NewMockUi()
		mockUI.InputReader = strings.NewReader(test.input)
		askForClaim(&test.claim, mockUI)

		assert.Equal(t, test.expectedClaim, test.claim)
	}
}
