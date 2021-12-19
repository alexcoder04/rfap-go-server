package main

import (
	"os"
	"path/filepath"
	"strings"
)

func CreateFile(path string) (HeaderMetadata, []byte, error) {
	metadata := HeaderMetadata{}
	metadata.Path = path
	body := make([]byte, 0)

	path, err := filepath.EvalSymlinks(PUBLIC_FOLDER + path)
	if err != nil {
		return retError(metadata, ERROR_UNKNOWN, "Unknown error while readlink"), body, err
	}
	if !strings.HasPrefix(path, PUBLIC_FOLDER) {
		return retError(metadata, ERROR_ACCESS_DENIED, "You are not permitted to create this file"), body, &ErrAccessDenied{}
	}

	_, err = os.Stat(path)
	if err == nil {
		return retError(metadata, ERROR_FILE_EXISTS, "File already exists"), body, os.ErrExist
	}
	if !os.IsNotExist(err) {
		return retError(metadata, ERROR_UNKNOWN, "Unknown error while stat file"), body, err
	}

	f, err := os.Create(path)
	if err != nil {
		return retError(metadata, ERROR_UNKNOWN, "Cannot create file"), body, err
	}
	f.Close()
	metadata.ErrorCode = 0
	metadata.Type = "f"

	return metadata, body, nil
}
