package ezFile

import (
	"errors"
	"os"
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
func CreateFile(path string, fileName string, overwrite bool, openFlag int) (*os.File, error) {
	if !IsDirExists(path) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return nil, err
		}
	}
	if IsFileExists(path + fileName) {
		if !overwrite {
			return nil, errors.New("file already exist")
		}
	}
	return os.OpenFile(path+fileName, openFlag, os.ModePerm)
}
