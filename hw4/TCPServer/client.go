package TCPServer

import "net"

type Client struct {
	ClientInfoCommand
	ClientNumber int32
	NthClientNumber int32
	Ip string
	Port string
	Name string
	Connection net.Conn
}

type ClientInfoCommand interface {
	SetName(string)
}

func (c *Client) SetName(nickName string) {
	c.Name = nickName
}