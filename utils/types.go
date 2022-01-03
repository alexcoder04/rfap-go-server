package utils

import (
	"net"

	"golang.org/x/crypto/openpgp"
)

type HeaderMetadata struct {
	ErrorCode    int    `yaml:"ErrorCode"`
	ErrorMessage string `yaml:"ErrorMessage"`

	RequestDetails []string `yaml:"RequestDetails"`

	Path        string `yaml:"Path"`
	Type        string `yaml:"Type"`
	Modified    int    `yaml:"Modified"`
	Destination string `yaml:"Destination"`

	FileSize int    `yaml:"FileSize"`
	FileType string `yaml:"FileType"`

	DirectorySize  int `yaml:"DirectorySize"`
	ElementsNumber int `yaml:"ElementsNumber"`
}

type Client struct {
	Conn         net.Conn
	PubkeyEntity *openpgp.Entity
	Address      string
}

type CommandExec func(string) (HeaderMetadata, []byte, error)

type CopySommandExec func(string, string, bool) (HeaderMetadata, []byte, error)
