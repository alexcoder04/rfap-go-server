package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexcoder04/rfap-go-server/settings"
)

func RetError(metadata HeaderMetadata, errorCode int, errorMsg string) HeaderMetadata {
	metadata.ErrorCode = errorCode
	metadata.ErrorMessage = errorMsg
	return metadata
}

func validatePath(path string) (string, error) {
	path, err := filepath.Abs(filepath.Join(settings.Config.PublicFolder, path))
	if err != nil {
		return path, err
	}
	if !strings.HasPrefix(path, settings.Config.PublicFolder) {
		return path, &ErrAccessDenied{}
	}
	return path, nil
}

func CheckFile(path string) (int, string, string, os.FileInfo, error) {
	path, err := validatePath(path)
	if err != nil {
		return settings.ERROR_ACCESS_DENIED, "You are not allowed to access this file", "", nil, &ErrAccessDenied{}
	}
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return settings.ERROR_FILE_NOT_EXISTS, "File or folder does not exist", path, stat, err
		}
		return settings.ERROR_WHILE_STAT, "Unknown error while stat", path, stat, err
	}
	return settings.ERROR_OK, "", path, stat, nil
}

func CalculateDirSize(path string) (int, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	if !stat.IsDir() {
		return 0, &ErrIsNotDir{}
	}
	var totalSize int
	filesList, err := ioutil.ReadDir(path)
	if err != nil {
		return 0, err
	}
	for _, file := range filesList {
		if file.IsDir() {
			thisSize, err := CalculateDirSize(path + "/" + file.Name())
			if err != nil {
				return 0, err
			}
			totalSize += thisSize
		} else {
			totalSize += int(file.Size())
		}
	}
	return totalSize, nil
}

func ConcatBytes(arrays ...[]byte) []byte {
	result := make([]byte, 0)
	for _, i := range arrays {
		result = append(result, i...)
	}
	return result
}

func Uint32ArrayContains(array []uint32, element uint32) bool {
	for _, i := range array {
		if i == element {
			return true
		}
	}
	return false
}
