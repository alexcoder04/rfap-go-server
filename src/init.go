package main

import (
	"os"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var logger *logrus.Logger

func Init() {
	operationMode := os.Getenv("RFAP_MODE")
	if operationMode == "testing" {
		InitStdoutLogger()
	} else {
		InitFileLogger()
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

func InitStdoutLogger() {
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
}

func InitFileLogger() {
	formatter := &prefixed.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	}

	file, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		logger = &logrus.Logger{
			Out:       os.Stdout,
			Level:     logrus.DebugLevel,
			Formatter: formatter,
		}
		logger.Fatal("Cannot open log file")
	}
	logger = &logrus.Logger{
		Out:       file,
		Level:     logrus.DebugLevel,
		Formatter: formatter,
	}

}

func CreateSharedFolder() {
	err := os.MkdirAll(PublicFolder, 0700)
	if err != nil {
		logger.Fatal("Cannot create shared folder: ", err.Error())
	}
	logger.Warning("Created shared folder", PublicFolder)
}
