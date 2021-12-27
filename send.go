package main

import (
	"crypto/sha256"
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
	if len(metadataBytes) > (1024 * 8) {
		return &ErrInvalidContentLength{}
	}

	// header length
	headerLength := uint32(COMMAND_LENGTH + len(metadataBytes) + CHECKSUM_LENGTH)
	headerLengthBytes := make([]byte, CONT_LEN_INDIC_LENGTH)
	binary.BigEndian.PutUint32(headerLengthBytes, headerLength)

	// checksum
	headerChecksum := sha256.Sum256(ConcatBytes(commandBytes, metadataBytes))

	// send header
	firstPart := ConcatBytes(version, headerLengthBytes, commandBytes, metadataBytes, headerChecksum[:])
	_, err = conn.Write(firstPart)
	if err != nil {
		return err
	}

	// body length
	bodyLength := uint32(len(body) + CHECKSUM_LENGTH)
	bodyLengthBytes := make([]byte, CONT_LEN_INDIC_LENGTH)
	binary.BigEndian.PutUint32(bodyLengthBytes, bodyLength)
	logger.WithFields(logrus.Fields{
		"client": conn.RemoteAddr().String(),
	}).Trace("body length: ", bodyLength)

	// send body length and body
	_, err = conn.Write(bodyLengthBytes)
	if err != nil {
		return err
	}
	i := 0
	for {
		if (i + MAX_BYTES_SEND_AT_ONCE) > len(body) {
			_, err := conn.Write(body[i:])
			if err != nil {
				return err
			}
			break
		}
		_, err := conn.Write(body[i : i+MAX_BYTES_SEND_AT_ONCE])
		if err != nil {
			return err
		}
		i += MAX_BYTES_SEND_AT_ONCE
	}
	bodyChecksum := sha256.Sum256(body)
	_, err = conn.Write(bodyChecksum[:])
	if err != nil {
		return err
	}

	logger.WithFields(logrus.Fields{
		"client": conn.RemoteAddr().String(),
	}).Info("sent packet 0x", hex.EncodeToString(commandBytes))
	return nil
}