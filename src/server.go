package main

import (
	"io"
	"log"
	"net"
	"utils"
)

func main() {
	address := utils.GetListenerIpAddress("0.0.0.0","8001")
	listener, err := utils.MakeTcpListener(address)
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		handle(conn)
		conn.Close()
	}
	return
}

func handle(connection net.Conn){
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
			log.Println(string(data))
		}
	}
	return
}
