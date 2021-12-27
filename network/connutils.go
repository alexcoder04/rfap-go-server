package network

import (
	"net"
	"runtime"

	"github.com/alexcoder04/rfap-go-server/log"
	"github.com/alexcoder04/rfap-go-server/settings"
	"github.com/alexcoder04/rfap-go-server/utils"
	"github.com/sirupsen/logrus"
)

func cleanErrorDisconnect(conn net.Conn) {
	header := utils.HeaderMetadata{}
	err := sendPacket(conn, settings.CMD_ERROR, header, make([]byte, 0))
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Error("send disconnect packet failed")
	}
	conn.Close()
	log.Logger.WithFields(logrus.Fields{
		"client": conn.RemoteAddr().String(),
	}).Info("connection closed")
	log.Logger.Info("running threads: ", runtime.NumGoroutine(), "/", settings.MAX_CLIENTS)
}

func runCommand(conn net.Conn, header utils.HeaderMetadata, cmd int, commandName string, fn utils.CommandExec) {
	log.Logger.WithFields(logrus.Fields{
		"client":  conn.RemoteAddr().String(),
		"command": commandName,
	}).Info("packet: ", header.Path)
	metadata, content, err := fn(header.Path)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Warning("error ", commandName, " ", header.Path, ": ", err.Error())
	}
	err = sendPacket(conn, cmd+1, metadata, content)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Error("error while response to ", commandName, ": ", err.Error())
	}
}

func runCopyCommand(conn net.Conn, header utils.HeaderMetadata, cmd int, commandName string, fn utils.CopySommandExec, move bool) {
	log.Logger.WithFields(logrus.Fields{
		"client":  conn.RemoteAddr().String(),
		"command": commandName,
	}).Info("packet: ", header.Path)
	metadata, content, err := fn(header.Path, header.Destination, move)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Warning("error ", commandName, " ", header.Path, ": ", err.Error())
	}
	err = sendPacket(conn, cmd+1, metadata, content)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Error("error while response to ", commandName, ": ", err.Error())
	}
}
