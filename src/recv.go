package main

import (
	"encoding/binary"
	//"encoding/hex"
	"log"
	"net"

	"gopkg.in/yaml.v3"
)

func RecvPacket(conn net.Conn) (uint32, uint32, HeaderValues, []byte, error) {
	// receive
	buffer := make([]byte, 2+4+(16*1024*1024)+4+(16*1024*1024))
	_, err := conn.Read(buffer[:])
	if err != nil {
		log.Println(err.Error())
		conn.Close()
		return 0, 0, HeaderValues{}, make([]byte, 0), err
	}

	// sort and split
	version := uint32(binary.BigEndian.Uint16(buffer[:2]))
	//log.Println("version:", version)
	if !Uint32ArrayContains(SUPPORTED_RFAP_VERSIONS, version) {
		return 0, 0, HeaderValues{}, make([]byte, 0), err
	}

	headerLength := binary.BigEndian.Uint32(buffer[2 : 2+4])
	//log.Println("header length:", headerLength)

	command := binary.BigEndian.Uint32(buffer[2+4 : 2+4+4])
	//log.Println("command:", command)

	headerRaw := buffer[2+4+4 : 2+4+(headerLength-32)]
	//log.Println("header:", hex.EncodeToString(headerRaw))

	headerChecksum := buffer[2+4+(headerLength-32) : 2+4+headerLength]
	_ = headerChecksum
	//log.Println("header checksum:", hex.EncodeToString(headerChecksum))

	bodyLength := binary.BigEndian.Uint32(buffer[2+4+headerLength : 2+4+headerLength+4])
	//log.Println("body length:", bodyLength)

	body := buffer[2+4+headerLength+4 : 2+4+headerLength+4+(bodyLength-32)]
	//log.Println("body:", hex.EncodeToString(body))

	bodyChecksum := buffer[2+4+headerLength+4+(bodyLength-32) : 2+4+headerLength+4+bodyLength]
	_ = bodyChecksum
	//log.Println("body checksum:", hex.EncodeToString(bodyChecksum))

	// parse
	header := HeaderValues{}
	yamlErr := yaml.Unmarshal([]byte(headerRaw), &header)
	if yamlErr != nil {
		log.Println("error encoding metadata")
		log.Println(err.Error())
		conn.Close()
		return 0, 0, HeaderValues{}, make([]byte, 0), err
	}

	// return
	return version, command, header, body, nil
}
