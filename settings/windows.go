//go:build windows

package settings

import (
	"fmt"
	"os"
	"path/filepath"
)

func getPublicFolder() string {
	path, err := filepath.Abs("./shared")
	if err != nil {
		fmt.Println("cannot resolve ./shared")
		os.Exit(1)
	}
	return path
}

func getLogFile() string {
	path, err := filepath.Abs(fmt.Sprintf("./rfap-go-server-%d.log", os.Getpid()))
	if err != nil {
		fmt.Println("cannot resolve ./shared")
		os.Exit(1)
	}
	return path
}
