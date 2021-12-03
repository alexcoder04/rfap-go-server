package main

import (
	"net"
)

const (
	connHost = "localhost"
	connPort = "6700"
	connType = "tcp"
)

var SUPPORTED_RFAP_VERSIONS = []uint32{1}

func main() {
	Init()
	logger.Info("Starting " + connType + " server on " + connHost + ":" + connPort)

	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		logger.Fatal("Error listening: ", err.Error())
	}

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			logger.Warning(c.RemoteAddr().String(), " error connecting: ", err.Error())
			c.Close()
			return
		}
		logger.Info(c.RemoteAddr().String(), " connected, starting thread to handle...")

		go HanleConnection(c)
	}
}
