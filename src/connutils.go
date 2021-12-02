package main

import (
	"log"
	"net"
)

func CleanErrorDisconnect(conn net.Conn) {
	header := HeaderMetadata{}
	SendPacket(conn, CMD_ERROR, header, make([]byte, 0))
	conn.Close()
	log.Println(conn.RemoteAddr().String() + ": connection closed")
}
