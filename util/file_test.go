package util

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbsPath(t *testing.T) {
	basePath, err := os.Getwd()
	assert.NoError(t, err)

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
			path.Join(basePath, "test/file_test.go"),
		},
	}

	for _, tc := range tests {
		restore := ReplaceOSArgs([]string{tc.execPath})
		elem := append(tc.subDirs, tc.fileName)
		absPath := AbsPath(tc.fromExec, elem...)
		assert.Equal(t, tc.expectedAbsPath, absPath)

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

	for _, tc := range tests {
		delete, err := MkDirs(tc.basePath, tc.dirs...)
		if tc.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

		delete()
	}
}

func TestCreateTempFile(t *testing.T) {
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

	for _, tc := range tests {
		path, delete, err := CreateTempFile(tc.content)
		defer delete()
		assert.NoError(t, err)

		content, err := os.ReadFile(path)
		assert.NoError(t, err)
		assert.Equal(t, tc.content, string(content))
	}
}

func TestConcatFiles(t *testing.T) {
	tests := []struct {
		name            string
		dest            string
		destContent     string
		fileContents    map[string]string
		append          bool
		expectError     bool
		expectedContent string
	}{
		{
			name:            "",
			dest:            "list",
			destContent:     "tangerine",
			fileContents:    map[string]string{},
			append:          false,
			expectError:     false,
			expectedContent: "",
		},
		{
			name:            "",
			dest:            "list",
			destContent:     "tangerine",
			fileContents:    map[string]string{},
			append:          true,
			expectError:     false,
			expectedContent: "tangerine",
		},
		{
			name:        "",
			dest:        "list",
			destContent: "apple ",
			fileContents: map[string]string{
				"item1": "pear ",
				"item2": "orange ",
			},
			append:          false,
			expectError:     false,
			expectedContent: "pear orange ",
		},
		{
			name:        "",
			dest:        "list",
			destContent: "apple ",
			fileContents: map[string]string{
				"item1": "pear ",
				"item2": "orange ",
			},
			append:          true,
			expectError:     false,
			expectedContent: "apple pear orange ",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := os.WriteFile(tc.dest, []byte(tc.destContent), 0644)
			assert.NoError(t, err)

			files := make([]string, 0)
			for file, content := range tc.fileContents {
				err := os.WriteFile(file, []byte(content), 0644)
				assert.NoError(t, err)

				files = append(files, file)
			}

			err = ConcatFiles(tc.dest, tc.append, files...)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				content, err := os.ReadFile(tc.dest)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedContent, string(content))
			}

			// Cleanup temporary files
			files = append(files, tc.dest)
			err = DeleteAll("", files...)
			assert.NoError(t, err)
		})
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

	for _, tc := range tests {
		_, err := MkDirs(tc.basePath, tc.dirs...)
		assert.NoError(t, err)

		for _, file := range tc.files {
			err = os.WriteFile(file, []byte(""), 0644)
			assert.NoError(t, err)
		}

		paths := append(tc.files, tc.dirs...)
		err = DeleteAll(tc.basePath, paths...)

		if tc.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
