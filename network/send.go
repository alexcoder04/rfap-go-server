package network

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"net"

	"github.com/alexcoder04/rfap-go-server/log"
	"github.com/alexcoder04/rfap-go-server/settings"
	"github.com/alexcoder04/rfap-go-server/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func sendPacket(conn net.Conn, command uint32, metadata utils.HeaderMetadata, body []byte) error {
	// version
	version := make([]byte, settings.VERSION_LENGTH)
	binary.BigEndian.PutUint16(version, settings.RFAP_VERSION)

	// command
	commandBytes := make([]byte, settings.COMMAND_LENGTH)
	binary.BigEndian.PutUint32(commandBytes, uint32(command))

	// header encode
	metadataBytes, err := yaml.Marshal(&metadata)
	if err != nil {
		return err
	}
	log.Logger.WithFields(logrus.Fields{
		"client": conn.RemoteAddr().String(),
	}).Trace("header length: ", len(metadataBytes))
	if len(metadataBytes) > (1024 * 8) {
		return &utils.ErrInvalidContentLength{}
	}

	// header length
	headerLength := uint32(settings.COMMAND_LENGTH + len(metadataBytes) + settings.CHECKSUM_LENGTH)
	headerLengthBytes := make([]byte, settings.CONT_LEN_INDIC_LENGTH)
	binary.BigEndian.PutUint32(headerLengthBytes, headerLength)

	// checksum
	headerChecksum := sha256.Sum256(utils.ConcatBytes(commandBytes, metadataBytes))

	// send header
	firstPart := utils.ConcatBytes(version, headerLengthBytes, commandBytes, metadataBytes, headerChecksum[:])
	_, err = conn.Write(firstPart)
	if err != nil {
		return err
	}

	// body length
	bodyLength := uint32(len(body) + settings.CHECKSUM_LENGTH)
	bodyLengthBytes := make([]byte, settings.CONT_LEN_INDIC_LENGTH)
	binary.BigEndian.PutUint32(bodyLengthBytes, bodyLength)
	log.Logger.WithFields(logrus.Fields{
		"client": conn.RemoteAddr().String(),
	}).Trace("body length: ", bodyLength)

	// send body length and body
	_, err = conn.Write(bodyLengthBytes)
	if err != nil {
		return err
	}
	i := 0
	for {
		if (i + settings.MAX_BYTES_SEND_AT_ONCE) > len(body) {
			_, err := conn.Write(body[i:])
			if err != nil {
				return err
			}
			break
		}
		_, err := conn.Write(body[i : i+settings.MAX_BYTES_SEND_AT_ONCE])
		if err != nil {
			return err
		}
		i += settings.MAX_BYTES_SEND_AT_ONCE
	}
	bodyChecksum := sha256.Sum256(body)
	_, err = conn.Write(bodyChecksum[:])
	if err != nil {
		return err
	}

	log.Logger.WithFields(logrus.Fields{
		"client": conn.RemoteAddr().String(),
	}).Info("sent packet 0x", hex.EncodeToString(commandBytes))
	return nil
}
