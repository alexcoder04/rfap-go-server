package main

import (
	"encoding/binary"
	"net"

	"gopkg.in/yaml.v3"
)

func SendPacket(conn net.Conn, command int, metadata HeaderValues, body []byte) error {
	version := make([]byte, 2)
	binary.BigEndian.PutUint16(version, ProtocolVersion)
	conn.Write(version)

	metadataBytes, err := yaml.Marshal(metadata)
	if err != nil {
		return err
	}

	headerLength := uint32(4 + len(metadataBytes) + 32)
	headerLengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(headerLengthBytes, headerLength)
	conn.Write(headerLengthBytes)

	commandBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(commandBytes, uint32(command))
	conn.Write(commandBytes)

	conn.Write(metadataBytes)

	checksum := make([]byte, 32)
	for i := 0; i < 32; i++ {
		checksum[i] = 0
	}
	conn.Write(checksum)

	bodyLength := len(body)
	bodyLengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bodyLengthBytes, uint32(bodyLength))
	conn.Write(bodyLengthBytes)

	conn.Write(body)

	return nil
}
