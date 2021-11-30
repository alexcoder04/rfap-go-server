package main

import (
	"fmt"
	"os"
)

func Init() {
	_, err := os.Stat(PublicFolder)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(PublicFolder, 0700)
			if err != nil {
				fmt.Println("Cannot create shared folder")
				os.Exit(1)
			}
			fmt.Println("Created shared folder")
			return
		}
		fmt.Println("unknown error while stat shared folder")
		os.Exit(1)
	}
}
