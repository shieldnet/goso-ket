package main

import (
	"io"
	"log"
	"net"
)

func main() {
	network := "tcp"
	ip := "0.0.0.0"
	port := "8001"
	address := ip + ":" + port

	listener, err := net.Listen(network, address)

	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		log.Fatalln(err)
	}

	buffer := make([]byte, 1024)
	for {
		count, err := conn.Read(buffer)
		if err != nil {
			if io.EOF == err {
				log.Printf("Connection End: %v", conn.RemoteAddr().String())
			} else {
				log.Printf("Receive Failed: %v", err)
			}
			break

		}
		if count > 0 {
			data := buffer[:count]
			log.Println(string(data))
		}
	}

}
