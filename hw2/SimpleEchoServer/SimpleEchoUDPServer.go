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

	// Make UDP Listener
	conn, err := utils.MakeUDPListener(address)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	for {
		handleUDP(&conn)
		conn.Close()
	}
}

func handleUDP(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	for {
		count, addr, err := conn.ReadFromUDP(buffer)
		if addr == nil {
			log.Fatalln("addr is nil.")
			return
		}
		if err != nil {
			log.Fatalln(err)
			return
		}

		remoteAddr := strings.Split(addr.String(), ":")
		fmt.Printf("Connection requested from ('%s', '%s')\n", remoteAddr[0], remoteAddr[1] )
		if err != nil {
			if io.EOF == err {
				log.Printf("Connection End: %v", addr.String())
			} else {
				log.Printf("Receive Failed: %v", err)
			}
			break
		}
		conn.WriteTo(bytes.ToUpper(buffer[:count]), addr)
	}
}
