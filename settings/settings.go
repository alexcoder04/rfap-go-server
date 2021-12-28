package settings

import (
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

type ServerConfiguration struct {
	ConnHost string
	ConnPort string
	ConnType string

	MaxClientsPerCore int

	PublicFolder string

	LogFile     string
	LogLevelStr string
	LogFormat   string
}

func (config *ServerConfiguration) MaxClients() int {
	return config.MaxClientsPerCore * runtime.NumCPU()
}

func (config *ServerConfiguration) LogLevel() logrus.Level {
	switch config.LogLevelStr {
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

func (config *ServerConfiguration) LoadDefaultConfig() {
	config.ConnHost = "localhost"
	config.ConnPort = "6700"
	config.ConnType = "tcp"

	config.MaxClientsPerCore = 4

	config.PublicFolder = getPublicFolder()

	config.LogFile = getLogFile()
	config.LogLevelStr = "info"
	config.LogFormat = "default"
}

func (config *ServerConfiguration) ApplyEnvConfig() {
	if connHost := os.Getenv("RFAP_CONN_HOST"); connHost != "" {
		config.ConnHost = connHost
	}
	if connPort := os.Getenv("RFAP_CONN_PORT"); connPort != "" {
		config.ConnPort = connPort
	}

	if publicFolder := os.Getenv("RFAP_PUBLIC_FOLDER"); publicFolder != "" {
		config.PublicFolder = publicFolder
	}

	if logFile := os.Getenv("RFAP_LOG_FILE"); logFile != "" {
		config.LogFile = logFile
	}
	if logLevel := os.Getenv("RFAP_LOG_LEVEL"); logLevel != "" {
		config.LogLevelStr = logLevel
	}
	if logFormat := os.Getenv("RFAP_LOG_FORMAT"); logFormat != "" {
		config.LogFormat = logFormat
	}
}

var Config = &ServerConfiguration{}
