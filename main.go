package main

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	connHost = "localhost"
	connPort = "3333"
	connType = "tcp"
)

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

func getVersion(conn net.Conn) ([]byte, error) {
	buffer := make([]byte, 2)
	_, err := conn.Read(buffer[:])
	if err != nil {
		conn.Close()
		return buffer, errors.New("cannot read version")
	}
	return buffer, nil
}

func getHeaderLength(conn net.Conn) (int, error) {
	buffer := make([]byte, 4)
	_, err := conn.Read(buffer[:])
	if err != nil {
		return 0, errors.New("cannot read header length")
	}
	return int(binary.BigEndian.Uint32(buffer)), nil
}

func hanleConnection(conn net.Conn) {
	version, err1 := getVersion(conn)
	if err1 != nil {
		conn.Close()
		log.Println(err1.Error())
		return
	}
	log.Println("version: 0x" + hex.EncodeToString(version))

	headerLength, err2 := getHeaderLength(conn)
	if err2 != nil {
		conn.Close()
		log.Println(err2.Error())
		return
	}
	log.Println("header length: ", headerLength)
}
