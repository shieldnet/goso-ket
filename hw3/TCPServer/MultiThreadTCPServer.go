package TCPServer

import (
	"fmt"
	"github.com/shieldnet/goso-ket/hw3/utils"
	"log"
	"strings"
	"time"
)

var tcpServerInfo = TCPServerInfo{}

func MultiThreadTCPServer() {
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

	totalClientsNumber := 0
	clientNumber := 1

	// Print Number of clients every minute
	go func(totalClientsNumber *int) {
		for true {
			fmt.Printf("Number of connected clients = %d\n", *totalClientsNumber)
			time.Sleep(time.Second * 60)
		}
	}(&totalClientsNumber)

	for {
		// TCP has to make connection
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		remoteAddr := strings.Split(conn.RemoteAddr().String(), ":")

		fmt.Printf("Connection requested from ('%s', '%s')\n", remoteAddr[0], remoteAddr[1])

		// Use `goroutine` instead of thread of OS
		go func(clientNumber, totalClientsNumber *int) {
			thisClientNumber := *clientNumber
			*totalClientsNumber++
			*clientNumber++
			fmt.Printf("Client %d connected. Number of connected clients = %d\n", thisClientNumber, *totalClientsNumber)

			handle(conn)
			conn.Close()

			*totalClientsNumber--
			fmt.Printf("Client %d disconnected. Number of connected clients = %d\n", thisClientNumber, *totalClientsNumber)
		}(&clientNumber, &totalClientsNumber)
	}
}