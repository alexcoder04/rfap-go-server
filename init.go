package main

import (
	"os"

	"github.com/alexcoder04/rfap-go-server/log"
	"github.com/alexcoder04/rfap-go-server/settings"
	"github.com/sirupsen/logrus"
)

func Init() {
	settings.Config.LoadDefaultConfig()
	settings.Config.ApplyEnvConfig()

	if settings.Config.LogFile == "[stdout]" {
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

	_, err := os.Stat(settings.Config.PublicFolder)
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
	err := os.MkdirAll(settings.Config.PublicFolder, 0700)
	if err != nil {
		log.Logger.Fatal("Cannot create shared folder: ", err.Error())
	}
	log.Logger.Warning("Created shared folder ", settings.Config.PublicFolder)
}
