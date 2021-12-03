package main

import (
	"encoding/binary"
	"encoding/hex"
	"net"

	"gopkg.in/yaml.v3"
)

func SendPacket(conn net.Conn, command int, metadata HeaderMetadata, body []byte) error {
	// version
	version := make([]byte, 2)
	binary.BigEndian.PutUint16(version, RFAP_VERSION)

	// header encode
	metadataBytes, err := yaml.Marshal(&metadata)
	if err != nil {
		return err
	}

	// header length
	headerLength := uint32(4 + len(metadataBytes) + 32)
	headerLengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(headerLengthBytes, headerLength)

	// command
	commandBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(commandBytes, uint32(command))

	// checksum
	checksum := make([]byte, 32)
	for i := 0; i < 32; i++ {
		checksum[i] = 0
	}

	// body length send
	bodyLength := uint32(len(body) + 32)
	bodyLengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bodyLengthBytes, bodyLength)

	// send everything
	result := ConcatBytes(version, headerLengthBytes, commandBytes, metadataBytes, checksum, bodyLengthBytes, body, checksum)
	_, err = conn.Write(result)
	if err != nil {
		return err
	}

	logger.Info(conn.RemoteAddr().String(), " sent packet 0x", hex.EncodeToString(commandBytes))
	return nil
}
