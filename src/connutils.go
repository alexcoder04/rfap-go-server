package main

import (
	"net"
)

func CleanErrorDisconnect(conn net.Conn) {
	header := HeaderMetadata{}
	err := SendPacket(conn, CMD_ERROR, header, make([]byte, 0))
	if err != nil {
		logger.Error(conn.RemoteAddr().String(), " send disconnect packet failed")
	}
	conn.Close()
	logger.Info(conn.RemoteAddr().String(), " connection closed")
}
