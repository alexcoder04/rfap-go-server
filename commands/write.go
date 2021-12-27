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

	path, err := utils.ValidatePath(path)
	if err != nil {
		return utils.RetError(metadata, settings.ERROR_ACCESS_DENIED, "You are not permitted to access to this file"), body, err
	}

	stat, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return utils.RetError(metadata, settings.ERROR_UNKNOWN, "Unknown error while stat"), body, err
	}

	if err == nil && stat.IsDir() {
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
