//go:build windows

package main

import (
	"fmt"
	"os"
)

const (
	PUBLIC_FOLDER = "./shared"
)

var LOG_FILE = fmt.Sprintf("./rfap-go-server-%d.log", os.Getpid())
