package main

import (
	"net"
	"log"
)

func main(){
	network := "tcp"
	ip := "0.0.0.0"
	port := "8001"
	address := ip + ":" + port

	conn,err := net.Dial(network, address)
	if err!=nil {
		log.Fatalln(err)
		return
	} else {
		_, err := conn.Write([]byte("test connection!!!"))
		if err!=nil{
			log.Fatalln(err)
			return
		}
		conn.Close()
	}


}