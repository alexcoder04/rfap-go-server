package main

import (
	"log"
	"net"
	"os"
)

func HanleConnection(conn net.Conn) {
	version, command, header, body, err := RecvPacket(conn)
	_, _ = version, body
	if err != nil {
		log.Println("error recieving packet from", conn.RemoteAddr().String())
		log.Println(err.Error())
		conn.Close()
		log.Println("closed connection to", conn.RemoteAddr().String())
		return
	}

	switch command {

	// server commands
	case CMD_PING:
		log.Println(conn.RemoteAddr().String(), "just ping")
		SendPacket(conn, CMD_PING, HeaderMetadata{}, make([]byte, 0))
		break

	case CMD_DISCONNECT:
		log.Println(conn.RemoteAddr().String(), "wants to disconnect")
		conn.Close()
		return

	case CMD_INFO:
		data, err := Info(header.Path)
		log.Println(conn.RemoteAddr().String(), "wants info on", header.Path)
		if err != nil {
			log.Println(err.Error())
			SendPacket(conn, CMD_INFO+1, data, make([]byte, 0))
			break
		}
		SendPacket(conn, CMD_INFO+1, data, make([]byte, 0))
		break

	case CMD_ERROR:
		// TODO what if the client sends us an error?
		break

	// file commands
	case CMD_FILE_READ:
		log.Println(conn.RemoteAddr().String(), "wants to read file", header.Path)
		metadata := HeaderMetadata{}
		// TODO set this stuff afterwards
		metadata.Path = header.Path
		metadata.Type = "f"
		content, err := ReadFile(header.Path)
		if err != nil {
			if os.IsNotExist(err) {
				metadata.ErrorCode = ERROR_FILE_NOT_EXISTS
				metadata.ErrorMessage = "File does not exist"
			} else {
				metadata.ErrorCode = ERROR_UNKNOWN
				metadata.ErrorMessage = "Unknown error while reading file"
			}
			SendPacket(conn, CMD_FILE_READ+1, metadata, make([]byte, 0))
			break
		}
		SendPacket(conn, CMD_FILE_READ+1, metadata, content)
		break

	// TODO optional file commands

	// directory commands
	case CMD_DIRECTORY_READ:
		log.Println(conn.RemoteAddr().String(), "wants to read directory", header.Path)
		metadata := HeaderMetadata{}
		// TODO set this stuff afterwards
		metadata.Path = header.Path
		metadata.Type = "d"
		content, err := ReadDirectory(header.Path)
		if err != nil {
			if os.IsNotExist(err) {
				metadata.ErrorCode = ERROR_FILE_NOT_EXISTS
				metadata.ErrorMessage = "Directory does not exist"
			} else {
				metadata.ErrorCode = ERROR_UNKNOWN
				metadata.ErrorMessage = "Unknown error while reading directory"
			}
			SendPacket(conn, CMD_DIRECTORY_READ+1, metadata, make([]byte, 0))
			break
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
