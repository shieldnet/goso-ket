package utils

import (
	"log"
	"net"
)

func GetListenerIpAddress(ip, port string) string {
	return ip + ":" + port
}

func MakeTCPListener(address string) (net.Listener, error) {
	return net.Listen("tcp", address)
}

func MakeUDPListener(address string) (net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
		return net.UDPConn{}, err
	}
	uconn, err := net.ListenUDP("udp", addr)
	return *uconn, err
}
