package ezFile

import (
	"errors"
	"os"
	"path/filepath"
)

func IsDirExists(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return fileInfo.IsDir()
}
func IsFileExists(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return !fileInfo.IsDir()
}
func CreateFile(dir string, fileName string, overwrite bool, openFlag int) (*os.File, error) {
	if !IsDirExists(dir) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}
	if IsFileExists(filepath.Join(dir, fileName)) {
		if !overwrite {
			return nil, errors.New("file already exist")
		}
	}
	return os.OpenFile(filepath.Join(dir, fileName), openFlag, os.ModePerm)
}
func CreateDir(path string) error {
	if !IsDirExists(path) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
