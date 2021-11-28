package main

import (
	"encoding/binary"
	"encoding/hex"
	"log"
	"net"

	"gopkg.in/yaml.v3"
)

func RecvPacket(conn net.Conn) (int, int, HeaderValues, []byte, error) {
	// receive
	buffer := make([]byte, 2+3+(16*1024*1024)+3+(16*1024*1024))
	_, err := conn.Read(buffer[:])
	if err != nil {
		log.Println(err.Error())
		conn.Close()
		return 0, 0, HeaderValues{}, make([]byte, 0), err
	}

	// sort and split
	version := binary.BigEndian.Uint16(buffer[:2])
	log.Println("version:", version)

	headerLength := binary.BigEndian.Uint32(buffer[2 : 2+3])
	log.Println("header length:", headerLength)

	command := binary.BigEndian.Uint32(buffer[2+3 : 2+3+4])
	log.Println("command:", command)

	headerRaw := buffer[2+3+4 : 2+3+4+(headerLength-32)]
	log.Println("header:", hex.EncodeToString(headerRaw))

	headerChecksum := buffer[2+3+(headerLength-32) : 2+3+headerLength]
	log.Println("header checksum:", hex.EncodeToString(headerChecksum))

	bodyLength := binary.BigEndian.Uint32(buffer[2+3+headerLength : 2+3+headerLength+3])
	log.Println("body length:", bodyLength)

	body := buffer[2+3+headerLength+3 : 2+3+headerLength+3+(bodyLength-32)]
	log.Println("body:", hex.EncodeToString(body))

	bodyChecksum := buffer[2+3+headerLength+3+(bodyLength-32) : 2+3+headerLength+3+bodyLength]
	log.Println("body checksum:", bodyChecksum)

	// parse
	header := HeaderValues{}
	yamlErr := yaml.Unmarshal([]byte(headerRaw), &header)
	if yamlErr != nil {
		log.Println(err.Error())
		conn.Close()
		return 0, 0, HeaderValues{}, make([]byte, 0), err
	}

	// return
	return int(version), int(command), header, body, nil
}
