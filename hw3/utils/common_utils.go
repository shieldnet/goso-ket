package utils

import (
	"net"
)

func GetListenerIpAddress(ip, port string) string {
	return ip + ":" + port
}

func MakeTCPListener(address string) (net.Listener, error) {
	return net.Listen("tcp", address)
}