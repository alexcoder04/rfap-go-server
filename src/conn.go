package main

import (
	"net"
)

func HanleConnection(conn net.Conn) {
	version, command, header, body, err := RecvPacket(conn)
	_ = body
	if err != nil {
		if _, ok := err.(*ErrUnsupportedRfapVersion); ok {
			logger.Error(conn.RemoteAddr().String(), " rfap version ", version, " unsupported ")
			CleanErrorDisconnect(conn)
			return
		}
		logger.Error(conn.RemoteAddr().String(), " error recieving packet: ", err.Error())
		CleanErrorDisconnect(conn)
		return
	}

	switch command {

	// server commands
	case CMD_PING:
		logger.Info(conn.RemoteAddr().String(), " packet: ping")
		err := SendPacket(conn, CMD_PING+1, HeaderMetadata{}, make([]byte, 0))
		if err != nil {
			logger.Error(conn.RemoteAddr().String(), " error while response to ping: ", err.Error())
		}
		break

	case CMD_DISCONNECT:
		logger.Info(conn.RemoteAddr().String(), " packet: disconnect")
		err := SendPacket(conn, CMD_DISCONNECT+1, HeaderMetadata{}, make([]byte, 0))
		if err != nil {
			logger.Error(conn.RemoteAddr().String(), " error while response to disconnect: ", err.Error())
		}
		conn.Close()
		logger.Info(conn.RemoteAddr().String(), " connection closed")
		return

	case CMD_INFO:
		logger.Info(conn.RemoteAddr().String(), " packet: info on ", header.Path)
		data := Info(header.Path, header.RequestDetails)
		err := SendPacket(conn, CMD_INFO+1, data, make([]byte, 0))
		if err != nil {
			logger.Error(conn.RemoteAddr().String(), "error while response to info: ", err.Error())
		}
		break

	case CMD_ERROR:
		logger.Warning(conn.LocalAddr().String(), " packet: error ", header.ErrorCode)
		break

	// file commands
	case CMD_FILE_READ:
		logger.Info(conn.RemoteAddr().String(), " packet: read file ", header.Path)
		metadata, content, err := ReadFile(header.Path)
		if err != nil {
			logger.Warning(conn.RemoteAddr().String(), " error reading file ", header.Path, ": ", err.Error())
		}
		err = SendPacket(conn, CMD_FILE_READ+1, metadata, content)
		if err != nil {
			logger.Error(conn.RemoteAddr().String(), " error while response to file_read: ", err.Error())
		}
		break

	// TODO optional file commands

	// directory commands
	case CMD_DIRECTORY_READ:
		logger.Info(conn.RemoteAddr().String(), " packet: read directory ", header.Path)
		metadata, content, err := ReadDirectory(header.Path, header.RequestDetails)
		if err != nil {
			logger.Warning(conn.RemoteAddr().String(), " error reading dir ", header.Path, ": ", err.Error())
		}
		err = SendPacket(conn, CMD_DIRECTORY_READ+1, metadata, content)
		if err != nil {
			logger.Error(conn.RemoteAddr().String(), " error while response to directory_read: ", err.Error())
		}
		break
	// TODO optional directory commands

	// unknown command
	default:
		logger.Warning(conn.RemoteAddr().String(), " packet: unknown command")
		metadata := HeaderMetadata{}
		metadata.ErrorCode = ERROR_INVALID_COMMAND
		metadata.ErrorMessage = "Unknown command requested"
		err := SendPacket(conn, CMD_ERROR+1, metadata, make([]byte, 0))
		if err != nil {
			logger.Error(conn.RemoteAddr().String(), " error while response to unknown packet: ", err.Error())
		}
		break
	}

	HanleConnection(conn)
}
