package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ip := "0.0.0.0"
	port := "11227"

	network := "tcp"

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Make Connection
	conn, err := net.Dial(network, ip+":"+port)

	// SIGNAL detector (goroutine)
	go func(conn net.Conn) {
		<-sigs
		println("Signal detected, ", sigs)
		conn.Close()
		os.Exit(1)
	}(conn)

	if err != nil {
		log.Fatalln(err)
		return
	} else {
		conn.Close()
	}
}

