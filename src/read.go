package main

import (
	"errors"
	"io/ioutil"
	"os"
)

func Info(originalPath string) (HeaderValues, error) {
	path := PublicFolder + originalPath

	h := HeaderValues{}

	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			h.ErrorCode = 1
			h.ErrorMessage = "File or folder does not exist"
			return h, errors.New("file or folder does not exist")
		}
		h.ErrorCode = 2
		h.ErrorMessage = "Unknown error while stat"
		return h, errors.New("unknown error")
	}

	h.ErrorCode = 0
	h.Path = originalPath
	if stat.IsDir() {
		h.Type = 'd'
	} else {
		h.Type = 'f'
	}
	h.Modified = int(stat.ModTime().Unix())

	// TODO DirectorySize and ElementsNumber
	return h, nil
}

func ReadFile(path string) ([]byte, error) {
	path = PublicFolder + path

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make([]byte, 0), err
		}
		return make([]byte, 0), errors.New("unknown error")
	}

	// TODO it could be a directory
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return make([]byte, 0), err
	}

	return content, nil
}
