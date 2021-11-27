package main

import (
	"encoding/binary"
	"errors"
	"net"
)

func GetVersion(conn net.Conn) ([]byte, error) {
	buffer := make([]byte, 2)
	_, err := conn.Read(buffer[:])
	if err != nil {
		return buffer, errors.New("cannot read version")
	}
	return buffer, nil
}

func GetContentLength(conn net.Conn) (int, error) {
	buffer := make([]byte, 4)
	_, err := conn.Read(buffer[:])
	if err != nil {
		return 0, errors.New("cannot read content length")
	}
	return int(binary.BigEndian.Uint32(buffer)), nil
}

func GetHeader(conn net.Conn, length int) ([]byte, string, []byte, error) {
	buffer := make([]byte, length)
	_, err := conn.Read(buffer[:])
	if err != nil {
		return buffer, "", buffer, errors.New("cannot read header")
	}
	command := buffer[:4]
	header := string(buffer[4:(length - 32)])
	checksum := buffer[(length - 32):]
	return command, header, checksum, nil
}

func GetBody(conn net.Conn, length int) ([]byte, error) {
	buffer := make([]byte, length)
	_, err := conn.Read(buffer[:])
	if err != nil {
		return buffer, errors.New("cannot read header")
	}
	return buffer, nil
}
