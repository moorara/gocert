package help

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockUI(t *testing.T) {
	tests := []struct {
		input          string
		expectedErrors []bool
		expectedValues []string
	}{
		{
			``,
			[]bool{true},
			[]string{""},
		},
		{
			`
			`,
			[]bool{false},
			[]string{""},
		},
		{
			`Token
			`,
			[]bool{false},
			[]string{"Token"},
		},
		{
			`API Token
			`,
			[]bool{false},
			[]string{"API Token"},
		},
		{
			`First
			`,
			[]bool{false, true},
			[]string{"First", ""},
		},
		{
			`First
			Second
			`,
			[]bool{false, false},
			[]string{"First", "Second"},
		},
	}

	for _, test := range tests {
		r := strings.NewReader(test.input)
		ui := NewMockUI(r)

		for i, expectedValue := range test.expectedValues {
			value, err := ui.Ask("Enter:")

			if test.expectedErrors[i] {
				assert.Error(t, err)
				assert.Empty(t, value)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedValue, value)
			}
		}
	}
}
