package BasicTCPServer

import (
	"fmt"
	"github.com/shieldnet/goso-ket/hw2/utils"
	"log"
	"strings"
	"time"
)

var tcpServerInfo = TCPServerInfo{}

func TCPServer() {
	ip := "0.0.0.0"
	port := "8001"

	// Set Basic Server Information
	tcpServerInfo.StartTime = time.Now()
	externalIP, _ := getMyIPAddrAndPort()
	tcpServerInfo.IP = externalIP
	tcpServerInfo.Port = port

	// ip:port string
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

		fmt.Printf("Connection requested from ('%s', '%s')\n", remoteAddr[0], remoteAddr[1])
		handle(conn)
		conn.Close()
	}
}