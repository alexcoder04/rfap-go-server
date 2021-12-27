//go:build windows

package settings

import (
	"fmt"
	"os"
)

const PUBLIC_FOLDER = "./shared"

var LOG_FILE = fmt.Sprintf("./rfap-go-server-%d.log", os.Getpid())
