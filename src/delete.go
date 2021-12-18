package main

import (
	"os"
	"path/filepath"
	"strings"
)

func DeleteFile(path string) (HeaderMetadata, error) {
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
		metadata.ErrorMessage = "You are not permitted to delete this file"
		return metadata, &ErrAccessDenied{}
	}

	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			metadata.ErrorCode = ERROR_FILE_NOT_EXISTS
			metadata.ErrorMessage = "File does not exist"
		} else {
			metadata.ErrorCode = ERROR_UNKNOWN
			metadata.ErrorMessage = "Unknown error while stat"
		}
		return metadata, err
	}

	if stat.IsDir() {
		metadata.ErrorCode = ERROR_INVALID_FILE_TYPE
		metadata.ErrorMessage = "Is a directory"
		metadata.Type = "d"
		return metadata, &ErrIsDir{}
	}
	metadata.Type = "f"
	metadata.FileSize = int(stat.Size())

	err = os.Remove(path)
	if err != nil {
		metadata.ErrorCode = ERROR_UNKNOWN
		metadata.ErrorMessage = "Cannot delete file"
		return metadata, err
	}

	return metadata, nil
}
