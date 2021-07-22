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
	Version            string
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

	cs.Version = "1.0.0"

	cs.Clients = map[string]*Client{}

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
		println(conn.RemoteAddr().String())
		remoteAddr := strings.Split(conn.RemoteAddr().String(), ":")
		fmt.Printf("Connection requested from ('%s', '%s')\n", remoteAddr[0], remoteAddr[1])

		go func() {
			remoteAddr := strings.Split(conn.RemoteAddr().String(), ":")
			cl := &Client{
				Ip:              remoteAddr[0],
				Port:            remoteAddr[1],
				NthClientNumber: c.NthClient,
				Connection:      conn,
			}

			c.NthClient++

			c.TotalClientsNumber += 1
			c.Handle(conn, cl)
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

			// Check strings contains "i hate professor"
			if req.Command == "\\say" {
				if strings.Contains(strings.ToLower(req.Param["message"].(string)), "i hate professor") {
					c.Exit(req, client)
					msg := fmt.Sprintf(
						"[%s is disconnected. There are %d users in the chat room.]",
						client.Name, GetChatServer().TotalClientsNumber,
					)
					c.Broadcast(msg)
					return
				}
			}

			// Handle Commands
			for k := range c.Commands {
				if req.Command == k {
					log.Println("Handle Command:: " + req.Command)
					resp = c.Commands[k](req, client)
				}
			}

			// Handle Hidden Commands
			for k := range c.HiddenCommands {
				if req.Command == k {
					log.Println("Handle Command:: " + req.Command)
					resp = c.HiddenCommands[k](req, client)
				}
			}

			if resp.Status == "" {
				resp.Status = "404"
				resp.Message = fmt.Sprintf(`The Command "%s" is not exist.`, req.Command)
			}

			//b, _ := json.Marshal(resp)
			//conn.Write(b)
		}

	}
	c.Exit(Request{}, client)
}

func (c *ChatServer) Broadcast(message string) {
	clientRequest := Request{
		Command: "\\notice",
		Param: map[string]interface{}{
			"message": message,
		},
	}

	b, _ := json.Marshal(clientRequest)

	for _, cl := range c.Clients {
		cl.Connection.Write(b)
	}
}

func (c *ChatServer) Join(req Request, client *Client) Response {
	res := Response{}

	if req.Command == "\\join" {
		log.Println(req.Param["name"].(string) + " is joined.")
		c.Clients[req.Param["name"].(string)] = client
		client.Name = req.Param["name"].(string)

		res.Status = "200"
		res.Message = "Ok, welcome " + req.Param["name"].(string)
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
		Param: map[string]interface{}{
			"from":    client.Name,
			"message": req.Param["message"],
		},
	}

	b, _ := json.Marshal(clientRequest)

	fmt.Println(len(c.Clients))

	for _, cl := range c.Clients {
		if cl.Name != client.Name {
			log.Println("[Say] " + cl.Name + " to send ::" + req.Param["message"].(string))
			cl.Connection.Write(b)
		}
	}

	res.Message = "Say Completed"
	res.Status = "200"

	return res
}

func (c *ChatServer) GetVersion(req Request, client *Client) Response {
	clientRequest := Request{
		Command: "\\version",
		Param: map[string]interface{}{
			"version": c.Version,
		},
	}

	b, _ := json.Marshal(clientRequest)
	client.Connection.Write(b)

	return Response{
		Status:  "200",
		Message: "OK",
	}
}

func (c *ChatServer) GetUserList(req Request, client *Client) Response {
	var (
		userList []map[string]string
	)

	for _, cl := range c.Clients {
		userList = append(userList, map[string]string{
			"name": cl.Name, "ip": cl.Ip, "port": cl.Port,
		})
	}

	clientRequest := Request{
		Command: "\\users",
		Param: map[string]interface{}{
			"users": userList,
		},
	}

	b, _ := json.Marshal(clientRequest)
	client.Connection.Write(b)

	return Response{
		Status:  "200",
		Message: "OK",
	}
}

func (c *ChatServer) Whisper(req Request, client *Client) Response {
	userName := req.Param["user"]

	clientRequest := Request{
		Command: "\\say",
		Param: map[string]interface{}{
			"from":    client.Name,
			"message": req.Param["message"],
		},
	}

	b, _ := json.Marshal(clientRequest)
	c.Clients[userName.(string)].Connection.Write(b)

	return Response{
		Status:  "200",
		Message: "OK",
	}
}

func (c *ChatServer) ChangeName(req Request, client *Client) Response {
	newName := req.Param["name"]

	delete(c.Clients, client.Name)
	c.Clients[newName.(string)] = client
	client.Name = newName.(string)

	return Response{
		Status:  "200",
		Message: "OK",
	}
}

func (c *ChatServer) GetRtt(req Request, client *Client) Response {

	clientRequest := Request{
		Command: "\\rtt",
		Param: map[string]interface{}{
			"time": int64(req.Param["time"].(float64)),
		},
	}

	b, _ := json.Marshal(clientRequest)
	client.Connection.Write(b)

	return Response{
		Status:  "200",
		Message: "OK",
	}
}

func (c *ChatServer) Exit(req Request, client *Client) Response {
	defer client.Connection.Close()
	c.TotalClientsNumber -= 1
	delete(c.Clients, client.Name)

	return Response{}
}
