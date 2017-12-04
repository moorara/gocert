package config

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
)

type toy struct {
	ID        string   `custom:"-"`
	Name      string   `custom:"name"`
	Available bool     `custom:"available"`
	Count     int      `custom:"count"`
	Serial    int64    `custom:"serial"`
	Parts     []string `custom:"parts"`
}

func TestFillIn(t *testing.T) {
	tests := []struct {
		tagKey         string
		includeOmitted bool
		toy            toy
		input          string
		expectedToy    toy
	}{
		{
			"custom", false,
			toy{},
			`Fidget
			false
			0
			1111
			cube,spinner
			`,
			toy{
				Name:      "Fidget",
				Available: false,
				Count:     0,
				Serial:    1111,
				Parts:     []string{"cube", "spinner"},
			},
		},
		{
			"custom", true,
			toy{
				Name: "Robo",
			},
			`bbbb
			true
			2
			2222
			body,head,hands,legs
			`,
			toy{
				ID:        "bbbb",
				Name:      "Robo",
				Available: true,
				Count:     2,
				Serial:    2222,
				Parts:     []string{"body", "head", "hands", "legs"},
			},
		},
		{
			"custom", true,
			toy{
				Name:      "Laser",
				Available: true,
				Count:     5,
			},
			`dddd
      4444
      emitter,filter
			`,
			toy{
				ID:        "dddd",
				Name:      "Laser",
				Available: true,
				Count:     5,
				Serial:    4444,
				Parts:     []string{"emitter", "filter"},
			},
		},
		{
			"custom", false,
			toy{
				Name:      "Car",
				Available: true,
				Count:     10,
				Serial:    248,
				Parts:     []string{"wheel", "engine"},
			},
			``,
			toy{
				Name:      "Car",
				Available: true,
				Count:     10,
				Serial:    248,
				Parts:     []string{"wheel", "engine"},
			},
		},
		{
			"custom", true,
			toy{
				ID:        "abcdef",
				Name:      "Drone",
				Available: true,
				Count:     1,
				Serial:    123456789,
				Parts:     []string{"wing", "camera"},
			},
			``,
			toy{
				ID:        "abcdef",
				Name:      "Drone",
				Available: true,
				Count:     1,
				Serial:    123456789,
				Parts:     []string{"wing", "camera"},
			},
		},
	}

	for _, test := range tests {
		mockUI := cli.NewMockUi()
		mockUI.InputReader = strings.NewReader(test.input)
		fillIn(&test.toy, test.tagKey, test.includeOmitted, mockUI)

		assert.Equal(t, test.expectedToy, test.toy)
	}
}
