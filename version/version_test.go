package version

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSemVer(t *testing.T) {
	tests := []struct {
		version, revision, branch, buildTime string
		expectedSemVer                       string
	}{
		{
			"0.1.0-0", "401f690", "local", "2017-11-28T10:00:00Z+0000",
			"0.1.0-0+401f690",
		},
		{
			"0.2.0-27", "365f39f", "ci", "2017-11-29T16:06:58Z+0000",
			"0.2.0-27+365f39f",
		},
		{
			"0.3.0", "b435957", "master", "2017-11-30T12:00:00Z+0000",
			"0.3.0+b435957",
		},
		{
			"1.0.0", "c9b0448", "master", "2017-12-01T09:09:00Z+0000",
			"1.0.0+c9b0448",
		},
	}

	for _, test := range tests {
		Version = test.version
		Revision = test.revision
		Branch = test.branch
		BuildTime = test.buildTime
		semVer := GetSemVer()

		assert.Equal(t, test.expectedSemVer, semVer)
	}
}

func TestGetFullSpec(t *testing.T) {
	tests := []struct {
		version, revision, branch, buildTime string
		expectedFullSpecPrefix               string
	}{
		{
			"0.1.0-0", "401f690", "local", "2017-11-28T10:00:00Z+0000",
			"0.1.0-0+401f690 local 2017-11-28T10:00:00Z+0000 go",
		},
		{
			"0.2.0-27", "365f39f", "ci", "2017-11-29T16:06:58Z+0000",
			"0.2.0-27+365f39f ci 2017-11-29T16:06:58Z+0000 go",
		},
		{
			"0.3.0", "b435957", "master", "2017-11-30T12:00:00Z+0000",
			"0.3.0+b435957 master 2017-11-30T12:00:00Z+0000 go",
		},
		{
			"1.0.0", "c9b0448", "master", "2017-12-01T09:09:00Z+0000",
			"1.0.0+c9b0448 master 2017-12-01T09:09:00Z+0000 go",
		},
	}

	for _, test := range tests {
		Version = test.version
		Revision = test.revision
		Branch = test.branch
		BuildTime = test.buildTime
		fullSpec := GetFullSpec()

		assert.True(t, strings.HasPrefix(fullSpec, test.expectedFullSpecPrefix))
	}
}
