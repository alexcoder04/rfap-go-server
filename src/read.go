package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

func Info(path string, requestDetails []string) (HeaderMetadata, error) {
	metadata := HeaderMetadata{}
	metadata.Path = path

	path, err := filepath.EvalSymlinks(PUBLIC_FOLDER + path)
	if err != nil {
		return retError(metadata, ERROR_UNKNOWN, "Unknown error while readlink"), err
	}
	if !strings.HasPrefix(path, PUBLIC_FOLDER) {
		return retError(metadata, ERROR_ACCESS_DENIED, "You are not permitted to read this folder"), &ErrAccessDenied{}
	}

	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return retError(metadata, ERROR_FILE_NOT_EXISTS, "File or folder does not exist"), err
		}
		return retError(metadata, ERROR_UNKNOWN, "Unknown error while stat"), err
	}

	metadata.Modified = int(stat.ModTime().Unix())
	if stat.IsDir() {
		metadata.Type = "d"
		for _, r := range requestDetails {
			switch r {
			case "DirectorySize":
				size, err := CalculateDirSize(path)
				if err != nil {
					return retError(metadata, ERROR_UNKNOWN, "Cannot calculate directory size"), &ErrCalculationFailed{}
				}
				metadata.DirectorySize = size
				break
			case "ElementsNumber":
				filesList, err := ioutil.ReadDir(path)
				if err != nil {
					return retError(metadata, ERROR_UNKNOWN, "Cannot list directory"), err
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
	return metadata, nil
}

func ReadFile(path string) (HeaderMetadata, []byte, error) {
	metadata := HeaderMetadata{}
	metadata.Path = path

	path, err := filepath.EvalSymlinks(PUBLIC_FOLDER + path)
	if err != nil {
		return retError(metadata, ERROR_UNKNOWN, "Unknown error while readlink"), make([]byte, 0), err
	}
	if !strings.HasPrefix(path, PUBLIC_FOLDER) {
		return retError(metadata, ERROR_ACCESS_DENIED, "You are not permitted to read this folder"), make([]byte, 0), &ErrAccessDenied{}
	}

	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return retError(metadata, ERROR_FILE_NOT_EXISTS, "File or folder does not exist"), make([]byte, 0), err
		}
		return retError(metadata, ERROR_UNKNOWN, "Unknown error while stat file"), make([]byte, 0), err
	}

	if stat.IsDir() {
		metadata.Type = "d"
		return retError(metadata, ERROR_INVALID_FILE_TYPE, "Is a directory"), make([]byte, 0), &ErrIsDir{}
	}
	metadata.Type = "f"
	metadata.FileSize = int(stat.Size())

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return retError(metadata, ERROR_UNKNOWN, "Cannot read file"), make([]byte, 0), err
	}

	mtype := mimetype.Detect(content)
	metadata.FileType = mtype.String()
	metadata.ErrorCode = 0
	metadata.Modified = int(stat.ModTime().Unix())
	return metadata, content, nil
}

func ReadDirectory(path string, requestDetails []string) (HeaderMetadata, []byte, error) {
	metadata := HeaderMetadata{}
	metadata.Path = path

	path, err := filepath.EvalSymlinks(PUBLIC_FOLDER + path)
	if err != nil {
		return retError(metadata, ERROR_UNKNOWN, "Unknown error while readlink"), make([]byte, 0), err
	}
	if !strings.HasPrefix(path, PUBLIC_FOLDER) {
		return retError(metadata, ERROR_ACCESS_DENIED, "You are not permitted to read this folder"), make([]byte, 0), &ErrAccessDenied{}
	}

	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return retError(metadata, ERROR_FILE_NOT_EXISTS, "Folder does not exist"), make([]byte, 0), err
		}
		return retError(metadata, ERROR_UNKNOWN, "Unknown error while stat"), make([]byte, 0), err
	}

	if !stat.Mode().IsDir() {
		metadata.Type = "f"
		return retError(metadata, ERROR_INVALID_FILE_TYPE, "Is not a directory"), make([]byte, 0), &ErrIsNotDir{}
	}

	filesList, err := ioutil.ReadDir(path)
	if err != nil {
		return retError(metadata, ERROR_UNKNOWN, "Cound not list folder"), make([]byte, 0), err
	}

	var result string
	for _, file := range filesList {
		result += "\n" + file.Name()
	}

	for _, r := range requestDetails {
		switch r {
		case "DirectorySize":
			size, err := CalculateDirSize(path)
			if err != nil {
				return retError(metadata, ERROR_UNKNOWN, "Cannot calculate directory size"), []byte(result), err
			}
			metadata.DirectorySize = size
			break
		case "ElementsNumber":
			filesList, err := ioutil.ReadDir(path)
			if err != nil {
				return retError(metadata, ERROR_UNKNOWN, "Cannot list directory"), []byte(result), err
			}
			metadata.DirectorySize = len(filesList)
			break
		}
	}

	return metadata, []byte(result), nil
}
