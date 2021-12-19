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
		return retError(metadata, ERROR_UNKNOWN, "Unknown error while readlink"), err
	}
	if !strings.HasPrefix(path, PUBLIC_FOLDER) {
		return retError(metadata, ERROR_ACCESS_DENIED, "You are not permitted to create this file"), &ErrAccessDenied{}
	}

	_, err = os.Stat(path)
	if err == nil {
		return retError(metadata, ERROR_FILE_EXISTS, "File already exists"), os.ErrExist
	}
	if !os.IsNotExist(err) {
		return retError(metadata, ERROR_UNKNOWN, "Unknown error while stat file"), err
	}

	f, err := os.Create(path)
	if err != nil {
		return retError(metadata, ERROR_UNKNOWN, "Cannot create file"), err
	}
	f.Close()
	metadata.ErrorCode = 0
	metadata.Type = "f"

	return metadata, nil
}
