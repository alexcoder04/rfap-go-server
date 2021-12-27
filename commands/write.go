package commands

import (
	"os"
	"time"

	"github.com/alexcoder04/rfap-go-server/settings"
	"github.com/alexcoder04/rfap-go-server/utils"
	"github.com/gabriel-vasile/mimetype"
)

func WriteFile(path string, content []byte) (utils.HeaderMetadata, []byte, error) {
	metadata := utils.HeaderMetadata{}
	metadata.Path = path
	body := make([]byte, 0)

	errCode, errMsg, path, stat, err := utils.CheckFile(path)
	if errCode == settings.ERROR_ACCESS_DENIED {
		return utils.RetError(metadata, errCode, errMsg), body, err
	}
	if errCode != settings.ERROR_OK && errCode != settings.ERROR_FILE_NOT_EXISTS {
		return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Unknown error while stat"), body, err
	}

	if errCode == settings.ERROR_OK && stat.IsDir() {
		metadata.Type = "d"
		return utils.RetError(metadata, settings.ERROR_INVALID_FILE_TYPE, "Is a directory"), body, &utils.ErrIsDir{}
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Cannot open file"), body, err
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Cannot write to file"), body, err
	}

	metadata.Type = "f"
	metadata.FileSize = len(content)
	mtype := mimetype.Detect(content)
	metadata.FileType = mtype.String()
	metadata.Modified = int(time.Now().Unix())
	return metadata, body, nil
}
