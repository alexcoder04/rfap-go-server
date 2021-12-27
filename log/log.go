package log

import (
	"os"

	"github.com/alexcoder04/rfap-go-server/settings"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var Logger *logrus.Logger

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

	Logger = &logrus.Logger{
		Out:       os.Stdout,
		Level:     getLogLevel(),
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

	file, err := os.OpenFile(settings.LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		Logger = &logrus.Logger{
			Out:       os.Stdout,
			Level:     getLogLevel(),
			Formatter: formatter,
		}
		Logger.Fatal("Cannot open log file")
	}
	Logger = &logrus.Logger{
		Out:       file,
		Level:     logrus.DebugLevel,
		Formatter: formatter,
	}

}

func getLogLevel() logrus.Level {
	switch os.Getenv("RFAP_LOG_LEVEL") {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}
