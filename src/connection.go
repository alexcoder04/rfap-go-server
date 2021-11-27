package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"net"

	"gopkg.in/yaml.v3"
)

func HanleConnection(conn net.Conn) {
	// TODO separate function for receive packet
	version, err1 := GetVersion(conn)
	if err1 != nil {
		conn.Close()
		log.Println(err1.Error())
		return
	}
	log.Println("version: 0x" + hex.EncodeToString(version))

	headerLength, err2 := GetContentLength(conn)
	if err2 != nil {
		conn.Close()
		log.Println(err2.Error())
		return
	}
	log.Println("header length: ", headerLength)

	command, header, headerChecksum, err3 := GetHeader(conn, headerLength)
	if err3 != nil {
		conn.Close()
		log.Println(err3.Error())
		return
	}
	log.Println("command: 0x" + hex.EncodeToString(command))
	log.Println("header: ", header)
	log.Println("header checksum: 0x" + hex.EncodeToString(headerChecksum))

	bodyLength, err4 := GetContentLength(conn)
	if err4 != nil {
		conn.Close()
		log.Println(err4.Error())
		return
	}
	log.Println("body length: ", bodyLength)

	body, err5 := GetBody(conn, bodyLength)
	if err5 != nil {
		conn.Close()
		log.Println(err5.Error())
		return
	}
	log.Println("body: 0x" + hex.EncodeToString(body))

	switch binary.BigEndian.Uint32(command) {

	// server commands
	case CMD_PING:
		log.Println(conn.RemoteAddr().String(), "just ping")
		SendPacket(conn, CMD_PING, "", make([]byte, 0))
	case CMD_DISCONNECT:
		log.Println(conn.RemoteAddr().String(), "wants to disconnect")
		conn.Close()
		return
	case CMD_INFO:
		log.Println(conn.RemoteAddr().String(), "wants info")
		// TODO
		return

	// file commands
	case CMD_FILE_READ:
		h := HeaderValues{}
		// TODO header decoding in separate function in separate file
		err := yaml.Unmarshal([]byte(header), &h)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		log.Println(conn.RemoteAddr().String(), "wants to read", h.FilePath)
		// TODO FileRead not generic Read
		fileContent, fileReadErr := Read(h.FilePath)
		if fileReadErr != nil {
			log.Println(fileReadErr.Error())
			return
		}
		fmt.Println(string(fileContent))
		sendErr := SendPacket(conn, 5, "FilePath: "+h.FilePath, fileContent)
		if sendErr != nil {
			log.Println(sendErr.Error())
			return
		}
	// TODO optional file commands

	// directory commands
	case CMD_DIRECTORY_READ:
		// TODO
		return
	// TODO optional directory commands

	// unknown command
	default:
		log.Println(conn.RemoteAddr().String(), "unknown command")
	}

	// TODO not close but listen
	conn.Close()
}
