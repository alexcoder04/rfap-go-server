package commands

import (
	"os"

	"github.com/alexcoder04/rfap-go-server/settings"
	"github.com/alexcoder04/rfap-go-server/utils"
)

func DeleteFile(path string) (utils.HeaderMetadata, []byte, error) {
	metadata := utils.HeaderMetadata{}
	metadata.Path = path
	body := make([]byte, 0)

	errCode, errMsg, path, stat, err := utils.CheckFile(path)
	if err != nil {
		return utils.RetError(metadata, errCode, errMsg), body, err
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

	errCode, errMsg, path, stat, err := utils.CheckFile(path)
	if err != nil {
		return utils.RetError(metadata, errCode, errMsg), body, err
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
