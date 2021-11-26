package main

import (
	"errors"
	"io/ioutil"
	"os"
)

func Read(path string) ([]byte, error) {
	path = publicFolder + path
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make([]byte, 0), errors.New("file does not exist")
		}
	}
	if stat.IsDir() {
		returnString := ""
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			returnString = returnString + "\n" + f.Name()
		}
		return []byte(returnString), nil
	} else {
		content, _ := ioutil.ReadFile(path)
		return content, nil
	}
}
