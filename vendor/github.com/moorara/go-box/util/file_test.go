package util

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathFromCurrent(t *testing.T) {
	tests := []struct {
		execPath        string
		fromExec        bool
		subDirs         []string
		fileName        string
		expectedAbsPath string
	}{
		{
			"/usr/local/bin/godo",
			true,
			[]string{"more"},
			"cert.pem",
			"/usr/local/bin/more/cert.pem",
		},
		{
			"/home/milad/go/bin/gotest",
			true,
			[]string{".src", ".bin"},
			"gobench",
			"/home/milad/go/bin/.src/.bin/gobench",
		},
		{
			"",
			false,
			[]string{"test"},
			"file_test.go",
			path.Join(os.Getenv("GOPATH"), "src/github.com/moorara/go-box/util", "test/file_test.go"),
		},
	}

	for _, test := range tests {
		restore := ReplaceOSArgs([]string{test.execPath})
		elem := append(test.subDirs, test.fileName)
		absPath := AbsPath(test.fromExec, elem...)
		assert.Equal(t, test.expectedAbsPath, absPath)

		restore()
	}
}

func TestMkDirs(t *testing.T) {
	tests := []struct {
		basePath    string
		dirs        []string
		expectError bool
	}{
		{
			"",
			[]string{},
			false,
		},
		{
			"",
			[]string{""},
			true,
		},
		{
			AbsPath(true),
			[]string{""},
			false,
		},
		{
			"",
			[]string{" "},
			false,
		},
		{
			"",
			[]string{"code"},
			false,
		},
		{
			"",
			[]string{"bin", "src"},
			false,
		},
		{
			"",
			[]string{"bin", "src", "test", "release"},
			false,
		},
		{
			"",
			[]string{"src/server", "src/client", "src/mocks"},
			false,
		},
	}

	for _, test := range tests {
		delete, err := MkDirs(test.basePath, test.dirs...)
		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

		delete()
	}
}

func TestWriteTempFile(t *testing.T) {
	tests := []struct {
		content string
	}{
		{
			"",
		},
		{
			"Example file content",
		},
		{
			`
      [settings]
        user = "milad"
        token = "api_token"
      `,
		},
	}

	for _, test := range tests {
		path, delete, err := WriteTempFile(test.content)
		defer delete()
		assert.NoError(t, err)

		content, err := ioutil.ReadFile(path)
		assert.NoError(t, err)
		assert.Equal(t, test.content, string(content))
	}
}

func TestDeleteAll(t *testing.T) {
	tests := []struct {
		basePath    string
		dirs        []string
		files       []string
		expectError bool
	}{
		{
			"",
			[]string{},
			[]string{},
			false,
		},
		{
			"",
			[]string{"src"},
			[]string{},
			false,
		},
		{
			"",
			[]string{},
			[]string{"test.txt"},
			false,
		},
		{
			"",
			[]string{"bin", "src"},
			[]string{"test.txt"},
			false,
		},
		{
			"",
			[]string{"bin", "src", "src/server", "src/client"},
			[]string{"README.md", "src/server/index.js", "src/client/index.js"},
			false,
		},
	}

	for _, test := range tests {
		_, err := MkDirs(test.basePath, test.dirs...)
		assert.NoError(t, err)

		for _, file := range test.files {
			err = ioutil.WriteFile(file, []byte(""), 0644)
			assert.NoError(t, err)
		}

		paths := append(test.files, test.dirs...)
		err = DeleteAll(test.basePath, paths...)

		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
