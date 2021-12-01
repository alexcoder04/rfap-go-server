package main

import (
	"os"
)

func FileOrDirExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
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
