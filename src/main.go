package main

import (
	"fmt"
	"net"
	"os"
)

const (
	connHost        = "localhost"
	connPort        = "6700"
	connType        = "tcp"
	ProtocolVersion = 1
)

func main() {
	Init()
	fmt.Println("Starting " + connType + " server on " + connHost + ":" + connPort)

	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}
		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

		go HanleConnection(c)
	}
}
