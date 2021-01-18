package main

import (
	"log"
	"net"
	"utils"
)

func main(){
	network := "tcp"
	address := utils.GetListenerIpAddress("0.0.0.0", "8001")

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