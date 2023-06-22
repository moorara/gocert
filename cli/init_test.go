package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/moorara/gocert/pki"
	"github.com/stretchr/testify/assert"
)

func TestNewInitCommand(t *testing.T) {
	tests := []struct {
		expectedSynopsis string
	}{
		{"Initializes a new workspace with desired configs and specs."},
		{"Initializes a new workspace with desired configs and specs."},
	}

	for _, test := range tests {
		cmd := NewInitCommand()

		assert.Equal(t, newColoredUI(), cmd.ui)
		assert.Equal(t, test.expectedSynopsis, cmd.Synopsis())
		assert.NotEmpty(t, cmd.Help())
	}
}

func TestInitCommand(t *testing.T) {
	tests := []struct {
		title                string
		args                 []string
		input                string
		expectedStateFixture string
		expectedSpecFixture  string
	}{
		{
			title: "DefaultStateSpec",
			args:  []string{},
			input: "\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n\n\n\n\n\n\n\n\n" +
				"\n\n" +
				"\n\n",
			expectedStateFixture: "./fixture/InitCommand/default.yaml",
			expectedSpecFixture:  "./fixture/InitCommand/default.toml",
		},
		{
			title: "CustomStateSpec",
			args:  []string{},
			input: "CA\n\n\nMilad\n\n\n\n\n-\n-\n" +
				"\n\n\n-\n-\n-\n" +
				"\n\nSRE\n-\n-\n-\n" +
				"Ontario\nOttawa\nR&D\nexample.org\n127.0.0.1\n\n" +
				"Ontario\nOttawa\nQE\n\n\n\n" +
				"Organization\nCommonName,OrganizationalUnit\n" +
				"Organization\nCommonName\n",
			expectedStateFixture: "./fixture/InitCommand/custom.yaml",
			expectedSpecFixture:  "./fixture/InitCommand/custom.toml",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			r := strings.NewReader(test.input)
			mockUI := newMockUI(r)
			cmd := &InitCommand{
				ui: mockUI,
			}

			exit := cmd.Run(test.args)
			assert.Zero(t, exit)

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

			err = pki.CleanupWorkspace()
			assert.NoError(t, err)
		})
	}
}

func TestInitCommandError(t *testing.T) {
	tests := []struct {
		title        string
		args         []string
		input        string
		expectedExit int
	}{
		{
			"InvalidFlag",
			[]string{"-invalid"},
			``,
			ErrorInvalidFlag,
		},
		{
			"NoInputForSpec",
			[]string{""},
			``,
			ErrorEnterSpec,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			r := strings.NewReader(test.input)
			mockUI := newMockUI(r)
			cmd := &InitCommand{
				ui: mockUI,
			}

			exit := cmd.Run(test.args)
			assert.Equal(t, test.expectedExit, exit)
		})
	}
}
