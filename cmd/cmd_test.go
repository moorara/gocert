package cmd

import (
	"testing"

	"github.com/mitchellh/cli"
	"github.com/moorara/go-box/test"
	"github.com/stretchr/testify/assert"
)

func TestNewColoredUi(t *testing.T) {
	tests := []struct {
		in, out, er string
	}{
		{"in", "out", "err"},
		{"input", "output", "error"},
	}

	for _, tc := range tests {
		inR, inW, outR, outW, errR, errW, restore, err := test.PipeStdAll()
		defer restore()
		assert.NoError(t, err)

		ui := newColoredUI()

		assert.Equal(t, cli.UiColorCyan, ui.OutputColor)
		assert.Equal(t, cli.UiColorGreen, ui.InfoColor)
		assert.Equal(t, cli.UiColorRed, ui.ErrorColor)
		assert.Equal(t, cli.UiColorYellow, ui.WarnColor)

		in, err := test.WriteToStdinPipe(inR, inW, tc.in)
		assert.NoError(t, err)
		assert.Equal(t, tc.in, in)

		out, err := test.ReadFromStdoutPipe(outR, outW, tc.out)
		assert.NoError(t, err)
		assert.Equal(t, tc.out, out)

		er, err := test.ReadFromStderrPipe(errR, errW, tc.er)
		assert.NoError(t, err)
		assert.Equal(t, tc.er, er)
	}
}
