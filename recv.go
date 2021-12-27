package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func RecvPacket(conn net.Conn) (uint32, uint32, HeaderMetadata, []byte, error) {
	// version
	versionBytes := make([]byte, VERSION_LENGTH)
	err := conn.SetReadDeadline(time.Now().Add(CONN_RECV_TIMEOUT_SECS * time.Second))
	if err != nil {
		return 0, 0, HeaderMetadata{}, make([]byte, 0), &ErrSetReadTimeoutFailed{}
	}
	_, err = conn.Read(versionBytes[:])
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return 0, 0, HeaderMetadata{}, make([]byte, 0), &ErrClientCrashed{}
		}
		return 0, 0, HeaderMetadata{}, make([]byte, 0), err
	}
	version := uint32(binary.BigEndian.Uint16(versionBytes))
	logger.Debug("version:", version)
	if !Uint32ArrayContains(SUPPORTED_RFAP_VERSIONS, version) {
		return version, 0, HeaderMetadata{}, make([]byte, 0), &ErrUnsupportedRfapVersion{}
	}

	// header length
	headerLengthBytes := make([]byte, CONT_LEN_INDIC_LENGTH)
	_, err = conn.Read(headerLengthBytes[:])
	if err != nil {
		return version, 0, HeaderMetadata{}, make([]byte, 0), err
	}
	headerLength := binary.BigEndian.Uint32(headerLengthBytes)
	logger.Debug("header length: ", headerLength, ":0x", hex.EncodeToString(headerLengthBytes))
	if headerLength > (1024 * 8) {
		return version, 0, HeaderMetadata{}, make([]byte, 0), &ErrInvalidContentLength{}
	}

	// raw header
	headerRaw := make([]byte, headerLength)
	_, err = conn.Read(headerRaw[:])
	if err != nil {
		return version, 0, HeaderMetadata{}, make([]byte, 0), err
	}

	// command
	command := binary.BigEndian.Uint32(headerRaw[:4])
	logger.Debug("command: 0x" + hex.EncodeToString(headerRaw[:4]))

	// metadata
	headerBytes := headerRaw[4 : len(headerRaw)-32]
	logger.Debug("header:", hex.EncodeToString(headerBytes))

	// header checksum
	headerChecksum := headerRaw[len(headerRaw)-32:]
	headerChecksumExpected := sha256.Sum256(headerRaw[:len(headerRaw)-32])
	if !bytes.Equal(headerChecksum, headerChecksumExpected[:]) {
		logger.Debug("header checksum:", hex.EncodeToString(headerChecksum))
		logger.Debug("header checksum expected:", hex.EncodeToString(headerChecksumExpected[:]))
		return version, 0, HeaderMetadata{}, make([]byte, 0), &ErrChecksumsNotMatching{}
	}

	// body length
	bodyLengthBytes := make([]byte, CONT_LEN_INDIC_LENGTH)
	_, err = conn.Read(bodyLengthBytes[:])
	if err != nil {
		return version, 0, HeaderMetadata{}, make([]byte, 0), err
	}
	bodyLength := binary.BigEndian.Uint32(bodyLengthBytes)
	logger.Debug("body length:", bodyLength)

	bodyRaw := make([]byte, bodyLength)
	// TODO recv in loop, body_length may be very big
	_, err = conn.Read(bodyRaw[:])
	if err != nil {
		return version, 0, HeaderMetadata{}, make([]byte, 0), err
	}
	body := bodyRaw[:len(bodyRaw)-32]
	bodyChecksum := bodyRaw[len(bodyRaw)-32:]
	bodyChecksumExpected := sha256.Sum256(body)
	if !bytes.Equal(bodyChecksum, bodyChecksumExpected[:]) {
		logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Debug("body checksum: ", hex.EncodeToString(bodyChecksum))
		logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Debug("body checksum expected: ", hex.EncodeToString(bodyChecksumExpected[:]))
		return version, 0, HeaderMetadata{}, make([]byte, 0), &ErrChecksumsNotMatching{}
	}

	// parse
	header := HeaderMetadata{}
	err = yaml.Unmarshal(headerRaw[4:len(headerRaw)-32], &header)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Error("error decoding metadata")
		return version, command, HeaderMetadata{}, body, err
	}

	// return
	return version, command, header, body, nil
}
