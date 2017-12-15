package util

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

const (
	defaultDirPerm  = 0755
	defaultFilePerm = 0644
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
		err := os.MkdirAll(absPath, defaultDirPerm)
		if err != nil {
			return deleteFunc, err
		}
		items = append(items, dir)
	}

	return deleteFunc, nil
}

// CreateTempFile writes a file in your os temp directory
func CreateTempFile(content string) (string, func(), error) {
	file, err := ioutil.TempFile(os.TempDir(), "gobox-")
	if err != nil {
		return "", nil, err
	}

	if len(content) > 0 {
		err = ioutil.WriteFile(file.Name(), []byte(content), defaultFilePerm)
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

// ConcatFiles concats a set files into a new or existing file
func ConcatFiles(dest string, append bool, files ...string) error {
	var flag int
	if append {
		flag = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	} else {
		flag = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}

	df, err := os.OpenFile(dest, flag, defaultFilePerm)
	if err != nil {
		return err
	}

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		_, err = io.Copy(df, f)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
	}

	err = df.Close()
	if err != nil {
		return err
	}

	return nil
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
