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
