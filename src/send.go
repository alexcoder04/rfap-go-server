package main

import (
	"encoding/binary"
	"encoding/hex"
	"net"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func SendPacket(conn net.Conn, command int, metadata HeaderMetadata, body []byte) error {
	if metadata.PacketsTotal == 0 {
		metadata.PacketsTotal = 1
	}
	if metadata.PacketNumber == 0 {
		metadata.PacketNumber = 1
	}
	// version
	version := make([]byte, VERSION_LENGTH)
	binary.BigEndian.PutUint16(version, RFAP_VERSION)

	// header encode
	metadataBytes, err := yaml.Marshal(&metadata)
	if err != nil {
		return err
	}

	// header length
	headerLength := uint32(COMMAND_LENGTH + len(metadataBytes) + CHECKSUM_LENGTH)
	headerLengthBytes := make([]byte, CONT_LEN_INDIC_LENGTH)
	binary.BigEndian.PutUint32(headerLengthBytes, headerLength)

	// command
	commandBytes := make([]byte, COMMAND_LENGTH)
	binary.BigEndian.PutUint32(commandBytes, uint32(command))

	// checksum
	checksum := make([]byte, CHECKSUM_LENGTH)
	for i := 0; i < 32; i++ {
		checksum[i] = 0
	}

	// body length send
	bodyLength := uint32(len(body) + CHECKSUM_LENGTH)
	bodyLengthBytes := make([]byte, CONT_LEN_INDIC_LENGTH)
	binary.BigEndian.PutUint32(bodyLengthBytes, bodyLength)
	logger.WithFields(logrus.Fields{
		"client": conn.RemoteAddr().String(),
	}).Trace("body length: ", bodyLength)

	// send everything
	result := ConcatBytes(version, headerLengthBytes, commandBytes, metadataBytes, checksum, bodyLengthBytes, body, checksum)
	_, err = conn.Write(result)
	if err != nil {
		return err
	}

	logger.WithFields(logrus.Fields{
		"client": conn.RemoteAddr().String(),
	}).Info("sent packet 0x", hex.EncodeToString(commandBytes))
	return nil
}

// same as SendPacket, but can handle bodies larger then max length (splitting)
func SendData(conn net.Conn, command int, metadata HeaderMetadata, body []byte) error {
	maxBodyLen := (MAX_CONT_LENGTH_MB * 1024 * 1024) - 32
	if len(body) <= maxBodyLen {
		metadata.PacketNumber = 1
		metadata.PacketsTotal = 1
		return SendPacket(conn, command, metadata, body)
	}

	var chunk []byte
	chunks := make([][]byte, 0, len(body)/(maxBodyLen+1))
	for len(body) >= maxBodyLen {
		chunk, body = body[:maxBodyLen], body[maxBodyLen:]
		chunks = append(chunks, chunk)
	}
	if len(body) > 0 {
		chunks = append(chunks, body[:])
	}

	metadata.PacketsTotal = len(chunks)
	for i := 0; i < len(chunks); i++ {
		metadata.PacketNumber = i + 1
		logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Debug("Sending packet ", metadata.PacketNumber, " of ", metadata.PacketsTotal, len(chunks[i]), "...")
		err := SendPacket(conn, command, metadata, chunks[i])
		if err != nil {
			return err
		}
	}

	return nil
}
