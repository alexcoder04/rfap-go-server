package main

import (
	"os"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var logger *logrus.Logger

func Init() {
	formatter := &prefixed.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	}
	formatter.SetColorScheme(&prefixed.ColorScheme{
		TimestampStyle: "white",
	})

	logger = &logrus.Logger{
		Out:       os.Stdout,
		Level:     logrus.TraceLevel,
		Formatter: formatter,
	}

	_, err := os.Stat(PublicFolder)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Warning("Shared folder does not exist, creating...")
			CreateSharedFolder()
		} else {
			logger.Fatal("Unknown error while stat shared folder: ", err.Error())
		}
	}
}

func CreateSharedFolder() {
	err := os.MkdirAll(PublicFolder, 0700)
	if err != nil {
		logger.Fatal("Cannot create shared folder: ", err.Error())
	}
	logger.Warning("Created shared folder", PublicFolder)
}
