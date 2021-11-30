package main

import (
	"encoding/binary"
	"log"
	"net"

	"gopkg.in/yaml.v3"
)

func SendPacket(conn net.Conn, command int, metadata HeaderValues, body []byte) error {
	// version
	version := make([]byte, 2)
	binary.BigEndian.PutUint16(version, ProtocolVersion)

	// header encode
	metadataBytes, err := yaml.Marshal(&metadata)
	if err != nil {
		return err
	}

	// header length
	headerLength := 4 + len(metadataBytes) + 32
	headerLengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(headerLengthBytes, uint32(headerLength))

	// command
	commandBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(commandBytes, uint32(command))

	// checksum
	checksum := make([]byte, 32)
	for i := 0; i < 32; i++ {
		checksum[i] = 0
	}

	// body length send
	bodyLength := len(body) + 32
	bodyLengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bodyLengthBytes, uint32(bodyLength))

	// send everything
	result := Concat(version, headerLengthBytes, commandBytes, metadataBytes, checksum, bodyLengthBytes, body)
	conn.Write(result)

	log.Println("sent packet to", conn.RemoteAddr().String())

	return nil
}
