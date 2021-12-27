package commands

import (
	"io/ioutil"
	"os"

	"github.com/alexcoder04/rfap-go-server/settings"
	"github.com/alexcoder04/rfap-go-server/utils"
	"github.com/otiai10/copy"
)

func CopyFile(source string, destin string, move bool) (utils.HeaderMetadata, []byte, error) {
	metadata := utils.HeaderMetadata{}
	metadata.Path = source
	body := make([]byte, 0)

	source, err := utils.ValidatePath(source)
	if err != nil {
		return utils.RetError(metadata, settings.ERROR_ACCESS_DENIED, "You are not permitted to access this file"), body, err
	}
	destin, err = utils.ValidatePath(destin)
	if err != nil {
		return utils.RetError(metadata, settings.ERROR_ACCESS_DENIED, "You are not permitted to access to this file"), body, &utils.ErrAccessDenied{}
	}

	errCode, errMsg, stat, err := utils.CheckFile(source)
	if err != nil {
		return utils.RetError(metadata, errCode, errMsg), body, err
	}
	if stat.IsDir() {
		metadata.Type = "d"
		return utils.RetError(metadata, settings.ERROR_INVALID_FILE_TYPE, "Is a directory"), body, &utils.ErrIsDir{}
	}

	_, err = os.Stat(destin)
	if err == nil {
		return utils.RetError(metadata, settings.ERROR_FILE_EXISTS, "File already exists"), body, os.ErrExist
	}
	if !os.IsNotExist(err) {
		return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Unknown error while stat file"), body, err
	}

	bytesRead, err := ioutil.ReadFile(source)
	if err != nil {
		return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Unknown error while read file"), body, err
	}
	err = ioutil.WriteFile(destin, bytesRead, 0600)
	if err != nil {
		return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Unknown error while write file"), body, err
	}
	if move {
		err = os.Remove(source)
		if err != nil {
			return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Cannot delete source file"), body, err
		}
	}

	metadata.ErrorCode = 0

	return metadata, body, nil
}

func CopyDirectory(source string, destin string, move bool) (utils.HeaderMetadata, []byte, error) {
	metadata := utils.HeaderMetadata{}
	metadata.Path = source
	body := make([]byte, 0)

	source, err := utils.ValidatePath(source)
	if err != nil {
		return utils.RetError(metadata, settings.ERROR_ACCESS_DENIED, "You are not permitted to access this file"), body, err
	}

	destin, err = utils.ValidatePath(destin)
	if err != nil {
		return utils.RetError(metadata, settings.ERROR_ACCESS_DENIED, "You are not permitted to access to this file"), body, &utils.ErrAccessDenied{}
	}

	errCode, errMsg, stat, err := utils.CheckFile(source)
	if err != nil {
		return utils.RetError(metadata, errCode, errMsg), body, err
	}
	if !stat.IsDir() {
		metadata.Type = "f"
		return utils.RetError(metadata, settings.ERROR_INVALID_FILE_TYPE, "Is not a directory"), body, &utils.ErrIsNotDir{}
	}

	_, err = os.Stat(destin)
	if err == nil {
		return utils.RetError(metadata, settings.ERROR_FILE_EXISTS, "File already exists"), body, os.ErrExist
	}
	if !os.IsNotExist(err) {
		return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Unknown error while stat file"), body, err
	}

	err = copy.Copy(source, destin)
	if err != nil {
		return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Cannot copy directory"), body, err
	}
	if move {
		err = os.RemoveAll(source)
		if err != nil {
			return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Cannot delete source directory"), body, err
		}
	}

	metadata.ErrorCode = 0

	return metadata, body, nil
}
