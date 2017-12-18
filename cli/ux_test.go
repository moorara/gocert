package cli

import (
	"net"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
)

type inner struct {
	Int    int
	String string
}

type example struct {
	unexported   int
	Bool         bool    `custom:"-"`
	Int          int     `custom:"int"`
	Int64        int64   `custom:"int64"`
	Float32      float32 `custom:"float32"`
	Float64      float64 `custom:"float64"`
	String       string  `custom:"-" secret:"required,6"`
	Text         string  `custom:"text,omitempty" secret:"optional"`
	IntSlice     []int
	Int64Slice   []int64
	Float32Slice []float32
	Float64Slice []float64
	StringSlice  []string
	IPSlice      []net.IP
	Inner        inner
}

func TestAskForStruct(t *testing.T) {
	tests := []struct {
		title           string
		tagKey          string
		ignoreOmitted   bool
		example         example
		input           string
		expectedExample example
	}{
		{
			"EmptyNoInput",
			"custom", false,
			example{},
			``,
			example{},
		},
		{
			"SimpleOmitEmpty",
			"custom", true,
			example{},
			`27
			64
			2.71
			3.1415
			`,
			example{
				Int:     27,
				Int64:   64,
				Float32: 2.71,
				Float64: 3.1415,
			},
		},
		{
			"EnterInvalidSecret",
			"custom", false,
			example{
				Bool:    true,
				Int:     27,
				Int64:   64,
				Float32: 2.71,
				Float64: 3.1415,
			},
			`five5
			`,
			example{
				Bool:    true,
				Int:     27,
				Int64:   64,
				Float32: 2.71,
				Float64: 3.1415,
			},
		},
		{
			"EnterNotConfirmedSecret",
			"custom", false,
			example{
				Bool:    true,
				Int:     27,
				Int64:   64,
				Float32: 2.71,
				Float64: 3.1415,
			},
			`secret
			`,
			example{
				Bool:    true,
				Int:     27,
				Int64:   64,
				Float32: 2.71,
				Float64: 3.1415,
			},
		},
		{
			"EnterNotMatchingSecret",
			"custom", false,
			example{
				Bool:    true,
				Int:     27,
				Int64:   64,
				Float32: 2.71,
				Float64: 3.1415,
			},
			`secret
			notMatched
			`,
			example{
				Bool:    true,
				Int:     27,
				Int64:   64,
				Float32: 2.71,
				Float64: 3.1415,
			},
		},
		{
			"EnterAll",
			"custom", false,
			example{},
			`true
			27
			64
			2.71
			3.1415
			secret
			secret
			password
			password
			1,2,3,4
			1000000,1000000000
			2.71,3.14
			2.7182818284,3.1415926535
			Milad,Mona
			8.8.8.8,127.0.0.1
			1001
			nested
			`,
			example{
				Bool:         true,
				Int:          27,
				Int64:        64,
				Float32:      2.71,
				Float64:      3.1415,
				String:       "secret",
				Text:         "password",
				IntSlice:     []int{1, 2, 3, 4},
				Int64Slice:   []int64{1000000, 1000000000},
				Float32Slice: []float32{2.71, 3.14},
				Float64Slice: []float64{2.7182818284, 3.1415926535},
				StringSlice:  []string{"Milad", "Mona"},
				IPSlice:      []net.IP{net.ParseIP("8.8.8.8"), net.ParseIP("127.0.0.1")},
				Inner: inner{
					Int:    1001,
					String: "nested",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			mockUI := cli.NewMockUi()
			mockUI.InputReader = strings.NewReader(test.input)
			askForStruct(&test.example, test.tagKey, test.ignoreOmitted, mockUI)

			assert.Equal(t, test.expectedExample, test.example)
		})
	}
}
