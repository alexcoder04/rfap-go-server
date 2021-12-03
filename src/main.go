package main

import (
	"net"

	"github.com/sirupsen/logrus"
)

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
			logger.WithFields(logrus.Fields{
				"client": c.RemoteAddr().String(),
			}).Warning("error connecting: ", err.Error())
			c.Close()
			return
		}
		logger.WithFields(logrus.Fields{
			"client": c.RemoteAddr().String(),
		}).Info("connected, starting thread to handle...")

		go HanleConnection(c)
	}
}
