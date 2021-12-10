package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

func Info(path string, requestDetails []string) HeaderMetadata {
	metadata := HeaderMetadata{}
	metadata.Path = path

	path, err := filepath.EvalSymlinks(PUBLIC_FOLDER + path)
	if err != nil {
		metadata.ErrorCode = ERROR_UNKNOWN
		metadata.ErrorMessage = "Unknown error while readlink"
		return metadata
	}
	if !strings.HasPrefix(path, PUBLIC_FOLDER) {
		metadata.ErrorCode = ERROR_ACCESS_DENIED
		metadata.ErrorMessage = "You are not permitted to read this folder"
		return metadata
	}

	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			metadata.ErrorCode = ERROR_FILE_NOT_EXISTS
			metadata.ErrorMessage = "File or folder does not exist"
		} else {
			metadata.ErrorCode = ERROR_UNKNOWN
			metadata.ErrorMessage = "Unknown error while stat"
		}
		return metadata
	}

	metadata.Modified = int(stat.ModTime().Unix())
	if stat.IsDir() {
		metadata.Type = "d"
		for _, r := range requestDetails {
			switch r {
			case "DirectorySize":
				size, err := CalculateDirSize(path)
				if err != nil {
					metadata.ErrorCode = ERROR_UNKNOWN
					metadata.ErrorMessage = "Cannot calculate directory size"
					break
				}
				metadata.DirectorySize = size
				break
			case "ElementsNumber":
				filesList, err := ioutil.ReadDir(path)
				if err != nil {
					metadata.ErrorCode = ERROR_UNKNOWN
					metadata.ErrorMessage = "Cannot list directory"
					break
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
	return metadata
}

func ReadFile(path string) (HeaderMetadata, []byte, error) {
	metadata := HeaderMetadata{}
	metadata.Path = path

	path, err := filepath.EvalSymlinks(PUBLIC_FOLDER + path)
	if err != nil {
		metadata.ErrorCode = ERROR_UNKNOWN
		metadata.ErrorMessage = "Unknown error while readlink"
		return metadata, make([]byte, 0), err
	}
	if !strings.HasPrefix(path, PUBLIC_FOLDER) {
		metadata.ErrorCode = ERROR_ACCESS_DENIED
		metadata.ErrorMessage = "You are not permitted to read this folder"
		return metadata, make([]byte, 0), &ErrAccessDenied{}
	}

	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			metadata.ErrorCode = ERROR_FILE_NOT_EXISTS
			metadata.ErrorMessage = "File or folder does not exist"
		} else {
			metadata.ErrorCode = ERROR_UNKNOWN
			metadata.ErrorMessage = "Unknown error while stat file"
		}
		return metadata, make([]byte, 0), err
	}

	if stat.IsDir() {
		metadata.ErrorCode = ERROR_INVALID_FILE_TYPE
		metadata.ErrorMessage = "Is a directory"
		metadata.Type = "d"
		return metadata, make([]byte, 0), &ErrIsDir{}
	}
	metadata.Type = "f"
	metadata.FileSize = int(stat.Size())

	content, err := ioutil.ReadFile(path)
	if err != nil {
		metadata.ErrorCode = ERROR_UNKNOWN
		metadata.ErrorMessage = "Cannot read file"
		return metadata, make([]byte, 0), err
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
		metadata.ErrorCode = ERROR_UNKNOWN
		metadata.ErrorMessage = "Unknown error while readlink"
		return metadata, make([]byte, 0), err
	}
	if !strings.HasPrefix(path, PUBLIC_FOLDER) {
		metadata.ErrorCode = ERROR_ACCESS_DENIED
		metadata.ErrorMessage = "You are not permitted to read this folder"
		return metadata, make([]byte, 0), &ErrAccessDenied{}
	}

	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			metadata.ErrorCode = ERROR_FILE_NOT_EXISTS
			metadata.ErrorMessage = "Folder does not exist"
		} else {
			metadata.ErrorCode = ERROR_UNKNOWN
			metadata.ErrorMessage = "Error while stat"
		}
		return metadata, make([]byte, 0), err
	}

	if !stat.Mode().IsDir() {
		metadata.ErrorCode = ERROR_INVALID_FILE_TYPE
		metadata.ErrorMessage = "Is not a directory"
		return metadata, make([]byte, 0), &ErrIsNotDir{}
	}

	filesList, err := ioutil.ReadDir(path)
	if err != nil {
		metadata.ErrorCode = ERROR_UNKNOWN
		metadata.ErrorMessage = "Cound not list folder"
		return metadata, make([]byte, 0), err
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
				metadata.ErrorCode = ERROR_UNKNOWN
				metadata.ErrorMessage = "Cannot calculate directory size"
				break
			}
			metadata.DirectorySize = size
			break
		case "ElementsNumber":
			filesList, err := ioutil.ReadDir(path)
			if err != nil {
				metadata.ErrorCode = ERROR_UNKNOWN
				metadata.ErrorMessage = "Cannot list directory"
				break
			}
			metadata.DirectorySize = len(filesList)
			break
		}
	}

	return metadata, []byte(result), nil
}
