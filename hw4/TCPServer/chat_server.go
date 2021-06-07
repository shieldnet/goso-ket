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
	Commands           map[string]func(Request, *Client) Response
	HiddenCommands     map[string]func(Request, *Client) Response
	TotalClientsNumber int32
	NthClient          int32
}

type Handler interface {
	GetUserList(req Request, client *Client) Response
	Whisper(req Request, client *Client) Response
	Exit(req Request, client *Client) Response
	GetVersion(req Request, client *Client) Response
	ChangeName(req Request, client *Client) Response
	GetRtt(req Request, client *Client) Response
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
	cs.Commands = map[string]func(request Request, client *Client) Response{
		"\\version": cs.GetVersion,
		"\\users":   cs.GetUserList,
		"\\wh":      cs.Whisper,
		"\\rename":  cs.ChangeName,
		"\\rtt":     cs.GetRtt,
	}

	// Hidden Command
	cs.HiddenCommands = map[string]func(request Request, client *Client) Response{
		"\\join": cs.Join,
		"\\say":  cs.Say,
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
			req := Request{}
			resp := Response{}

			if err := json.Unmarshal(data, &req); err != nil {
				log.Fatalln(err)
				return
			}

			// Handle Commands
			for k := range c.Commands {
				if req.Command == k {
					resp = c.Commands[k](req, client)
				}
			}

			// Handle Hidden Commands
			for k := range c.HiddenCommands {
				if req.Command == k {
					resp = c.HiddenCommands[k](req, client)
				}
			}

			if resp.Status == "" {
				resp.Status = "404"
				resp.Message = fmt.Sprintf(`The Command "%s" is not exist.`, req.Command)
			}

			b, _ := json.Marshal(resp)
			conn.Write(b)
		}

	}
}

func (c *ChatServer) Join(req Request, client *Client) Response {
	res := Response{}

	if req.Command == "\\join" {
		c.Clients[req.Param["name"]] = client
	} else {
		res.Status = "400"
		res.Message = ""
	}
	return res
}

func (c *ChatServer) Say(req Request, client *Client) Response {
	res := Response{}

	clientRequest := Request{
		Command: "\\say",
		Param: map[string]string{
			"from":    client.Name,
			"message": req.Param["message"],
		},
	}

	b, _ := json.Marshal(clientRequest)

	for _, cl := range c.Clients {
		cl.Connection.Write(b)
	}

	res.Message = "Completed"
	res.Status = "200"

	return res
}

func (c *ChatServer) GetVersion(req Request, client *Client) Response {
	return Response{}
}

func (c *ChatServer) GetUserList(req Request, client *Client) Response {
	return Response{}
}

func (c *ChatServer) Whisper(req Request, client *Client) Response {
	return Response{}
}

func (c *ChatServer) ChangeName(req Request, client *Client) Response {
	return Response{}
}

func (c *ChatServer) GetRtt(req Request, client *Client) Response {
	return Response{}
}

func (c *ChatServer) Exit(req Request, client *Client) Response {

	return Response{}
}
