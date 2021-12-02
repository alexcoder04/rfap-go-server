package main

import (
	"log"
	"net"
)

func HanleConnection(conn net.Conn) {
	version, command, header, body, err := RecvPacket(conn)
	_ = body
	if err != nil {
		if _, ok := err.(*ErrUnsupportedRfapVersion); ok {
			log.Println(conn.RemoteAddr().String(), "rfap version", version, "unsupported")
			CleanErrorDisconnect(conn)
			return
		}
		log.Println("error recieving packet from", conn.RemoteAddr().String(), err.Error())
		CleanErrorDisconnect(conn)
		return
	}

	switch command {

	// server commands
	case CMD_PING:
		log.Println(conn.RemoteAddr().String(), "just ping")
		SendPacket(conn, CMD_PING+1, HeaderMetadata{}, make([]byte, 0))
		break

	case CMD_DISCONNECT:
		log.Println(conn.RemoteAddr().String(), "wants to disconnect")
		SendPacket(conn, CMD_DISCONNECT+1, HeaderMetadata{}, make([]byte, 0))
		conn.Close()
		log.Println(conn.RemoteAddr().String() + ": connection closed")
		return

	case CMD_INFO:
		log.Println(conn.RemoteAddr().String(), "wants info on", header.Path)
		data := Info(header.Path, header.RequestDetails)
		SendPacket(conn, CMD_INFO+1, data, make([]byte, 0))
		break

	case CMD_ERROR:
		log.Println(conn.LocalAddr().String(), "sent error code", header.ErrorCode)
		break

	// file commands
	case CMD_FILE_READ:
		log.Println(conn.RemoteAddr().String(), "wants to read file", header.Path)
		metadata, content, err := ReadFile(header.Path)
		if err != nil {
			log.Println("error reading file", header.Path, err.Error())
		}
		SendPacket(conn, CMD_FILE_READ+1, metadata, content)
		break

	// TODO optional file commands

	// directory commands
	case CMD_DIRECTORY_READ:
		log.Println(conn.RemoteAddr().String(), "wants to read directory", header.Path)
		metadata, content, err := ReadDirectory(header.Path, header.RequestDetails)
		if err != nil {
			log.Println("error reading dir", header.Path, err.Error())
		}
		SendPacket(conn, CMD_DIRECTORY_READ+1, metadata, content)
		break
	// TODO optional directory commands

	// unknown command
	default:
		log.Println(conn.RemoteAddr().String(), "unknown command")
		metadata := HeaderMetadata{}
		metadata.ErrorCode = ERROR_INVALID_COMMAND
		metadata.ErrorMessage = "Unknown command requested"
		SendPacket(conn, CMD_ERROR+1, metadata, make([]byte, 0))
		break
	}

	HanleConnection(conn)
}
