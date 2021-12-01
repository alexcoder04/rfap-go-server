package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
)

func Info(originalPath string) (HeaderMetadata, error) {
	path := PublicFolder + originalPath
	log.Println("reading info on", originalPath, "=", path)

	h := HeaderMetadata{}

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
		h.Type = "d"
	} else {
		h.Type = "f"
	}
	h.Modified = int(stat.ModTime().Unix())

	// TODO DirectorySize and ElementsNumber
	return h, nil
}

func ReadFile(path string) ([]byte, error) {
	path = PublicFolder + path

	stat, err := os.Stat(path)
	if err != nil {
		return make([]byte, 0), err
	}

	if stat.IsDir() {
		return make([]byte, 0), errors.New("is a directory")
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return make([]byte, 0), err
	}

	return content, nil
}
func ReadDirectory(path string) ([]byte, error) {
	path = PublicFolder + path

	stat, err := os.Stat(path)
	if err != nil {
		return make([]byte, 0), err
	}

	if !stat.Mode().IsDir() {
		return make([]byte, 0), errors.New("is not a directory")
	}

	filesList, err := ioutil.ReadDir(path)
	if err != nil {
		return make([]byte, 0), err
	}

	var result string
	for _, file := range filesList {
		result += "\n" + file.Name()
	}
	log.Println(result)

	return []byte(result), nil
}
