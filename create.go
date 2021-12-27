package main

import (
	"os"

	"github.com/alexcoder04/rfap-go-server/settings"
)

func CreateFile(path string) (HeaderMetadata, []byte, error) {
	metadata := HeaderMetadata{}
	metadata.Path = path
	body := make([]byte, 0)

	path, err := ValidatePath(path)
	if err != nil {
		return retError(metadata, settings.ERROR_ACCESS_DENIED, "You are not permitted to access this file"), body, err
	}

	_, err = os.Stat(path)
	if err == nil {
		return retError(metadata, settings.ERROR_FILE_EXISTS, "File already exists"), body, os.ErrExist
	}
	if !os.IsNotExist(err) {
		return retError(metadata, settings.ERROR_UNKNOWN, "Unknown error while stat file"), body, err
	}

	f, err := os.Create(path)
	if err != nil {
		return retError(metadata, settings.ERROR_UNKNOWN, "Cannot create file"), body, err
	}
	f.Close()
	metadata.ErrorCode = 0
	metadata.Type = "f"

	return metadata, body, nil
}

func CreateDirectory(path string) (HeaderMetadata, []byte, error) {
	metadata := HeaderMetadata{}
	metadata.Path = path
	body := make([]byte, 0)

	path, err := ValidatePath(path)
	if err != nil {
		return retError(metadata, settings.ERROR_ACCESS_DENIED, "You are not permitted to access this file"), body, err
	}

	_, err = os.Stat(path)
	if err == nil {
		return retError(metadata, settings.ERROR_FILE_EXISTS, "File already exists"), body, os.ErrExist
	}
	if !os.IsNotExist(err) {
		return retError(metadata, settings.ERROR_UNKNOWN, "Unknown error while stat file"), body, err
	}

	err = os.Mkdir(path, 0700)
	if err != nil {
		return retError(metadata, settings.ERROR_UNKNOWN, "Cannot create folder"), body, err
	}
	metadata.ErrorCode = 0
	metadata.Type = "d"

	return metadata, body, nil
}
