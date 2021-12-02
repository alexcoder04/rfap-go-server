package main

import (
	"errors"
	"io/ioutil"
	"os"
)

func FileOrDirExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func CalculateDirSize(path string) (int, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	if !stat.IsDir() {
		return 0, errors.New("Is not a directory")
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
