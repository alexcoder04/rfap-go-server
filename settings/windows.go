//go:build windows

package settings

import (
	"fmt"
	"os"
	"path/filepath"
)

func getPublicFolder() string {
	return filepath.Abs("./shared")
}

func getLogFile() string {
	return filepath.Abs(fmt.Sprintf("./rfap-go-server-%d.log", od.Getpid()))
}
