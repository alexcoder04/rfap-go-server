package main

import (
	"os"
	"path/filepath"
	"strings"
)

func DeleteFile(path string) (HeaderMetadata, []byte, error) {
	metadata := HeaderMetadata{}
	metadata.Path = path
	body := make([]byte, 0)

	path, err := filepath.EvalSymlinks(PUBLIC_FOLDER + path)
	if err != nil {
		return retError(metadata, ERROR_UNKNOWN, "Unknown error while readlink"), body, err
	}
	if !strings.HasPrefix(path, PUBLIC_FOLDER) {
		return retError(metadata, ERROR_ACCESS_DENIED, "You are not permitted to delete this file"), body, &ErrAccessDenied{}
	}

	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return retError(metadata, ERROR_FILE_NOT_EXISTS, "File does not exist"), body, err
		}
		return retError(metadata, ERROR_UNKNOWN, "Unknown error while stat"), body, err
	}

	if stat.IsDir() {
		metadata.Type = "d"
		return retError(metadata, ERROR_INVALID_FILE_TYPE, "Is a directory"), body, &ErrIsDir{}
	}
	metadata.Type = "f"
	metadata.FileSize = int(stat.Size())

	err = os.Remove(path)
	if err != nil {
		return retError(metadata, ERROR_UNKNOWN, "Cannot delete file"), body, err
	}

	return metadata, body, nil
}
