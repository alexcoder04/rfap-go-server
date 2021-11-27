package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	connHost        = "localhost"
	connPort        = "3333"
	connType        = "tcp"
	publicFolder    = "/tmp/rfap-share"
	protocolVersion = 1
)

type HeaderValues struct {
	FilePath string `yaml:"FilePath"`
}

func main() {
	fmt.Println("Starting " + connType + " server on " + connHost + ":" + connPort)

	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}
		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

		go hanleConnection(c)
	}
}

func hanleConnection(conn net.Conn) {
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

	commandInt := binary.BigEndian.Uint32(command)
	switch commandInt {
	case 0:
		log.Println(conn.RemoteAddr().String(), "just ping")
	case 1:
		h := HeaderValues{}
		err := yaml.Unmarshal([]byte(header), &h)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		log.Println(conn.RemoteAddr().String(), "wants to read", h.FilePath)
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
	default:
		log.Println(conn.RemoteAddr().String(), "unknown command")
	}

	conn.Close()
}
