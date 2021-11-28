package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"

	"gopkg.in/yaml.v3"
)

func SendPacket(conn net.Conn, command int, metadata HeaderValues, body []byte) error {
	// version
	version := make([]byte, 2)
	binary.BigEndian.PutUint16(version, ProtocolVersion)
	conn.Write(version)

	// header encode
	metadataBytes, err := yaml.Marshal(metadata)
	if err != nil {
		return err
	}

	// header length send
	headerLength := uint32(4 + len(metadataBytes) + 32)
	headerLengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(headerLengthBytes, headerLength)
	conn.Write(headerLengthBytes)

	// command send
	commandBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(commandBytes, uint32(command))
	conn.Write(commandBytes)

	// header send
	conn.Write(metadataBytes)

	// checksum send
	checksum := make([]byte, 32)
	for i := 0; i < 32; i++ {
		checksum[i] = 0
	}
	conn.Write(checksum)

	// body length send
	bodyLength := len(body) + 32
	fmt.Println("send body length", bodyLength)
	bodyLengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bodyLengthBytes, uint32(bodyLength))
	conn.Write(bodyLengthBytes)

	// body send
	conn.Write(body)
	conn.Write(checksum)
	log.Println("sent packet to", conn.RemoteAddr().String())

	return nil
}
