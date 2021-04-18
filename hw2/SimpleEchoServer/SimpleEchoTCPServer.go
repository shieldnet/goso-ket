package main

import (
	"bytes"
	"fmt"
	"github.com/shieldnet/goso-ket/hw2/utils"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	ip := "0.0.0.0"
	port := "8001"

	address := utils.GetListenerIpAddress(ip, port)

	// Make TCP Listener
	listener, err := utils.MakeTCPListener(address)
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	for {
		// TCP has to make connection
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		remoteAddr := strings.Split(conn.RemoteAddr().String(), ":")

		fmt.Printf("Connection requested from ('%s', '%s')\n", remoteAddr[0], remoteAddr[1] )
		handleTCP(conn)
		conn.Close()
	}
}

func handleTCP(connection net.Conn) {
	buffer := make([]byte, 1024)
	for {
		count, err := connection.Read(buffer)
		if err != nil {
			if io.EOF == err {
				log.Printf("Connection End: %v", connection.RemoteAddr().String())
			} else {
				log.Printf("Receive Failed: %v", err)
			}
			break
		}
		if count > 0 {
			data := buffer[:count]
			connection.Write(bytes.ToUpper(data))
		}
	}
}
