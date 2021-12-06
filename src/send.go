package main

import (
	"encoding/binary"
	"encoding/hex"
	"net"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func SendPacket(conn net.Conn, command int, metadata HeaderMetadata, body []byte) error {
	// TODO refactor this into separate function
	// split if necessary
	if len(body) > MAX_CONT_LENGTH_MB*1024*1024 {
		for i := 0; i < len(body); i += MAX_CONT_LENGTH_MB * 1024 * 1024 {
			logger.WithFields(logrus.Fields{
				"client": conn.RemoteAddr().String(),
			}).Debug("Sending portion packet...")
			err := SendPacket(conn, command, metadata, body[i:MAX_CONT_LENGTH_MB*1024*1024])
			if err != nil {
				return err
			}
		}
		return nil
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
