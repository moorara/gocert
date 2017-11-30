package cmd

import (
	"testing"

	"github.com/moorara/go-box/test"
	"github.com/stretchr/testify/assert"
)

func TestCert(t *testing.T) {
	tests := []struct {
		args             []string
		expectedExit     int
		expectedHelp     string
		expectedSynopsis string
	}{
		{
			[]string{},
			0,
			"cert command help",
			"cert command short help",
		},
		{
			[]string{"-root"},
			0,
			"cert command help",
			"cert command short help",
		},
	}

	// Null out standard io streams temporarily
	_, _, _, _, _, _, restore, err := test.PipeStdAll()
	defer restore()
	assert.NoError(t, err)

	for _, tc := range tests {
		cmd := NewCert()
		exit := cmd.Run(tc.args)
		help := cmd.Help()
		synopsis := cmd.Synopsis()

		assert.Equal(t, tc.expectedExit, exit)
		assert.Equal(t, tc.expectedHelp, help)
		assert.Equal(t, tc.expectedSynopsis, synopsis)
	}
}
