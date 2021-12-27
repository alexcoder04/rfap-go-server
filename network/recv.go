package network

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"net"
	"time"

	"github.com/alexcoder04/rfap-go-server/log"
	"github.com/alexcoder04/rfap-go-server/settings"
	"github.com/alexcoder04/rfap-go-server/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func recvPacket(conn net.Conn) (uint32, uint32, utils.HeaderMetadata, []byte, error) {
	// version
	versionBytes := make([]byte, settings.VERSION_LENGTH)
	err := conn.SetReadDeadline(time.Now().Add(settings.CONN_RECV_TIMEOUT_SECS * time.Second))
	if err != nil {
		return 0, 0, utils.HeaderMetadata{}, make([]byte, 0), &utils.ErrSetReadTimeoutFailed{}
	}
	_, err = conn.Read(versionBytes[:])
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return 0, 0, utils.HeaderMetadata{}, make([]byte, 0), &utils.ErrClientCrashed{}
		}
		return 0, 0, utils.HeaderMetadata{}, make([]byte, 0), err
	}
	version := uint32(binary.BigEndian.Uint16(versionBytes))
	log.Logger.Debug("version:", version)
	if !utils.Uint32ArrayContains(settings.SUPPORTED_RFAP_VERSIONS, version) {
		return version, 0, utils.HeaderMetadata{}, make([]byte, 0), &utils.ErrUnsupportedRfapVersion{}
	}

	// header length
	headerLengthBytes := make([]byte, settings.CONT_LEN_INDIC_LENGTH)
	_, err = conn.Read(headerLengthBytes[:])
	if err != nil {
		return version, 0, utils.HeaderMetadata{}, make([]byte, 0), err
	}
	headerLength := binary.BigEndian.Uint32(headerLengthBytes)
	log.Logger.Debug("header length: ", headerLength, ":0x", hex.EncodeToString(headerLengthBytes))
	if headerLength > (1024 * 8) {
		return version, 0, utils.HeaderMetadata{}, make([]byte, 0), &utils.ErrInvalidContentLength{}
	}

	// raw header
	headerRaw := make([]byte, headerLength)
	_, err = conn.Read(headerRaw[:])
	if err != nil {
		return version, 0, utils.HeaderMetadata{}, make([]byte, 0), err
	}

	// command
	command := binary.BigEndian.Uint32(headerRaw[:4])
	log.Logger.Debug("command: 0x" + hex.EncodeToString(headerRaw[:4]))

	// metadata
	headerBytes := headerRaw[4 : len(headerRaw)-32]
	log.Logger.Debug("header:", hex.EncodeToString(headerBytes))

	// header checksum
	headerChecksum := headerRaw[len(headerRaw)-32:]
	headerChecksumExpected := sha256.Sum256(headerRaw[:len(headerRaw)-32])
	if !bytes.Equal(headerChecksum, headerChecksumExpected[:]) {
		log.Logger.Debug("header checksum:", hex.EncodeToString(headerChecksum))
		log.Logger.Debug("header checksum expected:", hex.EncodeToString(headerChecksumExpected[:]))
		return version, 0, utils.HeaderMetadata{}, make([]byte, 0), &utils.ErrChecksumsNotMatching{}
	}

	// body length
	bodyLengthBytes := make([]byte, settings.CONT_LEN_INDIC_LENGTH)
	_, err = conn.Read(bodyLengthBytes[:])
	if err != nil {
		return version, 0, utils.HeaderMetadata{}, make([]byte, 0), err
	}
	bodyLength := binary.BigEndian.Uint32(bodyLengthBytes)
	log.Logger.Debug("body length:", bodyLength)

	bodyRaw := make([]byte, bodyLength)
	// TODO recv in loop, body_length may be very big
	_, err = conn.Read(bodyRaw[:])
	if err != nil {
		return version, 0, utils.HeaderMetadata{}, make([]byte, 0), err
	}
	body := bodyRaw[:len(bodyRaw)-32]
	bodyChecksum := bodyRaw[len(bodyRaw)-32:]
	bodyChecksumExpected := sha256.Sum256(body)
	if !bytes.Equal(bodyChecksum, bodyChecksumExpected[:]) {
		log.Logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Debug("body checksum: ", hex.EncodeToString(bodyChecksum))
		log.Logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Debug("body checksum expected: ", hex.EncodeToString(bodyChecksumExpected[:]))
		return version, 0, utils.HeaderMetadata{}, make([]byte, 0), &utils.ErrChecksumsNotMatching{}
	}

	// parse
	header := utils.HeaderMetadata{}
	err = yaml.Unmarshal(headerRaw[4:len(headerRaw)-32], &header)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"client": conn.RemoteAddr().String(),
		}).Error("error decoding metadata")
		return version, command, utils.HeaderMetadata{}, body, err
	}

	// return
	return version, command, header, body, nil
}
