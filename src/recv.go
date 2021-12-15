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
	// TODO recv bit by bit
	// receive
	buffer := make([]byte, MAX_PACKET_LENGTH)
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

	headerLengthBegin := VERSION_LENGTH
	headerLengthEnd := headerLengthBegin + CONT_LEN_INDIC_LENGTH
	headerLength := binary.BigEndian.Uint32(buffer[headerLengthBegin:headerLengthEnd])
	logger.Debug("header length:", headerLength)

	commandBegin := VERSION_LENGTH + CONT_LEN_INDIC_LENGTH
	commandEnd := commandBegin + COMMAND_LENGTH
	command := binary.BigEndian.Uint32(buffer[commandBegin:commandEnd])
	logger.Debug("command: 0x" + hex.EncodeToString(buffer[commandBegin:commandEnd]))

	headerBegin := headerLengthEnd + COMMAND_LENGTH
	headerEnd := headerLengthEnd + (int(headerLength) - CHECKSUM_LENGTH)
	headerRaw := buffer[headerBegin:headerEnd]
	logger.Debug("header:", hex.EncodeToString(headerRaw))

	headerChecksumBegin := headerEnd
	headerChecksumEnd := headerChecksumBegin + CHECKSUM_LENGTH
	headerChecksum := buffer[headerChecksumBegin:headerChecksumEnd]
	_ = headerChecksum
	logger.Debug("header checksum:", hex.EncodeToString(headerChecksum))

	bodyLengthBegin := headerChecksumEnd
	bodyLengthEnd := bodyLengthBegin + CONT_LEN_INDIC_LENGTH
	bodyLength := binary.BigEndian.Uint32(buffer[bodyLengthBegin:bodyLengthEnd])
	logger.Debug("body length:", bodyLength)

	bodyBegin := bodyLengthEnd
	bodyEnd := bodyBegin + (int(bodyLength) - CHECKSUM_LENGTH)
	body := buffer[bodyBegin:bodyEnd]
	logger.Debug("body:", hex.EncodeToString(body))

	bodyChecksumBegin := bodyEnd
	bodyChecksumEnd := bodyChecksumBegin + CHECKSUM_LENGTH
	bodyChecksum := buffer[bodyChecksumBegin:bodyChecksumEnd]
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
