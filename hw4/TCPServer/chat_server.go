package TCPServer

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type ChatServer struct {
	Handler
	Ip                 string
	Port               string
	Listener           net.Listener
	Clients            map[string]*Client
	Commands           map[string]func(Request) Response
	HiddenCommands 	   map[string]func(Request) Response
	TotalClientsNumber int32
	NthClient          int32
}

type Handler interface {
	GetUserList(req Request) Response
	Whisper(req Request) Response
	Exit(req Request) Response
	GetVersion(req Request) Response
	ChangeName(req Request) Response
	GetRtt(req Request) Response
}

var basicChatServer *ChatServer = nil

func GetChatServer() *ChatServer {
	if basicChatServer == nil {
		InitChatServer("0.0.0.0", "11227")
	}
	return basicChatServer
}

func InitChatServer(ip, port string) *ChatServer {
	cs := &ChatServer{}

	cs.Ip = ip
	cs.Port = port

	listener, err := net.Listen("tcp", strings.Join([]string{cs.Ip, cs.Port}, ":"))
	if err != nil {
		log.Fatal(err)
	}
	cs.Listener = listener
	cs.TotalClientsNumber = 0
	cs.NthClient = 1
	basicChatServer = cs

	// Commmand
	cs.Commands = map[string]func(request Request) Response{
		"\\version": cs.GetVersion,
		"\\users":   cs.GetUserList,
		"\\wh":      cs.Whisper,
		"\\rename":  cs.ChangeName,
		"\\rtt":     cs.GetRtt,
	}
	return cs
}

func (c *ChatServer) Run() {
	for {
		if c.Listener == nil {
			listener, err := net.Listen("tcp", strings.Join([]string{c.Ip, c.Port}, ":"))
			if err != nil {
				log.Fatal(err)
			}
			c.Listener = listener
		}

		conn, err := c.Listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		remoteAddr := strings.Split(conn.RemoteAddr().String(), ":")
		fmt.Printf("Connection requested from ('%s', '%s')\n", remoteAddr[0], remoteAddr[1])

		go func() {
			cl := &Client{}
			c.Handle(conn, cl)
			defer conn.Close()
		}()

	}
}

func (c *ChatServer) Handle(conn net.Conn, client *Client) {
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
			msg := Request{}
			resp := Response{}

			if err := json.Unmarshal(data, &msg); err != nil {
				log.Fatalln(err)
				return
			}

			// Handle Commands


			b, _ := json.Marshal(resp)
			conn.Write(b)
		}

	}
}
