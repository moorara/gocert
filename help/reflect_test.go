package help

import (
	"net"
	"strings"
	"testing"

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
		example         *example
		input           string
		expectError     bool
		expectedExample *example
	}{
		{
			"ErrorNoBoolInput",
			"custom", false,
			&example{},
			``,
			true,
			nil,
		},
		{
			"ErrorNoIntInput",
			"custom", false,
			&example{},
			`true
			`,
			true,
			nil,
		},
		{
			"ErrorNoInt64Input",
			"custom", false,
			&example{},
			`true
			27
			`,
			true,
			nil,
		},
		{
			"ErrorNoFloat32Input",
			"custom", false,
			&example{},
			`true
			27
			64
			`,
			true,
			nil,
		},
		{
			"ErrorNoFloat64Input",
			"custom", false,
			&example{},
			`true
			27
			64
			2.71
			`,
			true,
			nil,
		},
		{
			"ErrorNoStringInput",
			"custom", false,
			&example{},
			`true
			27
			64
			2.71
			3.1415
			`,
			true,
			nil,
		},
		{
			"ErrorSecretInvalid",
			"custom", false,
			&example{
				Bool:    true,
				Int:     27,
				Int64:   64,
				Float32: 2.71,
				Float64: 3.1415,
			},
			`five5
			`,
			true,
			nil,
		},
		{
			"ErrorSecretNotConfirmed",
			"custom", false,
			&example{
				Bool:    true,
				Int:     27,
				Int64:   64,
				Float32: 2.71,
				Float64: 3.1415,
			},
			`secret
			`,
			true,
			nil,
		},
		{
			"ErrorSecretNotMatching",
			"custom", false,
			&example{
				Bool:    true,
				Int:     27,
				Int64:   64,
				Float32: 2.71,
				Float64: 3.1415,
			},
			`secret
			notMatched
			`,
			true,
			nil,
		},
		{
			"ErrorNoIntListInput",
			"custom", false,
			&example{
				Bool:    true,
				Int:     27,
				Int64:   64,
				Float32: 2.71,
				Float64: 3.1415,
				String:  "secret",
				Text:    "password",
			},
			``,
			true,
			nil,
		},
		{
			"ErrorNoInt64ListInput",
			"custom", false,
			&example{
				Bool:    true,
				Int:     27,
				Int64:   64,
				Float32: 2.71,
				Float64: 3.1415,
				String:  "secret",
				Text:    "password",
			},
			`1,2,3,4
			`,
			true,
			nil,
		},
		{
			"ErrorNoFloat32ListInput",
			"custom", false,
			&example{
				Bool:    true,
				Int:     27,
				Int64:   64,
				Float32: 2.71,
				Float64: 3.1415,
				String:  "secret",
				Text:    "password",
			},
			`1,2,3,4
			1000000,1000000000
			`,
			true,
			nil,
		},
		{
			"ErrorNoFloat64ListInput",
			"custom", false,
			&example{
				Bool:    true,
				Int:     27,
				Int64:   64,
				Float32: 2.71,
				Float64: 3.1415,
				String:  "secret",
				Text:    "password",
			},
			`1,2,3,4
			1000000,1000000000
			2.71,3.14
			`,
			true,
			nil,
		},
		{
			"ErrorNoStringListInput",
			"custom", false,
			&example{
				Bool:    true,
				Int:     27,
				Int64:   64,
				Float32: 2.71,
				Float64: 3.1415,
				String:  "secret",
				Text:    "password",
			},
			`1,2,3,4
			1000000,1000000000
			2.71,3.14
			2.7182818284,3.1415926535
			`,
			true,
			nil,
		},
		{
			"ErrorNoIPListInput",
			"custom", false,
			&example{
				Bool:    true,
				Int:     27,
				Int64:   64,
				Float32: 2.71,
				Float64: 3.1415,
				String:  "secret",
				Text:    "password",
			},
			`1,2,3,4
			1000000,1000000000
			2.71,3.14
			2.7182818284,3.1415926535
			Milad,Mona
			`,
			true,
			nil,
		},
		{
			"ErrorNoStructInput",
			"custom", false,
			&example{
				Bool:    true,
				Int:     27,
				Int64:   64,
				Float32: 2.71,
				Float64: 3.1415,
				String:  "secret",
				Text:    "password",
			},
			`1,2,3,4
			1000000,1000000000
			2.71,3.14
			2.7182818284,3.1415926535
			Milad,Mona
			8.8.8.8,127.0.0.1
			`,
			true,
			nil,
		},
		{
			"SuccessIgnoreOmitted",
			"custom", true,
			&example{},
			`27
			64
			2.71
			3.1415
			password
			password








			`,
			false,
			&example{
				Bool:    false,
				Int:     27,
				Int64:   64,
				Float32: 2.71,
				Float64: 3.1415,
				Text:    "password",
			},
		},
		{
			"SuccessEnterAll",
			"custom", false,
			&example{},
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
			false,
			&example{
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
			r := strings.NewReader(test.input)
			mockUI := NewMockUI(r)
			err := AskForStruct(test.example, test.tagKey, test.ignoreOmitted, mockUI)

			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedExample, test.example)
			}
		})
	}
}
