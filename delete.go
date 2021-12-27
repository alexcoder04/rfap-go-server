package main

import (
	"os"

	"github.com/alexcoder04/rfap-go-server/settings"
	"github.com/alexcoder04/rfap-go-server/utils"
)

func DeleteFile(path string) (utils.HeaderMetadata, []byte, error) {
	metadata := utils.HeaderMetadata{}
	metadata.Path = path
	body := make([]byte, 0)

	path, err := utils.ValidatePath(path)
	if err != nil {
		return utils.RetError(metadata, settings.ERROR_ACCESS_DENIED, "You are not permitted to access this file"), body, err
	}

	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return utils.RetError(metadata, settings.ERROR_FILE_NOT_EXISTS, "File does not exist"), body, err
		}
		return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Unknown error while stat"), body, err
	}

	if stat.IsDir() {
		metadata.Type = "d"
		return utils.RetError(metadata, settings.ERROR_INVALID_FILE_TYPE, "Is a directory"), body, &utils.ErrIsDir{}
	}
	metadata.Type = "f"
	metadata.FileSize = int(stat.Size())

	err = os.Remove(path)
	if err != nil {
		return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Cannot delete file"), body, err
	}

	return metadata, body, nil
}

func DeleteDirectory(path string) (utils.HeaderMetadata, []byte, error) {
	metadata := utils.HeaderMetadata{}
	metadata.Path = path
	body := make([]byte, 0)

	path, err := utils.ValidatePath(path)
	if err != nil {
		return utils.RetError(metadata, settings.ERROR_ACCESS_DENIED, "You are not permitted to access this file"), body, err
	}

	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return utils.RetError(metadata, settings.ERROR_FILE_NOT_EXISTS, "Folder does not exist"), body, err
		}
		return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Unknown error while stat"), body, err
	}

	if !stat.IsDir() {
		metadata.Type = "f"
		return utils.RetError(metadata, settings.ERROR_INVALID_FILE_TYPE, "Is not a directory"), body, &utils.ErrIsDir{}
	}

	metadata.Type = "d"

	err = os.RemoveAll(path)
	if err != nil {
		return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Cannot delete folder"), body, err
	}

	return metadata, body, nil
}
