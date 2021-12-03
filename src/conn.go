package main

import (
	"net"

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
			"client": conn.RemoteAddr().String(),
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
			"client": conn.RemoteAddr().String(),
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
		return

	case CMD_INFO:
		logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Info("packet: info on ", header.Path)
		data := Info(header.Path, header.RequestDetails)
		err := SendPacket(conn, CMD_INFO+1, data, make([]byte, 0))
		if err != nil {
			logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Error("error while response to info: ", err.Error())
		}
		break

	case CMD_ERROR:
		logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Warning("packet: error ", header.ErrorCode)
		break

	// file commands
	case CMD_FILE_READ:
		logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Info("packet: read file ", header.Path)
		metadata, content, err := ReadFile(header.Path)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Warning("error reading file ", header.Path, ": ", err.Error())
		}
		err = SendPacket(conn, CMD_FILE_READ+1, metadata, content)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Error("error while response to file_read: ", err.Error())
		}
		break

	// TODO optional file commands

	// directory commands
	case CMD_DIRECTORY_READ:
		logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
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
	// TODO optional directory commands

	// unknown command
	default:
		logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
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
