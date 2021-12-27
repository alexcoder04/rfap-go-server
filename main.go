package main

import (
	"net"
	"runtime"
	"time"

	"github.com/alexcoder04/rfap-go-server/log"
	"github.com/alexcoder04/rfap-go-server/settings"
	"github.com/sirupsen/logrus"
)

func main() {
	Init()
	log.Logger.Info("Starting " + settings.CONN_TYPE + " server on " + settings.CONN_HOST + ":" + settings.CONN_PORT + " sharing " + settings.PUBLIC_FOLDER + "...")

	l, err := net.Listen(settings.CONN_TYPE, settings.CONN_HOST+":"+settings.CONN_PORT)
	if err != nil {
		log.Logger.Fatal("Error listening: ", err.Error())
	}

	defer l.Close()

	for {
		// wait and don't accept new connections if max number of clients already connected
		if runtime.NumGoroutine() >= settings.MAX_CLIENTS {
			log.Logger.Warning("running threads: ", runtime.NumGoroutine(), "/", settings.MAX_CLIENTS)
			time.Sleep(settings.MAX_THREADS_WAIT_SECS * time.Second)
			continue
		}

		c, err := l.Accept()
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"client": c.RemoteAddr().String(),
			}).Warning("error connecting: ", err.Error())
			c.Close()
			return
		}

		log.Logger.WithFields(logrus.Fields{
			"client": c.RemoteAddr().String(),
		}).Info("connected, starting thread to handle...")

		go HanleConnection(c)
		log.Logger.Info("running threads: ", runtime.NumGoroutine(), "/", settings.MAX_CLIENTS)
	}
}
