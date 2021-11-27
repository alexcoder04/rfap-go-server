package main

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"log"
	"net"

	"gopkg.in/yaml.v3"
)

func RecvPacket(conn net.Conn) (int, int, HeaderValues, []byte, error) {
	version, err1 := GetVersion(conn)
	if err1 != nil {
		conn.Close()
		log.Println(err1.Error())
		return 0, 0, HeaderValues{}, make([]byte, 0), errors.New("cannot read version")
	}
	log.Println("version: 0x" + hex.EncodeToString(version))

	headerLength, err2 := GetContentLength(conn)
	if err2 != nil {
		conn.Close()
		log.Println(err2.Error())
		return 0, 0, HeaderValues{}, make([]byte, 0), errors.New("cannot read header length")
	}
	log.Println("header length: ", headerLength)

	command, headerBytes, headerChecksum, err3 := GetHeader(conn, headerLength)
	if err3 != nil {
		conn.Close()
		log.Println(err3.Error())
		return 0, 0, HeaderValues{}, make([]byte, 0), errors.New("cannot read header")
	}
	log.Println("command: 0x" + hex.EncodeToString(command))
	log.Println("header: ", headerBytes)
	log.Println("header checksum: 0x" + hex.EncodeToString(headerChecksum))
	header := HeaderValues{}
	err := yaml.Unmarshal([]byte(headerBytes), &header)
	if err != nil {
		return 0, 0, HeaderValues{}, make([]byte, 0), err
	}

	bodyLength, err4 := GetContentLength(conn)
	if err4 != nil {
		conn.Close()
		log.Println(err4.Error())
		return 0, 0, HeaderValues{}, make([]byte, 0), errors.New("cannot read body length")
	}
	log.Println("body length: ", bodyLength)

	body, err5 := GetBody(conn, bodyLength)
	if err5 != nil {
		conn.Close()
		log.Println(err5.Error())
		return 0, 0, HeaderValues{}, make([]byte, 0), errors.New("cannot read body")
	}
	log.Println("body: 0x" + hex.EncodeToString(body))

	return int(binary.BigEndian.Uint16(version)), int(binary.BigEndian.Uint32(command)), header, body, nil
}

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
