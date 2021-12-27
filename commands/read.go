package commands

import (
	"io/ioutil"
	"os"

	"github.com/alexcoder04/rfap-go-server/settings"
	"github.com/alexcoder04/rfap-go-server/utils"
	"github.com/gabriel-vasile/mimetype"
)

func Info(path string, requestDetails []string) (utils.HeaderMetadata, []byte, error) {
	metadata := utils.HeaderMetadata{}
	metadata.Path = path
	body := make([]byte, 0)

	path, err := utils.ValidatePath(path)
	if err != nil {
		return utils.RetError(metadata, settings.ERROR_ACCESS_DENIED, "You are not permitted to access this folder"), body, err
	}

	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return utils.RetError(metadata, settings.ERROR_FILE_NOT_EXISTS, "File or folder does not exist"), body, err
		}
		return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Unknown error while stat"), body, err
	}

	metadata.Modified = int(stat.ModTime().Unix())
	if stat.IsDir() {
		metadata.Type = "d"
		for _, r := range requestDetails {
			switch r {
			case "DirectorySize":
				size, err := utils.CalculateDirSize(path)
				if err != nil {
					return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Cannot calculate directory size"), body, &utils.ErrCalculationFailed{}
				}
				metadata.DirectorySize = size
				break
			case "ElementsNumber":
				filesList, err := ioutil.ReadDir(path)
				if err != nil {
					return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Cannot list directory"), body, err
				}
				metadata.ElementsNumber = len(filesList)
				break
			}
		}
	} else {
		metadata.Type = "f"
		metadata.FileSize = int(stat.Size())
		mtype, err := mimetype.DetectFile(path)
		if err != nil {
			metadata.FileType = "application/octet-stream"
		} else {
			metadata.FileType = mtype.String()
		}
	}

	metadata.ErrorCode = 0
	return metadata, body, nil
}

func ReadFile(path string) (utils.HeaderMetadata, []byte, error) {
	metadata := utils.HeaderMetadata{}
	metadata.Path = path

	path, err := utils.ValidatePath(path)
	if err != nil {
		return utils.RetError(metadata, settings.ERROR_ACCESS_DENIED, "You are not permitted to access this file"), make([]byte, 0), err
	}

	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return utils.RetError(metadata, settings.ERROR_FILE_NOT_EXISTS, "File or folder does not exist"), make([]byte, 0), err
		}
		return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Unknown error while stat file"), make([]byte, 0), err
	}

	if stat.IsDir() {
		metadata.Type = "d"
		return utils.RetError(metadata, settings.ERROR_INVALID_FILE_TYPE, "Is a directory"), make([]byte, 0), &utils.ErrIsDir{}
	}
	metadata.Type = "f"
	metadata.FileSize = int(stat.Size())

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Cannot read file"), make([]byte, 0), err
	}

	mtype := mimetype.Detect(content)
	metadata.FileType = mtype.String()
	metadata.ErrorCode = 0
	metadata.Modified = int(stat.ModTime().Unix())
	return metadata, content, nil
}

func ReadDirectory(path string, requestDetails []string) (utils.HeaderMetadata, []byte, error) {
	metadata := utils.HeaderMetadata{}
	metadata.Path = path

	path, err := utils.ValidatePath(path)
	if err != nil {
		return utils.RetError(metadata, settings.ERROR_ACCESS_DENIED, "You are not permitted to access this file"), make([]byte, 0), err
	}

	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return utils.RetError(metadata, settings.ERROR_FILE_NOT_EXISTS, "Folder does not exist"), make([]byte, 0), err
		}
		return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Unknown error while stat"), make([]byte, 0), err
	}

	if !stat.Mode().IsDir() {
		metadata.Type = "f"
		return utils.RetError(metadata, settings.ERROR_INVALID_FILE_TYPE, "Is not a directory"), make([]byte, 0), &utils.ErrIsNotDir{}
	}

	filesList, err := ioutil.ReadDir(path)
	if err != nil {
		return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Cound not list folder"), make([]byte, 0), err
	}

	var result string
	for _, file := range filesList {
		result += "\n" + file.Name()
	}

	for _, r := range requestDetails {
		switch r {
		case "DirectorySize":
			size, err := utils.CalculateDirSize(path)
			if err != nil {
				return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Cannot calculate directory size"), []byte(result), err
			}
			metadata.DirectorySize = size
			break
		case "ElementsNumber":
			filesList, err := ioutil.ReadDir(path)
			if err != nil {
				return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Cannot list directory"), []byte(result), err
			}
			metadata.DirectorySize = len(filesList)
			break
		}
	}

	return metadata, []byte(result), nil
}
