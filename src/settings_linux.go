//go:build !windows

package main

import (
	"fmt"
	"os"
)

const PublicFolder = "/tmp/rfap-share"

var LOG_FILE = fmt.Sprintf("/tmp/rfap-go-server-%d.log", os.Getpid())
