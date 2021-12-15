package main

import (
	"net"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	Init()
	logger.Info("Starting " + connType + " server on " + connHost + ":" + connPort + " sharing " + PUBLIC_FOLDER + "...")

	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		logger.Fatal("Error listening: ", err.Error())
	}

	defer l.Close()

	for {
		// wait and don't accept new connections if max number of clients already connected
		if runtime.NumGoroutine() >= MAX_CLIENTS {
			logger.Warning("running threads: ", runtime.NumGoroutine(), "/", MAX_CLIENTS)
			time.Sleep(MAX_THREADS_WAIT_SECS * time.Second)
			continue
		}

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
		logger.Info("running threads: ", runtime.NumGoroutine(), "/", MAX_CLIENTS)
	}
}
