//go:build !windows

package settings

import (
	"fmt"
	"os"
)

func getPublicFolder() string {
	return "/srv/rfap/shared"
}

func getLogFile() string {
	return fmt.Sprintf("/tmp/rfap-go-server-%d.log", os.Getpid())
}
