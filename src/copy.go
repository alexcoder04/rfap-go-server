package main

import (
	"io/ioutil"
	"os"
)

func CopyFile(source string, destin string, move bool) (HeaderMetadata, error) {
	metadata := HeaderMetadata{}
	metadata.Path = source

	source, err := ValidatePath(source)
	if err != nil {
		return retError(metadata, ERROR_ACCESS_DENIED, "You are not permitted to access this file"), err
	}

	destin, err = ValidatePath(destin)
	if err != nil {
		return retError(metadata, ERROR_ACCESS_DENIED, "You are not permitted to access to this file"), &ErrAccessDenied{}
	}

	stat, err := os.Stat(source)
	if err != nil {
		if os.IsNotExist(err) {
			return retError(metadata, ERROR_FILE_NOT_EXISTS, "File or folder does not exist"), err
		}
		return retError(metadata, ERROR_UNKNOWN, "Unknown error while stat"), err
	}
	if stat.IsDir() {
		metadata.Type = "d"
		return retError(metadata, ERROR_INVALID_FILE_TYPE, "Is a directory"), &ErrIsDir{}
	}

	_, err = os.Stat(destin)
	if err == nil {
		return retError(metadata, ERROR_FILE_EXISTS, "File already exists"), os.ErrExist
	}
	if !os.IsNotExist(err) {
		return retError(metadata, ERROR_UNKNOWN, "Unknown error while stat file"), err
	}

	bytesRead, err := ioutil.ReadFile(source)
	if err != nil {
		return retError(metadata, ERROR_UNKNOWN, "Unknown error while read file"), err
	}
	err = ioutil.WriteFile(destin, bytesRead, 0600)
	if err != nil {
		return retError(metadata, ERROR_UNKNOWN, "Unknown error while write file"), err
	}
	if move {
		err = os.Remove(source)
		if err != nil {
			return retError(metadata, ERROR_UNKNOWN, "Cannot delete source file"), err
		}
	}

	metadata.ErrorCode = 0
	metadata.Type = "f"

	return metadata, nil
}
