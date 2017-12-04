package common

import (
	"testing"

	"github.com/mitchellh/cli"
	"github.com/moorara/go-box/util"
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

		ui := NewColoredUI()

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
