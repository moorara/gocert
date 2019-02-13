package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	tests := []struct {
		expectedValues []string
	}{
		{
			[]string{"version:", "revision:", "branch:", "goVersion:", "buildTool:", "buildTime:"},
		},
	}

	for _, tc := range tests {
		str := String()
		for _, expected := range tc.expectedValues {
			assert.Contains(t, str, expected)
		}
	}
}
