package main

import (
	"net"
	"runtime"
	"time"

	"github.com/alexcoder04/rfap-go-server/log"
	"github.com/alexcoder04/rfap-go-server/network"
	"github.com/alexcoder04/rfap-go-server/settings"
	"github.com/sirupsen/logrus"
)

func main() {
	Init()
	log.Logger.Info("Starting " + settings.Config.ConnType + " server on " + settings.Config.ConnHost + ":" + settings.Config.ConnPort + " sharing " + settings.Config.PublicFolder + "...")

	l, err := net.Listen(settings.Config.ConnType, settings.Config.ConnHost+":"+settings.Config.ConnPort)
	if err != nil {
		log.Logger.Fatal("Error listening: ", err.Error())
	}

	defer l.Close()

	for {
		// wait and don't accept new connections if max number of clients already connected
		if runtime.NumGoroutine() >= settings.Config.MaxClients() {
			log.Logger.Warning("running threads: ", runtime.NumGoroutine(), "/", settings.Config.MaxClients)
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

		go network.HanleConnection(c)
		log.Logger.Info("running threads: ", runtime.NumGoroutine(), "/", settings.Config.MaxClients)
	}
}
