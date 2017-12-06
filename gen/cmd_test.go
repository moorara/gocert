package gen

import (
	"errors"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/moorara/go-box/util"
	"github.com/moorara/gocert/cert"
	"github.com/moorara/gocert/common"
	"github.com/moorara/gocert/config"
	"github.com/stretchr/testify/assert"
)

type mockedCertManager struct {
	GenRootCAError     error
	GenIntermCAError   error
	GenServerCertError error
	GenClientCertError error

	GenRootCACalled     bool
	GenIntermCACalled   bool
	GenServerCertCalled bool
	GenClientCertCalled bool
}

func (m *mockedCertManager) GenRootCA(config.SettingsCA, config.Claim) error {
	m.GenRootCACalled = true
	return m.GenRootCAError
}

func (m *mockedCertManager) GenIntermCA(config.SettingsCA, config.Claim) error {
	m.GenIntermCACalled = true
	return m.GenIntermCAError
}

func (m *mockedCertManager) GenServerCert(config.Settings, config.Claim) error {
	m.GenServerCertCalled = true
	return m.GenServerCertError
}

func (m *mockedCertManager) GenClientCert(config.Settings, config.Claim) error {
	m.GenClientCertCalled = true
	return m.GenClientCertError
}

func mockWorkspace(state *config.State, spec *config.Spec) (func() error, error) {
	items := make([]string, 0)
	deleteFunc := func() error {
		return util.DeleteAll("", items...)
	}

	// Mock sub-directories
	_, err := util.MkDirs("", config.DirNameRoot, config.DirNameInterm, config.DirNameServer, config.DirNameClient)
	items = append(items, config.DirNameRoot, config.DirNameInterm, config.DirNameServer, config.DirNameClient)
	if err != nil {
		return deleteFunc, err
	}

	// Mock state file
	if state != nil {
		err := config.SaveState(state, config.FileNameState)
		if err != nil {
			return deleteFunc, err
		}
		items = append(items, config.FileNameState)
	}

	// Mock spec file
	if spec != nil {
		err := config.SaveSpec(spec, config.FileNameSpec)
		if err != nil {
			return deleteFunc, err
		}
		items = append(items, config.FileNameSpec)
	}

	return deleteFunc, nil
}

func TestNewCommand(t *testing.T) {
	cmdGen := NewCommand()

	assert.Equal(t, common.NewColoredUI(), cmdGen.ui)
	assert.Equal(t, cert.NewX509Manager(), cmdGen.cert)
}

func TestGenError(t *testing.T) {
	tests := []struct {
		state        *config.State
		spec         *config.Spec
		args         []string
		expectedExit int
	}{
		{
			nil,
			nil,
			[]string{"-invalid"},
			ErrorInvalidFlags,
		},
		{
			nil,
			&config.Spec{},
			[]string{},
			ErrorReadState,
		},
		{
			config.NewState(),
			nil,
			[]string{},
			ErrorReadSpec,
		},
	}

	for _, test := range tests {
		cleanup, err := mockWorkspace(test.state, test.spec)
		assert.NoError(t, err)

		cmd := &Command{
			ui:   cli.NewMockUi(),
			cert: &mockedCertManager{},
		}
		exit := cmd.Run(test.args)

		assert.Equal(t, textSynopsis, cmd.Synopsis())
		assert.Equal(t, textHelp, cmd.Help())
		assert.Equal(t, test.expectedExit, exit)

		err = cleanup()
		assert.NoError(t, err)
	}
}

func TestGenRootCA(t *testing.T) {
	tests := []struct {
		state          *config.State
		spec           *config.Spec
		args           []string
		input          string
		GenRootCAError error
		expectedExit   int
	}{
		{
			&config.State{},
			&config.Spec{},
			[]string{"-root"},
			``,
			errors.New("error"),
			ErrorRootCA,
		},
		{
			config.NewState(),
			&config.Spec{},
			[]string{"-root"},
			``,
			nil,
			0,
		},
		{
			config.NewState(),
			&config.Spec{
				Root: config.Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Organization: []string{"Milad"},
				},
			},
			[]string{"-root"},
			`secret
			RootCA
			Ottawa,Toronto
			`,
			nil,
			0,
		},
	}

	for _, test := range tests {
		cleanup, err := mockWorkspace(test.state, test.spec)
		assert.NoError(t, err)

		mockUI := cli.NewMockUi()
		mockUI.InputReader = strings.NewReader(test.input)

		cmd := &Command{
			ui: mockUI,
			cert: &mockedCertManager{
				GenRootCAError: test.GenRootCAError,
			},
		}
		exit := cmd.Run(test.args)

		assert.Equal(t, textSynopsis, cmd.Synopsis())
		assert.Equal(t, textHelp, cmd.Help())
		assert.Equal(t, test.expectedExit, exit)

		err = cleanup()
		assert.NoError(t, err)
	}
}

func TestGenIntermCA(t *testing.T) {
	tests := []struct {
		state            *config.State
		spec             *config.Spec
		args             []string
		input            string
		GenIntermCAError error
		expectedExit     int
	}{
		{
			&config.State{},
			&config.Spec{},
			[]string{"-intermediate"},
			``,
			errors.New("error"),
			ErrorIntermCA,
		},
		{
			config.NewState(),
			&config.Spec{},
			[]string{"-intermediate"},
			``,
			nil,
			0,
		},
		{
			config.NewState(),
			&config.Spec{
				Interm: config.Claim{
					Country:      []string{"CA"},
					Province:     []string{"Ontario"},
					Organization: []string{"Milad"},
				},
			},
			[]string{"-intermediate"},
			`secret
			IntermediateCA
			Ottawa,Toronto
			R&D
			`,
			nil,
			0,
		},
	}

	for _, test := range tests {
		cleanup, err := mockWorkspace(test.state, test.spec)
		assert.NoError(t, err)

		mockUI := cli.NewMockUi()
		mockUI.InputReader = strings.NewReader(test.input)

		cmd := &Command{
			ui: mockUI,
			cert: &mockedCertManager{
				GenIntermCAError: test.GenIntermCAError,
			},
		}
		exit := cmd.Run(test.args)

		assert.Equal(t, textSynopsis, cmd.Synopsis())
		assert.Equal(t, textHelp, cmd.Help())
		assert.Equal(t, test.expectedExit, exit)

		err = cleanup()
		assert.NoError(t, err)
	}
}
