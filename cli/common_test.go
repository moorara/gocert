package cli

import (
	"net"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/moorara/gocert/pki"
	"github.com/moorara/gocert/util"
	"github.com/stretchr/testify/assert"
)

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
		title          string
		stateFixture   string
		specFixture    string
		expectedStatus int
		expectedState  *pki.State
		expectedSpec   *pki.Spec
	}{
		{
			title:          "Empty",
			stateFixture:   "./fixture/loadWorkspace/empty.yaml",
			specFixture:    "./fixture/loadWorkspace/empty.toml",
			expectedStatus: 0,
			expectedState:  &pki.State{},
			expectedSpec:   &pki.Spec{},
		},
		{
			title:          "InvalidYAML",
			stateFixture:   "./fixture/loadWorkspace/invalid.yaml",
			specFixture:    "./fixture/loadWorkspace/empty.toml",
			expectedStatus: ErrorReadState,
			expectedState:  &pki.State{},
			expectedSpec:   &pki.Spec{},
		},
		{
			title:          "InvalidTOML",
			stateFixture:   "./fixture/loadWorkspace/empty.yaml",
			specFixture:    "./fixture/loadWorkspace/invalid.toml",
			expectedStatus: ErrorReadSpec,
			expectedState:  &pki.State{},
			expectedSpec:   &pki.Spec{},
		},
		{
			title:          "Simple",
			stateFixture:   "./fixture/loadWorkspace/simple.yaml",
			specFixture:    "./fixture/loadWorkspace/simple.toml",
			expectedStatus: 0,
			expectedState: &pki.State{
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
			},
			expectedSpec: &pki.Spec{
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
				RootPolicy: pki.Policy{
					Supplied: []string{"CommonName"},
				},
				IntermPolicy: pki.Policy{
					Supplied: []string{"CommonName"},
				},
			},
		},
		{
			title:          "Complex",
			stateFixture:   "./fixture/loadWorkspace/complex.yaml",
			specFixture:    "./fixture/loadWorkspace/complex.toml",
			expectedStatus: 0,
			expectedState: &pki.State{
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
			expectedSpec: &pki.Spec{
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
					DNSName:      []string{"example.com"},
					IPAddress:    []net.IP{net.ParseIP("127.0.0.1")},
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
				Metadata: pki.Metadata{
					"RootSkip":   []string{"IPAddress", "StreetAddress", "PostalCode"},
					"IntermSkip": []string{"IPAddress", "StreetAddress", "PostalCode"},
					"ServerSkip": []string{"StreetAddress", "PostalCode"},
					"ClientSkip": []string{"StreetAddress", "PostalCode"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			stateYAML, err := os.ReadFile(test.stateFixture)
			assert.NoError(t, err)
			err = os.WriteFile(pki.FileState, stateYAML, 0644)
			assert.NoError(t, err)

			specTOML, err := os.ReadFile(test.specFixture)
			assert.NoError(t, err)
			err = os.WriteFile(pki.FileSpec, specTOML, 0644)
			assert.NoError(t, err)

			mockUI := newMockUI(nil)
			state, spec, status := loadWorkspace(mockUI)

			assert.Equal(t, test.expectedStatus, status)
			if test.expectedStatus == 0 {
				assert.Equal(t, test.expectedState, state)
				assert.Equal(t, test.expectedSpec, spec)
			}

			err = util.DeleteAll("", pki.FileState, pki.FileSpec)
			assert.NoError(t, err)
		})
	}
}

func TestSaveWorkspace(t *testing.T) {
	tests := []struct {
		title                string
		state                *pki.State
		spec                 *pki.Spec
		expectedStatus       int
		expectedStateFixture string
		expectedSpecFixture  string
	}{
		{
			title:                "ZeroStateSpec",
			state:                &pki.State{},
			spec:                 &pki.Spec{},
			expectedStatus:       0,
			expectedStateFixture: "./fixture/saveWorkspace/zero.yaml",
			expectedSpecFixture:  "./fixture/saveWorkspace/zero.toml",
		},
		{
			title:                "DefaultStateSpec",
			state:                pki.NewState(),
			spec:                 pki.NewSpec(),
			expectedStatus:       0,
			expectedStateFixture: "./fixture/saveWorkspace/default.yaml",
			expectedSpecFixture:  "./fixture/saveWorkspace/default.toml",
		},
		{
			title: "CustomStateSpec",
			state: &pki.State{
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
			spec: &pki.Spec{
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
				Metadata: pki.Metadata{
					"RootSkip":   []string{"IPAddress", "StreetAddress", "PostalCode"},
					"IntermSkip": []string{"IPAddress", "StreetAddress", "PostalCode"},
					"ServerSkip": []string{"StreetAddress", "PostalCode"},
					"ClientSkip": []string{"StreetAddress", "PostalCode"},
				},
			},
			expectedStatus:       0,
			expectedStateFixture: "./fixture/saveWorkspace/custom.yaml",
			expectedSpecFixture:  "./fixture/saveWorkspace/custom.toml",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			mockUI := newMockUI(nil)
			status := saveWorkspace(test.state, test.spec, mockUI)

			assert.Equal(t, test.expectedStatus, status)
			if test.expectedStatus == 0 {
				stateYAML, err := os.ReadFile(pki.FileState)
				assert.NoError(t, err)

				expectedStateYAML, err := os.ReadFile(test.expectedStateFixture)
				assert.NoError(t, err)
				assert.Equal(t, string(expectedStateYAML), string(stateYAML))

				specTOML, err := os.ReadFile(pki.FileSpec)
				assert.NoError(t, err)

				expectedSpecTOML, err := os.ReadFile(test.expectedSpecFixture)
				assert.NoError(t, err)
				assert.Equal(t, string(expectedSpecTOML), string(specTOML))
			}

			err := util.DeleteAll("", pki.FileState, pki.FileSpec)
			assert.NoError(t, err)
		})
	}
}

func TestResolveByName(t *testing.T) {
	tests := []struct {
		title        string
		keyFile      string
		name         string
		expectedCert pki.Cert
	}{
		{
			"InvalidRootName",
			path.Join(pki.DirRoot, "top.ca.key"),
			"top",
			pki.Cert{},
		},
		{
			"ResolveRoot",
			path.Join(pki.DirRoot, "root.ca.key"),
			"root",
			pki.Cert{Name: "root", Type: pki.CertTypeRoot},
		},
		{
			"ResolveIntermediate",
			path.Join(pki.DirInterm, "interm.ca.key"),
			"interm",
			pki.Cert{Name: "interm", Type: pki.CertTypeInterm},
		},
		{
			"ResolveServer",
			path.Join(pki.DirServer, "server.key"),
			"server",
			pki.Cert{Name: "server", Type: pki.CertTypeServer},
		},
		{
			"ResolveClient",
			path.Join(pki.DirClient, "client.key"),
			"client",
			pki.Cert{Name: "client", Type: pki.CertTypeClient},
		},
	}

	err := pki.NewWorkspace(nil, nil)
	assert.NoError(t, err)
	defer pki.CleanupWorkspace() // nolint: errcheck

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			err := os.WriteFile(test.keyFile, []byte("mocked cert"), 0644)
			assert.NoError(t, err)

			c := resolveByName(test.name)
			assert.Equal(t, test.expectedCert, c)

			os.Remove(test.keyFile)
		})
	}
}

func TestAskForNewState(t *testing.T) {
	tests := []struct {
		title         string
		input         string
		expectError   bool
		expectedState *pki.State
	}{
		{
			"ErrorNoInputForRoot",
			``,
			true,
			nil,
		},
		{
			"ErrorNoInputForInterm",
			`10
			4096
			7300
			`,
			true,
			nil,
		},
		{
			"ErrorNoInputForServer",
			`10
			4096
			7300
			100
			4096
			3650
			`,
			true,
			nil,
		},
		{
			"ErrorNoInputForClient",
			`10
			4096
			7300
			100
			4096
			3650
			1000
			2048
			375
			`,
			true,
			nil,
		},
		{
			"SuccessEnterSome",
			`10
			4096
			7300
			100
			4096
			3650






			`,
			false,
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
			},
		},
		{
			"SuccessEnterAll",
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
			false,
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
					Length: 2048, Days: 375,
				},
				Client: pki.Config{
					Serial: 10000,
					Length: 2048,
					Days:   40,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			r := strings.NewReader(test.input)
			mockUI := newMockUI(r)
			state, err := askForNewState(mockUI)

			if test.expectError {
				assert.Error(t, err)
				assert.Nil(t, state)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedState, state)
			}
		})
	}
}

func TestAskForNewSpec(t *testing.T) {
	tests := []struct {
		title        string
		input        string
		expectError  bool
		expectedSpec *pki.Spec
	}{
		{
			"ErrorNoInputForCommon",
			"",
			true,
			nil,
		},
		{
			"ErrorNoInputForRoot",
			"\n\n\n\n\n\n\n\n\n\n",
			true,
			nil,
		},
		{
			"ErrorNoInputForInterm",
			"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n",
			true,
			nil,
		},
		{
			"ErrorNoInputForServer",
			"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n",
			true,
			nil,
		},
		{
			"ErrorNoInputForClient",
			"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n",
			true,
			nil,
		},
		{
			"ErrorNoInputForRootPolicy",
			"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n",
			true,
			nil,
		},
		{
			"ErrorNoInputForIntermPolicy",
			"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n",
			true,
			nil,
		},
		{
			"SuccessSimple",
			"CA\n\n\nMilad\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n" +
				"\n\n" +
				"\n\n",
			false,
			&pki.Spec{
				Root: pki.Claim{
					Country:      []string{"CA"},
					Organization: []string{"Milad"},
				},
				Interm: pki.Claim{
					Country:      []string{"CA"},
					Organization: []string{"Milad"},
				},
				Server: pki.Claim{
					Country:      []string{"CA"},
					Organization: []string{"Milad"},
				},
				Client: pki.Claim{
					Country:      []string{"CA"},
					Organization: []string{"Milad"},
				},
				RootPolicy: pki.Policy{
					Supplied: []string{"CommonName"},
				},
				IntermPolicy: pki.Policy{
					Supplied: []string{"CommonName"},
				},
				Metadata: pki.Metadata{},
			},
		},
		{
			"SuccessComplex",
			"CA\nOntario\n\nMilad\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n" +
				"Ottawa\nSRE\n\n\n\n\n\n" +
				"Toronto,Montreal\nR&D\nexample.com\n127.0.0.1\n\n\n\n" +
				"Ottawa\n\n\n\nmilad@example.com\n\n\n" +
				"Country,Organization\nCommonName\n" +
				"Organization\nCommonName\n",
			false,
			&pki.Spec{
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
					OrganizationalUnit: []string{"SRE"},
				},
				Server: pki.Claim{
					Country:            []string{"CA"},
					Province:           []string{"Ontario"},
					Locality:           []string{"Toronto", "Montreal"},
					Organization:       []string{"Milad"},
					OrganizationalUnit: []string{"R&D"},
					DNSName:            []string{"example.com"},
					IPAddress:          []net.IP{net.ParseIP("127.0.0.1")},
				},
				Client: pki.Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Locality:     []string{"Ottawa"},
					Organization: []string{"Milad"},
					EmailAddress: []string{"milad@example.com"},
				},
				RootPolicy: pki.Policy{
					Match:    []string{"Country", "Organization"},
					Supplied: []string{"CommonName"},
				},
				IntermPolicy: pki.Policy{
					Match:    []string{"Organization"},
					Supplied: []string{"CommonName"},
				},
				Metadata: pki.Metadata{},
			},
		},
		{
			"SuccessWithSkip",
			"CA\n\n\nMilad\n\n\n\n\n-\n-\n" +
				"\n\n\n-\n-\n\n" +
				"\n\nSRE\n-\n-\n\n" +
				"\nToronto,Montreal\nR&D\nexample.com\n127.0.0.1\n\n" +
				"\nOttawa\n\n\n\nmilad@example.com\n" +
				"Country,Organization\nCommonName\n" +
				"Organization\nCommonName\n",
			false,
			&pki.Spec{
				Root: pki.Claim{
					Country:      []string{"CA"},
					Organization: []string{"Milad"},
				},
				Interm: pki.Claim{
					Country:            []string{"CA"},
					Organization:       []string{"Milad"},
					OrganizationalUnit: []string{"SRE"},
				},
				Server: pki.Claim{
					Country:            []string{"CA"},
					Locality:           []string{"Toronto", "Montreal"},
					Organization:       []string{"Milad"},
					OrganizationalUnit: []string{"R&D"},
					DNSName:            []string{"example.com"},
					IPAddress:          []net.IP{net.ParseIP("127.0.0.1")},
				},
				Client: pki.Claim{
					Country:      []string{"CA"},
					Locality:     []string{"Ottawa"},
					Organization: []string{"Milad"},
					EmailAddress: []string{"milad@example.com"},
				},
				RootPolicy: pki.Policy{
					Match:    []string{"Country", "Organization"},
					Supplied: []string{"CommonName"},
				},
				IntermPolicy: pki.Policy{
					Match:    []string{"Organization"},
					Supplied: []string{"CommonName"},
				},
				Metadata: pki.Metadata{
					mdRootSkip:   []string{"Claim.StreetAddress", "Claim.PostalCode", "Claim.DNSName", "Claim.IPAddress"},
					mdIntermSkip: []string{"Claim.StreetAddress", "Claim.PostalCode", "Claim.DNSName", "Claim.IPAddress"},
					mdServerSkip: []string{"Claim.StreetAddress", "Claim.PostalCode"},
					mdClientSkip: []string{"Claim.StreetAddress", "Claim.PostalCode"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			r := strings.NewReader(test.input)
			mockUI := newMockUI(r)
			spec, err := askForNewSpec(mockUI)

			if test.expectError {
				assert.Error(t, err)
				assert.Nil(t, spec)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedSpec, spec)
			}
		})
	}
}

func TestAskForConfig(t *testing.T) {
	tests := []struct {
		title          string
		config         *pki.Config
		c              pki.Cert
		skipList       []string
		input          string
		expectError    bool
		expectedConfig *pki.Config
	}{
		{
			"ErrorNoInput",
			&pki.Config{},
			pki.Cert{},
			nil,
			``,
			true,
			nil,
		},
		{
			"SucessAskForRoot",
			&pki.Config{},
			pki.Cert{
				Type: pki.CertTypeRoot,
			},
			nil,
			`10
			4096
			7300
			secret
			secret
			`,
			false,
			&pki.Config{
				Serial:   10,
				Length:   4096,
				Days:     7300,
				Password: "secret",
			},
		},
		{
			"SuccessAskForInterm",
			&pki.Config{
				Password: "password",
			},
			pki.Cert{
				Type: pki.CertTypeInterm,
			},
			nil,
			`100
			4096
			3650
			`,
			false,
			&pki.Config{
				Serial:   100,
				Length:   4096,
				Days:     3650,
				Password: "password",
			},
		},
		{
			"SuccessAskForServer",
			&pki.Config{},
			pki.Cert{
				Type: pki.CertTypeServer,
			},
			nil,
			`1000
			2048
			375
			`,
			false,
			&pki.Config{
				Serial: 1000,
				Length: 2048,
				Days:   375,
			},
		},
		{
			"SuccessAskForClient",
			&pki.Config{},
			pki.Cert{
				Type: pki.CertTypeClient,
			},
			nil,
			`10000
			2048
			40
			`,
			false,
			&pki.Config{
				Serial: 10000,
				Length: 2048,
				Days:   40,
			},
		},
		{
			"SuccessWithSkip",
			&pki.Config{},
			pki.Cert{},
			[]string{"Config.Serial", "Config.Password"},
			`2048
			365
			`,
			false,
			&pki.Config{
				Length: 2048,
				Days:   365,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			r := strings.NewReader(test.input)
			mockUI := newMockUI(r)
			err := askForConfig(test.config, test.c, &test.skipList, mockUI)

			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedConfig, test.config)
			}
		})
	}
}

func TestAskForClaim(t *testing.T) {
	tests := []struct {
		title         string
		claim         *pki.Claim
		c             pki.Cert
		skipList      []string
		input         string
		expectError   bool
		expectedClaim *pki.Claim
	}{
		{
			"ErrorNoInput",
			&pki.Claim{},
			pki.Cert{},
			nil,
			"",
			true,
			nil,
		},
		{
			"SuccessAskForRoot",
			&pki.Claim{
				Country:      []string{"CA"},
				Organization: []string{"Milad"},
			},
			pki.Cert{
				Type: pki.CertTypeRoot,
			},
			[]string{"Claim.StreetAddress", "Claim.PostalCode"},
			"RootCA\n\n\n\n\n\n\n",
			false,
			&pki.Claim{
				CommonName:   "RootCA",
				Country:      []string{"CA"},
				Organization: []string{"Milad"},
			},
		},
		{
			"SuccessAskForInterm",
			&pki.Claim{
				Country:      []string{"CA"},
				Organization: []string{"Milad"},
			},
			pki.Cert{
				Type: pki.CertTypeInterm,
			},
			[]string{"Claim.StreetAddress", "Claim.PostalCode"},
			"IntermediateCA\n\n\nSRE\n\n\n\n",
			false,
			&pki.Claim{
				CommonName:         "IntermediateCA",
				Country:            []string{"CA"},
				Organization:       []string{"Milad"},
				OrganizationalUnit: []string{"SRE"},
			},
		},
		{
			"SuccessAskForServer",
			&pki.Claim{
				Country:      []string{"CA"},
				Organization: []string{"Milad"},
			},
			pki.Cert{
				Type: pki.CertTypeServer,
			},
			[]string{"Claim.StreetAddress", "Claim.PostalCode"},
			"Server\n\n\nR&D\nexample.com\n127.0.0.1,8.8.8.8\n\n",
			false,
			&pki.Claim{
				CommonName:         "Server",
				Country:            []string{"CA"},
				Organization:       []string{"Milad"},
				OrganizationalUnit: []string{"R&D"},
				DNSName:            []string{"example.com"},
				IPAddress:          []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("8.8.8.8")},
			},
		},
		{
			"SuccessAskForClient",
			&pki.Claim{
				Country:      []string{"CA"},
				Organization: []string{"Milad"},
			},
			pki.Cert{
				Type: pki.CertTypeClient,
			},
			[]string{"Claim.StreetAddress", "Claim.PostalCode"},
			"Client\nOntario\nOttawa\nQE\n\n\n\n",
			false,
			&pki.Claim{
				CommonName:         "Client",
				Country:            []string{"CA"},
				Province:           []string{"Ontario"},
				Locality:           []string{"Ottawa"},
				Organization:       []string{"Milad"},
				OrganizationalUnit: []string{"QE"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			r := strings.NewReader(test.input)
			mockUI := newMockUI(r)
			err := askForClaim(test.claim, test.c, &test.skipList, mockUI)

			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedClaim, test.claim)
			}
		})
	}
}
