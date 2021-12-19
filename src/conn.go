package main

import (
	"net"
	"runtime"

	"github.com/sirupsen/logrus"
)

func HanleConnection(conn net.Conn) {
	version, command, header, body, err := RecvPacket(conn)
	_ = body
	if err != nil {
		if _, ok := err.(*ErrUnsupportedRfapVersion); ok {
			logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Error("rfap version ", version, " unsupported ")
			CleanErrorDisconnect(conn)
			return
		}
		if _, ok := err.(*ErrClientCrashed); ok {
			logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Error("client crashed")
			CleanErrorDisconnect(conn)
			return
		}
		logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Error("error recieving packet: ", err.Error())
		CleanErrorDisconnect(conn)
		return
	}

	switch command {

	// server commands
	case CMD_PING:
		logger.WithFields(logrus.Fields{
			"client":  conn.RemoteAddr().String(),
			"command": "ping",
		}).Info("packet: ping")
		err := SendPacket(conn, CMD_PING+1, HeaderMetadata{}, make([]byte, 0))
		if err != nil {
			logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Error("error while response to ping: ", err.Error())
		}
		break

	case CMD_DISCONNECT:
		logger.WithFields(logrus.Fields{
			"client":  conn.RemoteAddr().String(),
			"command": "disconnect",
		}).Info("packet: disconnect")
		err := SendPacket(conn, CMD_DISCONNECT+1, HeaderMetadata{}, make([]byte, 0))
		if err != nil {
			logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Error("error while response to disconnect: ", err.Error())
		}
		conn.Close()
		logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Info("connection closed")
		logger.Info("running threads: ", runtime.NumGoroutine(), "/", MAX_CLIENTS)
		return

	case CMD_INFO:
		logger.WithFields(logrus.Fields{
			"client":  conn.RemoteAddr().String(),
			"command": "info",
		}).Info("packet: info on ", header.Path)
		data, body, err := Info(header.Path, header.RequestDetails)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Warning("error info on file ", header.Path, ": ", err.Error())
		}
		err = SendPacket(conn, CMD_INFO+1, data, body)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Error("error while response to info: ", err.Error())
		}
		break

	case CMD_ERROR:
		logger.WithFields(logrus.Fields{
			"client":  conn.RemoteAddr().String(),
			"command": "error",
		}).Warning("packet: error ", header.ErrorCode)
		break

	// file commands
	case CMD_FILE_READ:
		RunCommand(conn, header, CMD_FILE_READ, "file_read", ReadFile)
		break

	case CMD_FILE_DELETE:
		RunCommand(conn, header, CMD_FILE_DELETE, "file_delete", DeleteFile)
		break

	case CMD_FILE_CREATE:
		RunCommand(conn, header, CMD_FILE_CREATE, "file_create", CreateFile)
		break

	case CMD_FILE_COPY:
		logger.WithFields(logrus.Fields{
			"client":  conn.RemoteAddr().String(),
			"command": "file_copy",
		}).Info("packet: copy file ", header.Path, " to ", header.Destination)
		metadata, err := CopyFile(header.Path, header.Destination, false)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Warning("error copying file ", header.Path, " to ", header.Destination, ": ", err.Error())
		}
		err = SendPacket(conn, CMD_FILE_COPY+1, metadata, make([]byte, 0))
		if err != nil {
			logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Error("error while response to file_copy: ", err.Error())
		}
		break

	case CMD_FILE_MOVE:
		logger.WithFields(logrus.Fields{
			"client":  conn.RemoteAddr().String(),
			"command": "file_move",
		}).Info("packet: move file ", header.Path, " to ", header.Destination)
		metadata, err := CopyFile(header.Path, header.Destination, true)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Warning("error moving file ", header.Path, " to ", header.Destination, ": ", err.Error())
		}
		err = SendPacket(conn, CMD_FILE_MOVE+1, metadata, make([]byte, 0))
		if err != nil {
			logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Error("error while response to file_move: ", err.Error())
		}
		break

	// directory commands
	case CMD_DIRECTORY_READ:
		logger.WithFields(logrus.Fields{
			"client":  conn.RemoteAddr().String(),
			"command": "directory_read",
		}).Info("packet: read directory ", header.Path)
		metadata, content, err := ReadDirectory(header.Path, header.RequestDetails)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Warning("error reading dir ", header.Path, ": ", err.Error())
		}
		err = SendPacket(conn, CMD_DIRECTORY_READ+1, metadata, content)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Error("error while response to directory_read: ", err.Error())
		}
		break

	// unknown command
	default:
		logger.WithFields(logrus.Fields{
			"client":  conn.RemoteAddr().String(),
			"command": "unknown",
		}).Warning("packet: unknown command")
		metadata := HeaderMetadata{}
		metadata.ErrorCode = ERROR_INVALID_COMMAND
		metadata.ErrorMessage = "Unknown command requested"
		err := SendPacket(conn, CMD_ERROR+1, metadata, make([]byte, 0))
		if err != nil {
			logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Error("error while response to unknown packet: ", err.Error())
		}
		break
	}

	HanleConnection(conn)
}
