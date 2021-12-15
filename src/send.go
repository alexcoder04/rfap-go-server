package main

import (
	"encoding/binary"
	"encoding/hex"
	"net"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func SendPacket(conn net.Conn, command int, metadata HeaderMetadata, body []byte) error {
	// version
	version := make([]byte, VERSION_LENGTH)
	binary.BigEndian.PutUint16(version, RFAP_VERSION)

	// command
	commandBytes := make([]byte, COMMAND_LENGTH)
	binary.BigEndian.PutUint32(commandBytes, uint32(command))

	// header encode
	metadataBytes, err := yaml.Marshal(&metadata)
	if err != nil {
		return err
	}
	logger.WithFields(logrus.Fields{
		"client": conn.RemoteAddr().String(),
	}).Trace("header length: ", len(metadataBytes))

	// header length
	headerLength := uint32(COMMAND_LENGTH + len(metadataBytes) + CHECKSUM_LENGTH)
	headerLengthBytes := make([]byte, CONT_LEN_INDIC_LENGTH)
	binary.BigEndian.PutUint32(headerLengthBytes, headerLength)

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
	logger.WithFields(logrus.Fields{
		"client": conn.RemoteAddr().String(),
	}).Trace("packet length: ", len(result))
	// TODO send bit by bit
	_, err = conn.Write(result)
	if err != nil {
		return err
	}

	logger.WithFields(logrus.Fields{
		"client": conn.RemoteAddr().String(),
	}).Info("sent packet 0x", hex.EncodeToString(commandBytes))
	return nil
}
