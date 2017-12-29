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
	Bool         bool    `custom:"-" default:"true"`
	Int          int     `custom:"int" default:"27"`
	Int64        int64   `custom:"int64" default:"64"`
	Float32      float32 `custom:"float32" default:"2.71"`
	Float64      float64 `custom:"float64" default:"3.1415"`
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
		title            string
		tagKey           string
		ignoreOmitted    bool
		skipList         []string
		example          *example
		input            string
		expectError      bool
		expectedExample  *example
		expectedSkipList []string
	}{
		{
			"ErrorNoBoolInput",
			"custom", false,
			nil,
			&example{},
			``,
			true,
			nil,
			nil,
		},
		{
			"ErrorNoIntInput",
			"custom", false,
			nil,
			&example{},
			`true
			`,
			true,
			nil,
			nil,
		},
		{
			"ErrorNoInt64Input",
			"custom", false,
			nil,
			&example{},
			`true
			27
			`,
			true,
			nil,
			nil,
		},
		{
			"ErrorNoFloat32Input",
			"custom", false,
			nil,
			&example{},
			`true
			27
			64
			`,
			true,
			nil,
			nil,
		},
		{
			"ErrorNoFloat64Input",
			"custom", false,
			nil,
			&example{},
			`true
			27
			64
			2.71
			`,
			true,
			nil,
			nil,
		},
		{
			"ErrorNoStringInput",
			"custom", false,
			nil,
			&example{},
			`true
			27
			64
			2.71
			3.1415
			`,
			true,
			nil,
			nil,
		},
		{
			"ErrorSecretInvalid",
			"custom", false,
			nil,
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
			nil,
		},
		{
			"ErrorSecretNotConfirmed",
			"custom", false,
			nil,
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
			nil,
		},
		{
			"ErrorSecretNotMatching",
			"custom", false,
			nil,
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
			nil,
		},
		{
			"ErrorNoIntListInput",
			"custom", false,
			nil,
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
			nil,
		},
		{
			"ErrorNoInt64ListInput",
			"custom", false,
			nil,
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
			nil,
		},
		{
			"ErrorNoFloat32ListInput",
			"custom", false,
			nil,
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
			nil,
		},
		{
			"ErrorNoFloat64ListInput",
			"custom", false,
			nil,
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
			nil,
		},
		{
			"ErrorNoStringListInput",
			"custom", false,
			nil,
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
			nil,
		},
		{
			"ErrorNoIPListInput",
			"custom", false,
			nil,
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
			nil,
		},
		{
			"ErrorNoStructInput",
			"custom", false,
			nil,
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
			nil,
		},
		{
			"SuccessIgnoreOmitted",
			"custom", true,
			[]string{},
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
			[]string{},
		},
		{
			"SuccessEnterAll",
			"custom", false,
			[]string{},
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
			[]string{},
		},
		{
			"SuccessUsingSkipFields",
			"custom", true,
			[]string{"example.String", "example.Text"},
			&example{},
			`27
			64
			-
			-
			1,2,3,4
			-
			-
			-
			-
			-
			-
			-
			`,
			false,
			&example{
				Int:      27,
				Int64:    64,
				IntSlice: []int{1, 2, 3, 4},
			},
			[]string{
				"example.String", "example.Text",
				"example.Float32", "example.Float64",
				"example.Int64Slice", "example.Float32Slice", "example.Float64Slice", "example.StringSlice", "example.IPSlice",
				"inner.Int", "inner.String"},
		},
		{
			"SuccessUsingDefaultsAndSkip",
			"custom", false,
			[]string{},
			&example{},
			`




			secret
			secret
			password
			password
			-
			-
			-
			-
			-
			-
			-
			-
			`,
			false,
			&example{
				Bool:    true,
				Int:     27,
				Int64:   64,
				Float32: 2.71,
				Float64: 3.1415,
				String:  "secret",
				Text:    "password",
			},
			[]string{
				"example.IntSlice", "example.Int64Slice", "example.Float32Slice", "example.Float64Slice", "example.StringSlice", "example.IPSlice",
				"inner.Int", "inner.String",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			r := strings.NewReader(test.input)
			mockUI := NewMockUI(r)
			err := AskForStruct(test.example, test.tagKey, test.ignoreOmitted, &test.skipList, mockUI)

			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedExample, test.example)
				assert.Equal(t, test.expectedSkipList, test.skipList)
			}
		})
	}
}
