//go:build !windows

package settings

import (
	"fmt"
	"os"
)

const PUBLIC_FOLDER = "/srv/rfap/shared"

var LOG_FILE = fmt.Sprintf("/tmp/rfap-go-server-%d.log", os.Getpid())
