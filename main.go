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

	log.Logger.WithFields(logrus.Fields{
		"host":          settings.Config.ConnHost,
		"port":          settings.Config.ConnPort,
		"shared folder": settings.Config.PublicFolder,
	}).Info("Starting  server ...")

	l, err := net.Listen(settings.Config.ConnType, settings.Config.ConnHost+":"+settings.Config.ConnPort)
	if err != nil {
		log.Logger.Fatal("Error listening: ", err.Error())
	}

	defer l.Close()

	for {
		// wait and don't accept new connections if max number of clients already connected
		if runtime.NumGoroutine() >= settings.Config.MaxClients() {
			log.Logger.Warning("running threads: ", runtime.NumGoroutine(), "/", settings.Config.MaxClients)
			time.Sleep(time.Duration(settings.Config.SecsWaitIfMaxThreads) * time.Second)
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
