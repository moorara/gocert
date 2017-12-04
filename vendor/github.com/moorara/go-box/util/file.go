package util

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// AbsPath returns absolute path from current execution path
func AbsPath(fromExec bool, elem ...string) string {
	var basePath string

	if fromExec {
		execPath, _ := filepath.Abs(os.Args[0])
		basePath = filepath.Dir(execPath)
	} else {
		basePath, _ = os.Getwd()
	}

	relPath := filepath.Join(elem...)
	absPath := filepath.Join(basePath, relPath)

	return absPath
}

// MkDirs creates directories
func MkDirs(basePath string, dirs ...string) (func(), error) {
	items := make([]string, 0)
	deleteFunc := func() {
		for _, item := range items {
			os.RemoveAll(item)
		}
	}

	for _, dir := range dirs {
		absPath := path.Join(basePath, dir)
		err := os.MkdirAll(absPath, 0755)
		if err != nil {
			return deleteFunc, err
		}
		items = append(items, dir)
	}

	return deleteFunc, nil
}

// WriteTempFile writes a file in your os temp directory
func WriteTempFile(content string) (string, func(), error) {
	file, err := ioutil.TempFile(os.TempDir(), "gobox-")
	if err != nil {
		return "", nil, err
	}

	if len(content) > 0 {
		err = ioutil.WriteFile(file.Name(), []byte(content), 0644)
		if err != nil {
			return "", nil, err
		}
	}

	err = file.Close()
	if err != nil {
		return "", nil, err
	}

	filePath := file.Name()
	deleteFunc := func() {
		os.Remove(filePath)
	}

	return filePath, deleteFunc, nil
}

// DeleteAll deletes all files and directories
func DeleteAll(basePath string, items ...string) error {
	for _, item := range items {
		absPath := path.Join(basePath, item)
		err := os.RemoveAll(absPath)
		if err != nil {
			return err
		}
	}

	return nil
}
