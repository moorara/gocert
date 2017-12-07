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
	GenRootCAError     error
	GenIntermCAError   error
	GenServerCertError error
	GenClientCertError error

	GenRootCACalled     bool
	GenIntermCACalled   bool
	GenServerCertCalled bool
	GenClientCertCalled bool
}

func (m *mockedManager) GenRootCA(pki.ConfigCA, pki.Claim) error {
	m.GenRootCACalled = true
	return m.GenRootCAError
}

func (m *mockedManager) GenIntermCA(pki.ConfigCA, pki.Claim) error {
	m.GenIntermCACalled = true
	return m.GenIntermCAError
}

func (m *mockedManager) GenServerCert(pki.Config, pki.Claim) error {
	m.GenServerCertCalled = true
	return m.GenServerCertError
}

func (m *mockedManager) GenClientCert(pki.Config, pki.Claim) error {
	m.GenClientCertCalled = true
	return m.GenClientCertError
}

func mockWorkspace(state *pki.State, spec *pki.Spec) (func() error, error) {
	items := make([]string, 0)
	deleteFunc := func() error {
		return util.DeleteAll("", items...)
	}

	// Mock sub-directories
	_, err := util.MkDirs("", pki.DirRoot, pki.DirInterm, pki.DirServer, pki.DirClient, pki.DirCSR)
	items = append(items, pki.DirRoot, pki.DirInterm, pki.DirServer, pki.DirClient, pki.DirCSR)
	if err != nil {
		return deleteFunc, err
	}

	// Mock state file
	if state != nil {
		err := pki.SaveState(state, pki.FileState)
		if err != nil {
			return deleteFunc, err
		}
		items = append(items, pki.FileState)
	}

	// Mock spec file
	if spec != nil {
		err := pki.SaveSpec(spec, pki.FileSpec)
		if err != nil {
			return deleteFunc, err
		}
		items = append(items, pki.FileSpec)
	}

	return deleteFunc, nil
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
      `,
			0,
			&pki.State{
				Root: pki.ConfigCA{
					Config: pki.Config{
						Serial: int64(10),
						Length: 4096,
						Days:   7300,
					},
				},
				Interm: pki.ConfigCA{
					Config: pki.Config{
						Serial: int64(100),
						Length: 4096,
						Days:   3650,
					},
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
				}},
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
				}},
		},
	}

	for _, test := range tests {
		stateYAML := strings.Replace(test.stateYAML, "\t", "  ", -1)
		err := ioutil.WriteFile(pki.FileState, []byte(stateYAML), 0644)
		assert.NoError(t, err)
		err = ioutil.WriteFile(pki.FileSpec, []byte(test.specTOML), 0644)
		assert.NoError(t, err)

		mockUI := cli.NewMockUi()
		state, spec, status := LoadWorkspace(mockUI)

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
			`,
		},
		{
			&pki.State{
				Root: pki.ConfigCA{
					Config: pki.Config{
						Serial: 10,
						Length: 4096,
						Days:   7300,
					},
				},
				Interm: pki.ConfigCA{
					Config: pki.Config{
						Serial: 100,
						Length: 4096,
						Days:   3650,
					},
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
			`,
		},
	}

	for _, test := range tests {
		mockUI := cli.NewMockUi()
		status := SaveWorkspace(test.state, test.spec, mockUI)

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
				Root:   pki.ConfigCA{},
				Interm: pki.ConfigCA{},
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
				Root: pki.ConfigCA{
					Config: pki.Config{
						Serial: int64(10),
						Length: 4096,
						Days:   7300,
					},
				},
				Interm: pki.ConfigCA{
					Config: pki.Config{
						Serial: int64(100),
						Length: 4096,
						Days:   3650,
					},
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
		state := AskForNewState(mockUI)

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
			pki.Spec{
				Root:   pki.Claim{},
				Interm: pki.Claim{},
				Server: pki.Claim{},
				Client: pki.Claim{},
			},
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
			},
		},
	}

	for _, test := range tests {
		mockUI := cli.NewMockUi()
		mockUI.InputReader = strings.NewReader(test.input)
		spec := AskForNewSpec(mockUI)

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
		AskForConfig(&test.config, mockUI)

		assert.Equal(t, test.expectedConfig, test.config)
	}
}

func TestAskForConfigCA(t *testing.T) {
	tests := []struct {
		configCA         pki.ConfigCA
		input            string
		expectedConfigCA pki.ConfigCA
	}{
		{
			pki.ConfigCA{},
			``,
			pki.ConfigCA{},
		},
		{
			pki.ConfigCA{},
			`100
			4096
			3650
			secret
			secret
			`,
			pki.ConfigCA{
				Config: pki.Config{
					Serial: int64(100),
					Length: 4096,
					Days:   3650,
				},
				Password: "secret",
			},
		},
	}

	for _, test := range tests {
		mockUI := cli.NewMockUi()
		mockUI.InputReader = strings.NewReader(test.input)
		AskForConfigCA(&test.configCA, mockUI)

		assert.Equal(t, test.expectedConfigCA, test.configCA)
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
		AskForClaim(&test.claim, mockUI)

		assert.Equal(t, test.expectedClaim, test.claim)
	}
}
