package network

import (
	"net"
	"runtime"

	"github.com/alexcoder04/rfap-go-server/commands"
	"github.com/alexcoder04/rfap-go-server/log"
	"github.com/alexcoder04/rfap-go-server/settings"
	"github.com/alexcoder04/rfap-go-server/utils"
	"github.com/jchavannes/go-pgp/pgp"
	"github.com/sirupsen/logrus"
)

func connectionLoop(client utils.Client) {
	for {
		version, command, header, body, err := recvPacket(client.Conn)
		if err != nil {
			errRecvPacket(client.Conn, err, version)
			return
		}

		switch command {

		// server commands
		case settings.CMD_PING:
			log.Logger.WithFields(logrus.Fields{
				"client":  client.Address,
				"command": "ping",
			}).Info("packet: ping")
			err := sendPacket(client.Conn, settings.CMD_PING+1, utils.HeaderMetadata{}, make([]byte, 0))
			if err != nil {
				log.Logger.WithFields(logrus.Fields{
					"client": client.Address,
				}).Error("error while response to ping: ", err.Error())
			}
			break

		case settings.CMD_DISCONNECT:
			log.Logger.WithFields(logrus.Fields{
				"client":  client.Address,
				"command": "disconnect",
			}).Info("packet: disconnect")
			err := sendPacket(client.Conn, settings.CMD_DISCONNECT+1, utils.HeaderMetadata{}, make([]byte, 0))
			if err != nil {
				log.Logger.WithFields(logrus.Fields{
					"client": client.Address,
				}).Error("error while response to disconnect: ", err.Error())
			}
			client.Conn.Close()
			log.Logger.WithFields(logrus.Fields{
				"client": client.Address,
			}).Info("connection closed")
			log.Logger.Info("running threads: ", runtime.NumGoroutine(), "/", settings.Config.MaxClients())
			return

		case settings.CMD_INFO:
			log.Logger.WithFields(logrus.Fields{
				"client":  client.Address,
				"command": "info",
			}).Info("packet: info on ", header.Path)
			data, respBody, err := commands.Info(header.Path, header.RequestDetails)
			if err != nil {
				log.Logger.WithFields(logrus.Fields{
					"client": client.Address,
				}).Warning("error info on file ", header.Path, ": ", err.Error())
			}
			err = sendPacket(client.Conn, settings.CMD_INFO+1, data, respBody)
			if err != nil {
				log.Logger.WithFields(logrus.Fields{
					"client": client.Address,
				}).Error("error while response to info: ", err.Error())
			}
			break

		case settings.CMD_ERROR:
			log.Logger.WithFields(logrus.Fields{
				"client":  client.Address,
				"command": "error",
			}).Warning("packet: error ", header.ErrorCode)
			break

		// file commands
		case settings.CMD_FILE_READ:
			runCommand(client.Conn, header, settings.CMD_FILE_READ, "file_read", commands.ReadFile)
			break

		case settings.CMD_FILE_DELETE:
			runCommand(client.Conn, header, settings.CMD_FILE_DELETE, "file_delete", commands.DeleteFile)
			break

		case settings.CMD_FILE_CREATE:
			runCommand(client.Conn, header, settings.CMD_FILE_CREATE, "file_create", commands.CreateFile)
			break

		case settings.CMD_FILE_COPY:
			runCopyCommand(client.Conn, header, settings.CMD_FILE_COPY, "file_copy", commands.CopyFile, false)
			break

		case settings.CMD_FILE_MOVE:
			runCopyCommand(client.Conn, header, settings.CMD_FILE_MOVE, "file_move", commands.CopyFile, true)
			break

		case settings.CMD_FILE_WRITE:
			log.Logger.WithFields(logrus.Fields{
				"client":  client.Address,
				"command": "file_write",
			}).Info("packet: write file ", header.Path)
			metadata, respBody, err := commands.WriteFile(header.Path, body)
			if err != nil {
				log.Logger.WithFields(logrus.Fields{
					"client": client.Address,
				}).Warning("error writing file ", header.Path, ": ", err.Error())
			}
			err = sendPacket(client.Conn, settings.CMD_FILE_WRITE+1, metadata, respBody)
			if err != nil {
				log.Logger.WithFields(logrus.Fields{
					"client": client.Address,
				}).Error("error while response to file_write: ", err.Error())
			}
			break

		// directory commands
		case settings.CMD_DIRECTORY_READ:
			log.Logger.WithFields(logrus.Fields{
				"client":  client.Address,
				"command": "directory_read",
			}).Info("packet: read directory ", header.Path)
			metadata, content, err := commands.ReadDirectory(header.Path, header.RequestDetails)
			if err != nil {
				log.Logger.WithFields(logrus.Fields{
					"client": client.Address,
				}).Warning("error reading dir ", header.Path, ": ", err.Error())
			}
			err = sendPacket(client.Conn, settings.CMD_DIRECTORY_READ+1, metadata, content)
			if err != nil {
				log.Logger.WithFields(logrus.Fields{
					"client": client.Address,
				}).Error("error while response to directory_read: ", err.Error())
			}
			break

		case settings.CMD_DIRECTORY_DELETE:
			runCommand(client.Conn, header, settings.CMD_DIRECTORY_DELETE, "directory_delete", commands.DeleteDirectory)
			break

		case settings.CMD_DIRECTORY_CREATE:
			runCommand(client.Conn, header, settings.CMD_DIRECTORY_CREATE, "directory_create", commands.CreateDirectory)
			break

		case settings.CMD_DIRECTORY_COPY:
			runCopyCommand(client.Conn, header, settings.CMD_DIRECTORY_COPY, "directory_copy", commands.CopyDirectory, false)
			break

		case settings.CMD_DIRECTORY_MOVE:
			runCopyCommand(client.Conn, header, settings.CMD_DIRECTORY_MOVE, "directory_move", commands.CopyDirectory, true)
			break

		// unknown command
		default:
			log.Logger.WithFields(logrus.Fields{
				"client":  client.Address,
				"command": "unknown",
			}).Warning("packet: unknown command")
			metadata := utils.HeaderMetadata{}
			metadata.ErrorCode = int(settings.ERROR_INVALID_COMMAND)
			metadata.ErrorMessage = "Unknown command requested"
			err := sendPacket(client.Conn, settings.CMD_ERROR+1, metadata, make([]byte, 0))
			if err != nil {
				log.Logger.WithFields(logrus.Fields{
					"client": client.Address,
				}).Error("error while response to unknown packet: ", err.Error())
			}
			break
		}
	}
}

func HanleConnection(conn net.Conn) {
	version, command, _, body, err := recvPacket(conn)
	if err != nil {
		errRecvPacket(conn, err, version)
		return
	}

	if command != settings.CMD_PUB_KEY {
		log.Logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Warning("no pub key exchanged, disconnecting")
		metadata := utils.HeaderMetadata{}
		metadata.ErrorCode = int(settings.ERROR_NO_PUB_KEY)
		metadata.ErrorMessage = "First message must be public key exchange"
		err := sendPacket(conn, settings.CMD_ERROR+1, metadata, make([]byte, 0))
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Error("error while send pubkey error: ", err.Error())
		}
		conn.Close()
		return
	}

	entity, err := pgp.GetEntity(body, []byte{})
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Warning("cannot get pub key, disconnecting")
		metadata := utils.HeaderMetadata{}
		metadata.ErrorCode = int(settings.ERROR_NO_PUB_KEY)
		metadata.ErrorMessage = "Cannot extract public key from message body"
		err := sendPacket(conn, settings.CMD_ERROR+1, metadata, make([]byte, 0))
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Error("error while send pubkey error: ", err.Error())
		}
		conn.Close()
		return
	}
	client := utils.Client{}
	client.Conn = conn
	client.PubkeyEntity = entity
	client.Address = conn.RemoteAddr().String()

	connectionLoop(client)
}
