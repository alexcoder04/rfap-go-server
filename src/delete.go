package main

import (
	"os"
)

func DeleteFile(path string) (HeaderMetadata, []byte, error) {
	metadata := HeaderMetadata{}
	metadata.Path = path
	body := make([]byte, 0)

	path, err := ValidatePath(path)
	if err != nil {
		return retError(metadata, ERROR_ACCESS_DENIED, "You are not permitted to access this file"), body, err
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

func DeleteDirectory(path string) (HeaderMetadata, []byte, error) {
	metadata := HeaderMetadata{}
	metadata.Path = path
	body := make([]byte, 0)

	path, err := ValidatePath(path)
	if err != nil {
		return retError(metadata, ERROR_ACCESS_DENIED, "You are not permitted to access this file"), body, err
	}

	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return retError(metadata, ERROR_FILE_NOT_EXISTS, "Folder does not exist"), body, err
		}
		return retError(metadata, ERROR_UNKNOWN, "Unknown error while stat"), body, err
	}

	if !stat.IsDir() {
		metadata.Type = "f"
		return retError(metadata, ERROR_INVALID_FILE_TYPE, "Is not a directory"), body, &ErrIsDir{}
	}

	metadata.Type = "d"

	err = os.RemoveAll(path)
	if err != nil {
		return retError(metadata, ERROR_UNKNOWN, "Cannot delete folder"), body, err
	}

	return metadata, body, nil
}
