package log

import (
	"os"

	"github.com/alexcoder04/rfap-go-server/settings"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var Logger *logrus.Logger

func getLogOut() *os.File {
	if settings.Config.LogFile == "[stdout]" {
		return os.Stdout
	}
	file, err := os.OpenFile(settings.Config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		logrus.Fatalln("Cannot open log file")
	}
	return file
}

func getFormatter() logrus.Formatter {
	switch settings.Config.LogFormat {
	case "color":
		formatter := &prefixed.TextFormatter{
			DisableColors:   false,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		}
		formatter.SetColorScheme(&prefixed.ColorScheme{
			TimestampStyle: "white",
		})
		return formatter
	case "json":
		formatter := &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		}
		return formatter
	default:
		formatter := &prefixed.TextFormatter{
			DisableColors:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		}
		return formatter
	}
}

func InitLogger() {
	Logger = &logrus.Logger{
		Out:       getLogOut(),
		Level:     settings.Config.LogLevel(),
		Formatter: getFormatter(),
	}
}
