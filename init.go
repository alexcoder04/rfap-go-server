package main

import (
	"os"

	"github.com/alexcoder04/rfap-go-server/log"
	"github.com/alexcoder04/rfap-go-server/settings"
	"github.com/sirupsen/logrus"
)

func Init() {
	operationMode := os.Getenv("RFAP_MODE")
	if operationMode == "testing" {
		log.InitStdoutLogger()
	} else {
		log.InitFileLogger()
	}

	log.Logger.WithFields(logrus.Fields{
		"commit":     settings.GIT_COMMIT,
		"build time": settings.BUILD_TIMESTAMP,
		"version":    settings.SERVER_VERSION,
		"os":         settings.BUILD_OS,
	}).Info("build info")

	_, err := os.Stat(settings.PUBLIC_FOLDER)
	if err != nil {
		if os.IsNotExist(err) {
			log.Logger.Warning("Shared folder does not exist, creating...")
			CreateSharedFolder()
		} else {
			log.Logger.Fatal("Unknown error while stat shared folder: ", err.Error())
		}
	}
}

func CreateSharedFolder() {
	err := os.MkdirAll(settings.PUBLIC_FOLDER, 0700)
	if err != nil {
		log.Logger.Fatal("Cannot create shared folder: ", err.Error())
	}
	log.Logger.Warning("Created shared folder ", settings.PUBLIC_FOLDER)
}
