package main

import (
	"os"
	"path/filepath"
	"strings"
)

func CreateFile(path string) (HeaderMetadata, error) {
	metadata := HeaderMetadata{}
	metadata.Path = path

	path, err := filepath.EvalSymlinks(PUBLIC_FOLDER + path)
	if err != nil {
		metadata.ErrorCode = ERROR_UNKNOWN
		metadata.ErrorMessage = "Unknown error while readlink"
		return metadata, err
	}
	if !strings.HasPrefix(path, PUBLIC_FOLDER) {
		metadata.ErrorCode = ERROR_ACCESS_DENIED
		metadata.ErrorMessage = "You are not permitted to create this file"
		return metadata, &ErrAccessDenied{}
	}

	_, err = os.Stat(path)
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

	f, err := os.Create(path)
	if err != nil {
		metadata.ErrorCode = ERROR_UNKNOWN
		metadata.ErrorMessage = "Cannot create file"
		return metadata, err
	}
	f.Close()
	metadata.ErrorCode = 0
	metadata.Type = "f"

	return metadata, nil
}
