package main

import (
	"net"
	"runtime"

	"github.com/sirupsen/logrus"
)

func CleanErrorDisconnect(conn net.Conn) {
	header := HeaderMetadata{}
	err := SendPacket(conn, CMD_ERROR, header, make([]byte, 0))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Error("send disconnect packet failed")
	}
	conn.Close()
	logger.WithFields(logrus.Fields{
		"client": conn.RemoteAddr().String(),
	}).Info("connection closed")
	logger.Info("running threads: ", runtime.NumGoroutine(), "/", MAX_CLIENTS)
}

func RunCommand(conn net.Conn, header HeaderMetadata, cmd int, commandName string, fn commandExec) {
	logger.WithFields(logrus.Fields{
		"client":  conn.RemoteAddr().String(),
		"command": commandName,
	}).Info("packet: ", header.Path)
	metadata, content, err := fn(header.Path)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Warning("error ", commandName, " ", header.Path, ": ", err.Error())
	}
	err = SendPacket(conn, cmd+1, metadata, content)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Error("error while response to ", commandName, ": ", err.Error())
	}
}

func RunCopyCommand(conn net.Conn, header HeaderMetadata, cmd int, commandName string, fn copySommandExec, move bool) {
	logger.WithFields(logrus.Fields{
		"client":  conn.RemoteAddr().String(),
		"command": commandName,
	}).Info("packet: ", header.Path)
	metadata, content, err := fn(header.Path, header.Destination, move)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Warning("error ", commandName, " ", header.Path, ": ", err.Error())
	}
	err = SendPacket(conn, cmd+1, metadata, content)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Error("error while response to ", commandName, ": ", err.Error())
	}
}
