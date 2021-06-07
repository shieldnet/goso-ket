package main

import (
	"github.com/shieldnet/goso-ket/hw4/TCPServer"
)

func main(){
	cs := TCPServer.InitChatServer("0.0.0.0", "11227")
	cs.Run()
}