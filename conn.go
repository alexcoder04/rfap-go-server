package main

import (
	"net"
	"runtime"

	"github.com/alexcoder04/rfap-go-server/commands"
	"github.com/alexcoder04/rfap-go-server/log"
	"github.com/alexcoder04/rfap-go-server/settings"
	"github.com/alexcoder04/rfap-go-server/utils"
	"github.com/sirupsen/logrus"
)

func HanleConnection(conn net.Conn) {
	version, command, header, body, err := RecvPacket(conn)
	if err != nil {
		if _, ok := err.(*utils.ErrUnsupportedRfapVersion); ok {
			log.Logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Error("rfap version ", version, " unsupported ")
			CleanErrorDisconnect(conn)
			return
		}
		if _, ok := err.(*utils.ErrClientCrashed); ok {
			log.Logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Error("client crashed")
			CleanErrorDisconnect(conn)
			return
		}
		log.Logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Error("error recieving packet: ", err.Error())
		CleanErrorDisconnect(conn)
		return
	}

	switch command {

	// server commands
	case settings.CMD_PING:
		log.Logger.WithFields(logrus.Fields{
			"client":  conn.RemoteAddr().String(),
			"command": "ping",
		}).Info("packet: ping")
		err := SendPacket(conn, settings.CMD_PING+1, utils.HeaderMetadata{}, make([]byte, 0))
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Error("error while response to ping: ", err.Error())
		}
		break

	case settings.CMD_DISCONNECT:
		log.Logger.WithFields(logrus.Fields{
			"client":  conn.RemoteAddr().String(),
			"command": "disconnect",
		}).Info("packet: disconnect")
		err := SendPacket(conn, settings.CMD_DISCONNECT+1, utils.HeaderMetadata{}, make([]byte, 0))
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Error("error while response to disconnect: ", err.Error())
		}
		conn.Close()
		log.Logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Info("connection closed")
		log.Logger.Info("running threads: ", runtime.NumGoroutine(), "/", settings.MAX_CLIENTS)
		return

	case settings.CMD_INFO:
		log.Logger.WithFields(logrus.Fields{
			"client":  conn.RemoteAddr().String(),
			"command": "info",
		}).Info("packet: info on ", header.Path)
		data, respBody, err := commands.Info(header.Path, header.RequestDetails)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Warning("error info on file ", header.Path, ": ", err.Error())
		}
		err = SendPacket(conn, settings.CMD_INFO+1, data, respBody)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Error("error while response to info: ", err.Error())
		}
		break

	case settings.CMD_ERROR:
		log.Logger.WithFields(logrus.Fields{
			"client":  conn.RemoteAddr().String(),
			"command": "error",
		}).Warning("packet: error ", header.ErrorCode)
		break

	// file commands
	case settings.CMD_FILE_READ:
		RunCommand(conn, header, settings.CMD_FILE_READ, "file_read", commands.ReadFile)
		break

	case settings.CMD_FILE_DELETE:
		RunCommand(conn, header, settings.CMD_FILE_DELETE, "file_delete", commands.DeleteFile)
		break

	case settings.CMD_FILE_CREATE:
		RunCommand(conn, header, settings.CMD_FILE_CREATE, "file_create", commands.CreateFile)
		break

	case settings.CMD_FILE_COPY:
		RunCopyCommand(conn, header, settings.CMD_FILE_COPY, "file_copy", commands.CopyFile, false)
		break

	case settings.CMD_FILE_MOVE:
		RunCopyCommand(conn, header, settings.CMD_FILE_MOVE, "file_move", commands.CopyFile, true)
		break

	case settings.CMD_FILE_WRITE:
		log.Logger.WithFields(logrus.Fields{
			"client":  conn.RemoteAddr().String(),
			"command": "file_write",
		}).Info("packet: write file ", header.Path)
		metadata, respBody, err := commands.WriteFile(header.Path, body)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Warning("error writing file ", header.Path, ": ", err.Error())
		}
		err = SendPacket(conn, settings.CMD_FILE_WRITE+1, metadata, respBody)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Error("error while response to file_write: ", err.Error())
		}
		break

	// directory commands
	case settings.CMD_DIRECTORY_READ:
		log.Logger.WithFields(logrus.Fields{
			"client":  conn.RemoteAddr().String(),
			"command": "directory_read",
		}).Info("packet: read directory ", header.Path)
		metadata, content, err := commands.ReadDirectory(header.Path, header.RequestDetails)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Warning("error reading dir ", header.Path, ": ", err.Error())
		}
		err = SendPacket(conn, settings.CMD_DIRECTORY_READ+1, metadata, content)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Error("error while response to directory_read: ", err.Error())
		}
		break

	case settings.CMD_DIRECTORY_DELETE:
		RunCommand(conn, header, settings.CMD_DIRECTORY_DELETE, "directory_delete", commands.DeleteDirectory)
		break

	case settings.CMD_DIRECTORY_CREATE:
		RunCommand(conn, header, settings.CMD_DIRECTORY_CREATE, "directory_create", commands.CreateDirectory)
		break

	case settings.CMD_DIRECTORY_COPY:
		RunCopyCommand(conn, header, settings.CMD_DIRECTORY_COPY, "directory_copy", commands.CopyDirectory, false)
		break

	case settings.CMD_DIRECTORY_MOVE:
		RunCopyCommand(conn, header, settings.CMD_DIRECTORY_MOVE, "directory_move", commands.CopyDirectory, true)
		break

	// unknown command
	default:
		log.Logger.WithFields(logrus.Fields{
			"client":  conn.RemoteAddr().String(),
			"command": "unknown",
		}).Warning("packet: unknown command")
		metadata := utils.HeaderMetadata{}
		metadata.ErrorCode = settings.ERROR_INVALID_COMMAND
		metadata.ErrorMessage = "Unknown command requested"
		err := SendPacket(conn, settings.CMD_ERROR+1, metadata, make([]byte, 0))
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Error("error while response to unknown packet: ", err.Error())
		}
		break
	}

	HanleConnection(conn)
}
