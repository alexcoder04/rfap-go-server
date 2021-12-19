package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func CopyFile(source string, destin string) (HeaderMetadata, error) {
	metadata := HeaderMetadata{}
	metadata.Path = source

	source, err := filepath.EvalSymlinks(PUBLIC_FOLDER + source)
	if err != nil {
		metadata.ErrorCode = ERROR_UNKNOWN
		metadata.ErrorMessage = "Unknown error while readlink"
		return metadata, err
	}
	if !strings.HasPrefix(source, PUBLIC_FOLDER) {
		metadata.ErrorCode = ERROR_ACCESS_DENIED
		metadata.ErrorMessage = "You are not permitted to read this file"
		return metadata, &ErrAccessDenied{}
	}

	destin, err = filepath.EvalSymlinks(PUBLIC_FOLDER + destin)
	if err != nil {
		metadata.ErrorCode = ERROR_UNKNOWN
		metadata.ErrorMessage = "Unknown error while readlink"
		return metadata, err
	}
	if !strings.HasPrefix(destin, PUBLIC_FOLDER) {
		metadata.ErrorCode = ERROR_ACCESS_DENIED
		metadata.ErrorMessage = "You are not permitted to write to this file"
		return metadata, &ErrAccessDenied{}
	}

	stat, err := os.Stat(source)
	if err != nil {
		if os.IsNotExist(err) {
			metadata.ErrorCode = ERROR_FILE_NOT_EXISTS
			metadata.ErrorMessage = "File or folder does not exist"
		} else {
			metadata.ErrorCode = ERROR_UNKNOWN
			metadata.ErrorMessage = "Unknown error while stat"
		}
		return metadata, err
	}
	if stat.IsDir() {
		metadata.ErrorCode = ERROR_INVALID_FILE_TYPE
		metadata.ErrorMessage = "Is a directory"
		return metadata, &ErrIsDir{}
	}

	_, err = os.Stat(destin)
	if err == nil {
		metadata.ErrorCode = ERROR_FILE_EXISTS
		metadata.ErrorMessage = "File already exists"
		return metadata, os.ErrExist
	}
	if !os.IsNotExist(err) {
		metadata.ErrorCode = ERROR_UNKNOWN
		metadata.ErrorMessage = "Unknown error while stat file"
		return metadata, err
	}

	bytesRead, err := ioutil.ReadFile(source)
	if err != nil {
		metadata.ErrorCode = ERROR_UNKNOWN
		metadata.ErrorMessage = "Unknown error while read file"
		return metadata, err
	}
	err = ioutil.WriteFile(destin, bytesRead, 0644)
	if err != nil {
		metadata.ErrorCode = ERROR_UNKNOWN
		metadata.ErrorMessage = "Unknown error while write file"
		return metadata, err
	}

	metadata.ErrorCode = 0
	metadata.Type = "f"

	return metadata, nil
}
