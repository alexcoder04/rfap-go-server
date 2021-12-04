package main

import (
	"encoding/binary"
	"encoding/hex"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func RecvPacket(conn net.Conn) (uint32, uint32, HeaderMetadata, []byte, error) {
	// receive
	buffer := make([]byte, VERSION_LENGTH+CONT_LEN_INDIC_LENGTH+(MAX_CONT_LENGTH_MB*1024*1024)+CONT_LEN_INDIC_LENGTH+(MAX_CONT_LENGTH_MB*1024*1024))
	err := conn.SetReadDeadline(time.Now().Add(CONN_RECV_TIMEOUT_SECS * time.Second))
	if err != nil {
		return 0, 0, HeaderMetadata{}, make([]byte, 0), &ErrSetReadTimeoutFailed{}
	}
	_, err = conn.Read(buffer[:])
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return 0, 0, HeaderMetadata{}, make([]byte, 0), &ErrClientCrashed{}
		}
		return 0, 0, HeaderMetadata{}, make([]byte, 0), err
	}

	// sort and split
	version := uint32(binary.BigEndian.Uint16(buffer[:VERSION_LENGTH]))
	logger.Debug("version:", version)
	if !Uint32ArrayContains(SUPPORTED_RFAP_VERSIONS, version) {
		return version, 0, HeaderMetadata{}, make([]byte, 0), &ErrUnsupportedRfapVersion{}
	}

	headerLength := binary.BigEndian.Uint32(buffer[VERSION_LENGTH : VERSION_LENGTH+CONT_LEN_INDIC_LENGTH])
	logger.Debug("header length:", headerLength)

	command := binary.BigEndian.Uint32(buffer[VERSION_LENGTH+CONT_LEN_INDIC_LENGTH : VERSION_LENGTH+CONT_LEN_INDIC_LENGTH+COMMAND_LENGTH])
	logger.Debug("command: 0x" + hex.EncodeToString(buffer[VERSION_LENGTH+CONT_LEN_INDIC_LENGTH:VERSION_LENGTH+CONT_LEN_INDIC_LENGTH+COMMAND_LENGTH]))

	headerRaw := buffer[VERSION_LENGTH+CONT_LEN_INDIC_LENGTH+COMMAND_LENGTH : VERSION_LENGTH+CONT_LEN_INDIC_LENGTH+(headerLength-CHECKSUM_LENGTH)]
	logger.Debug("header:", hex.EncodeToString(headerRaw))

	headerChecksum := buffer[VERSION_LENGTH+CONT_LEN_INDIC_LENGTH+(headerLength-CHECKSUM_LENGTH) : VERSION_LENGTH+CONT_LEN_INDIC_LENGTH+headerLength]
	_ = headerChecksum
	logger.Debug("header checksum:", hex.EncodeToString(headerChecksum))

	bodyLength := binary.BigEndian.Uint32(buffer[VERSION_LENGTH+CONT_LEN_INDIC_LENGTH+headerLength : VERSION_LENGTH+CONT_LEN_INDIC_LENGTH+headerLength+CONT_LEN_INDIC_LENGTH])
	logger.Debug("body length:", bodyLength)

	body := buffer[VERSION_LENGTH+CONT_LEN_INDIC_LENGTH+headerLength+CONT_LEN_INDIC_LENGTH : VERSION_LENGTH+CONT_LEN_INDIC_LENGTH+headerLength+CONT_LEN_INDIC_LENGTH+(bodyLength-CHECKSUM_LENGTH)]
	logger.Debug("body:", hex.EncodeToString(body))

	bodyChecksum := buffer[VERSION_LENGTH+CONT_LEN_INDIC_LENGTH+headerLength+CONT_LEN_INDIC_LENGTH+(bodyLength-CHECKSUM_LENGTH) : VERSION_LENGTH+CONT_LEN_INDIC_LENGTH+headerLength+CONT_LEN_INDIC_LENGTH+bodyLength]
	_ = bodyChecksum
	logger.Debug("body checksum:", hex.EncodeToString(bodyChecksum))

	// parse
	header := HeaderMetadata{}
	err = yaml.Unmarshal([]byte(headerRaw), &header)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Error("error decoding metadata")
		return version, command, HeaderMetadata{}, body, err
	}

	// return
	return version, command, header, body, nil
}
