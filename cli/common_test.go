package cli

import (
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/moorara/go-box/util"
	"github.com/moorara/gocert/help"
	"github.com/moorara/gocert/pki"
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
		stateYAML      string
		specTOML       string
		expectedStatus int
		expectedState  *pki.State
		expectedSpec   *pki.Spec
	}{
		{
			"Empty",
			``,
			``,
			0,
			&pki.State{},
			&pki.Spec{},
		},
		{
			"InvalidYAML",
			`invalid yaml`,
			``,
			ErrorReadState,
			&pki.State{},
			&pki.Spec{},
		},
		{
			"InvalidTOML",
			``,
			`invalid toml`,
			ErrorReadSpec,
			&pki.State{},
			&pki.Spec{},
		},
		{
			"Simple",
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
			[root_policy]
				supplied = ["CommonName"]
			[intermediate_policy]
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
				RootPolicy: pki.Policy{
					Supplied: []string{"CommonName"},
				},
				IntermPolicy: pki.Policy{
					Supplied: []string{"CommonName"},
				},
			},
		},
		{
			"Complex",
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
				dns_name = [ "example.com" ]
				ip_address = [ "127.0.0.1" ]
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
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			stateYAML := strings.Replace(test.stateYAML, "\t", "  ", -1)
			err := ioutil.WriteFile(pki.FileState, []byte(stateYAML), 0644)
			assert.NoError(t, err)
			err = ioutil.WriteFile(pki.FileSpec, []byte(test.specTOML), 0644)
			assert.NoError(t, err)

			mockUI := help.NewMockUI(nil)
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
		title             string
		state             *pki.State
		spec              *pki.Spec
		expectedStatus    int
		expectedStateYAML string
		expectedSpecTOML  string
	}{
		{
			"Simple",
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
			"Complex",
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
		t.Run(test.title, func(t *testing.T) {
			mockUI := help.NewMockUI(nil)
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
			"Empty",
			"",
			"",
			pki.Cert{},
		},
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
	defer pki.CleanupWorkspace()
	assert.NoError(t, err)

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			ioutil.WriteFile(test.keyFile, []byte("mocked cert"), 0644)

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
			mockUI := help.NewMockUI(r)
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
			"ErrorNoInput",
			"",
			true,
			nil,
		},
		{
			"SuccessEnterCommon",
			"CA\n\n\nMilad\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n",
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
			},
		},
		{
			"SuccessEnterMore",
			"CA\nOntario\n\nMilad\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n" +
				"Ottawa\nR&D\n\n\n\n\n\n" +
				"Toronto,Montreal\n\nexample.com\n127.0.0.1\n\n\n\n" +
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
					OrganizationalUnit: []string{"R&D"},
				},
				Server: pki.Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Locality:     []string{"Toronto", "Montreal"},
					Organization: []string{"Milad"},
					DNSName:      []string{"example.com"},
					IPAddress:    []net.IP{net.ParseIP("127.0.0.1")},
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
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			r := strings.NewReader(test.input)
			mockUI := help.NewMockUI(r)
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
		input          string
		expectError    bool
		expectedConfig *pki.Config
	}{
		{
			"ErrorNoInput",
			&pki.Config{},
			pki.Cert{},
			``,
			true,
			nil,
		},
		{
			"SucessWithPassword",
			&pki.Config{},
			pki.Cert{
				Type: pki.CertTypeRoot,
			},
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
			"SuccessWithoutPassword",
			&pki.Config{
				Password: "dummy",
			},
			pki.Cert{
				Type: pki.CertTypeInterm,
			},
			`100
			4096
			3650
			`,
			false,
			&pki.Config{
				Serial:   100,
				Length:   4096,
				Days:     3650,
				Password: "dummy",
			},
		},
		{
			"SuccessWithoutPassword",
			&pki.Config{},
			pki.Cert{
				Type: pki.CertTypeServer,
			},
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
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			r := strings.NewReader(test.input)
			mockUI := help.NewMockUI(r)
			err := askForConfig(test.config, test.c, mockUI)

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
		input         string
		expectError   bool
		expectedClaim *pki.Claim
	}{
		{
			"ErrorNoInput",
			&pki.Claim{},
			pki.Cert{},
			"",
			true,
			nil,
		},
		{
			"SuccessSimple",
			&pki.Claim{},
			pki.Cert{
				Type: pki.CertTypeRoot,
			},
			"RootCA\nCA\n\n\nMilad\n\n\n\n\n\n\n",
			false,
			&pki.Claim{
				CommonName:   "RootCA",
				Country:      []string{"CA"},
				Organization: []string{"Milad"},
			},
		},
		{
			"SuccessComplex",
			&pki.Claim{
				Country:  []string{"CA"},
				Province: []string{"Ontario"},
				Locality: []string{"Ottawa"},
			},
			pki.Cert{
				Type: pki.CertTypeInterm,
			},
			"IntermediateCA\nMilad\nSRE\nexample.com\n8.8.8.8,127.0.0.1\n\n\n\n",
			false,
			&pki.Claim{
				CommonName:         "IntermediateCA",
				Country:            []string{"CA"},
				Province:           []string{"Ontario"},
				Locality:           []string{"Ottawa"},
				Organization:       []string{"Milad"},
				OrganizationalUnit: []string{"SRE"},
				DNSName:            []string{"example.com"},
				IPAddress:          []net.IP{net.ParseIP("8.8.8.8"), net.ParseIP("127.0.0.1")},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			r := strings.NewReader(test.input)
			mockUI := help.NewMockUI(r)
			err := askForClaim(test.claim, test.c, mockUI)

			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedClaim, test.claim)
			}
		})
	}
}
