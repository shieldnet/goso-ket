package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Request struct {
	Command string            `json:"command"`
	Param   map[string]string `json:"param"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"answer"`
	Items   map[string]interface{} `json:"items"`
}

func main() {
	ip := "127.0.0.1"
	port := "11227"

	network := "tcp"

	name := "kim-client"

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
		Join(conn, name)
		time.Sleep(10*time.Second)
	}
}

func Join(conn net.Conn, name string) {
	var clientRequest = Request{
		Command: "\\join",
		Param: map[string]string{
			"name": name,
		},
	}
	b, _ := json.Marshal(clientRequest)

	conn.Write(b)
}

