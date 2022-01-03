package network

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"

	"github.com/alexcoder04/rfap-go-server/log"
	"github.com/alexcoder04/rfap-go-server/settings"
	"github.com/alexcoder04/rfap-go-server/utils"
	"github.com/jchavannes/go-pgp/pgp"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func sendDataEncrypted(client utils.Client, data []byte) error {
	checksum := sha256.Sum256(data)
	dataEncrypted, err := pgp.Encrypt(client.PubkeyEntity, utils.ConcatBytes(data, checksum[:]))
	if err != nil {
		return err
	}
	length := uint32(len(dataEncrypted))
	lengthBytes := make([]byte, settings.CONT_LEN_INDIC_LENGTH)
	binary.BigEndian.PutUint32(lengthBytes, length)
	allData := utils.ConcatBytes(lengthBytes, dataEncrypted)

	i := 0
	for {
		if (i + settings.MAX_BYTES_SEND_AT_ONCE) > len(allData) {
			_, err := client.Conn.Write(allData[i:])
			if err != nil {
				return err
			}
			break
		}
		_, err := client.Conn.Write(allData[i : i+settings.MAX_BYTES_SEND_AT_ONCE])
		if err != nil {
			return err
		}
		i += settings.MAX_BYTES_SEND_AT_ONCE
	}
	return nil
}

func sendPacket(client utils.Client, command uint32, metadata utils.HeaderMetadata, body []byte) error {
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
		"client": client.Address,
	}).Trace("header length: ", len(metadataBytes))
	if len(metadataBytes) > (1024 * 8) {
		return &utils.ErrInvalidContentLength{}
	}

	// send version
	_, err = client.Conn.Write(version)
	if err != nil {
		return err
	}

	// send header
	err = sendDataEncrypted(client, utils.ConcatBytes(commandBytes, metadataBytes))
	if err != nil {
		return err
	}

	// send body
	err = sendDataEncrypted(client, body)
	if err != nil {
		return err
	}

	// success
	log.Logger.WithFields(logrus.Fields{
		"client": client.Address,
	}).Info("sent packet 0x", hex.EncodeToString(commandBytes))
	return nil
}
