//go:build !windows

package main

import (
	"fmt"
	"os"
)

const PUBLIC_FOLDER = "/srv/rfap/shared"

var LOG_FILE = fmt.Sprintf("/tmp/rfap-go-server-%d.log", os.Getpid())
