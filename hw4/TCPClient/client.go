package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Request struct {
	Command string                 `json:"command"`
	Param   map[string]interface{} `json:"param"`
}

type Response struct {
	Status  string                 `json:"status"`
	Message string                 `json:"answer"`
	Items   map[string]interface{} `json:"items"`
}

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Usage:\n")
		fmt.Printf("\tgo run client.go <nickname>\n")
		os.Exit(0)
	}

	fmt.Printf("Your nickname is set to %s\n", args[1])

	ip := "127.0.0.1"
	port := "11227"

	network := "tcp"

	rand.Seed(time.Now().UnixNano())
	n, _ := uuid.NewV4()
	name := n.String()

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGPIPE, syscall.SIGINT, syscall.SIGTERM)

	// Make Connection
	conn, err := net.Dial(network, ip+":"+port)

	// SIGNAL detector (goroutine)
	go func(conn net.Conn) {
		go Listener(conn)
		<-sigs
		println("adios~")
		conn.Close()
		os.Exit(1)
	}(conn)

	if err != nil {
		log.Fatalln(err)
		return
	} else {
		Join(conn, name)
	}

	for {
		in := bufio.NewReader(os.Stdin)
		msg, err := in.ReadString('\n')
		msg = strings.ReplaceAll(msg, "\n", "")

		if err != nil {
			// 에러처리
		}

		Parse(conn, msg)

	}

}

func Parse(conn net.Conn, msg string) {
	if strings.Contains(msg, "\\wh") {
		s := strings.Split(msg, " ")
		if len(s) < 3 {
			return
		}
		originMsg := strings.Join(s[2:], " ")
		Whisper(conn, s[1], originMsg)
	} else if strings.Contains(msg, "\\rename") {
		s := strings.Split(msg, " ")
		if len(s) < 2 {
			return
		}
		nickName := strings.Join(s[1:], " ")
		Rename(conn, nickName)
	} else if strings.Contains(msg, "\\users") {
		GetUserList(conn)
	} else if strings.Contains(msg, "\\version") {
		GetVersion(conn)
	} else if strings.Contains(msg, "\\rtt") {
		GetRtt(conn)
	} else if strings.Contains(msg, "\\exit") {
		println("adios~")
		conn.Close()
		os.Exit(1)
	} else {
		Say(conn, msg)
	}

}

func Join(conn net.Conn, name string) {
	var clientRequest = Request{
		Command: "\\join",
		Param: map[string]interface{}{
			"name": name,
		},
	}
	b, _ := json.Marshal(clientRequest)

	conn.Write(b)
}

func Say(conn net.Conn, say string) {
	var clientRequest = Request{
		Command: "\\say",
		Param: map[string]interface{}{
			"message": say,
		},
	}
	b, _ := json.Marshal(clientRequest)

	conn.Write(b)
}

func Whisper(conn net.Conn, to, msg string) {
	var clientRequest = Request{
		Command: "\\wh",
		Param: map[string]interface{}{
			"message": msg,
			"user":    to,
		},
	}
	b, _ := json.Marshal(clientRequest)

	conn.Write(b)

}

func Rename(conn net.Conn, nickName string) {
	var clientRequest = Request{
		Command: "\\rename",
		Param: map[string]interface{}{
			"name": nickName,
		},
	}
	b, _ := json.Marshal(clientRequest)

	conn.Write(b)
}

func GetUserList(conn net.Conn) {
	var clientRequest = Request{
		Command: "\\users",
		Param:   map[string]interface{}{},
	}
	b, _ := json.Marshal(clientRequest)

	conn.Write(b)
}

func GetVersion(conn net.Conn) {
	var clientRequest = Request{
		Command: "\\version",
		Param:   map[string]interface{}{},
	}
	b, _ := json.Marshal(clientRequest)

	conn.Write(b)
}

func GetRtt(conn net.Conn) {
	var clientRequest = Request{
		Command: "\\rtt",
		Param:   map[string]interface{}{
			"time": time.Now().UnixNano(),
		},
	}
	b, _ := json.Marshal(clientRequest)

	conn.Write(b)
}

func Listener(conn net.Conn) {
	buffer := make([]byte, 1024)
	for {
		count, err := conn.Read(buffer)
		if err != nil {
			if io.EOF == err {
				log.Printf("Connection End: %v", conn.RemoteAddr().String())
				os.Exit(0)
			} else {
				log.Printf("Receive Failed: %v", err)
				os.Exit(0)
			}
		}

		if count > 0 {
			data := buffer[:count]
			req := Request{}

			if err := json.Unmarshal(data, &req); err != nil {
				log.Fatalln(err)
				return
			}

			//fmt.Println("got packet: " + string(data))

			if req.Command == "\\say" {
				fmt.Printf("%s> %s\n", req.Param["from"], req.Param["message"].(string))
			} else if req.Command == "\\rtt" {
				now := time.Now().UnixNano()
				sub := float64(now - int64(req.Param["time"].(float64))) / 1000000.0
				fmt.Printf("Your rtt is %s ms.\n", fmt.Sprintf("%f",sub))
			} else if req.Command == "\\notice" {
				fmt.Printf("%s\n", req.Param["message"])
			} else if req.Command == "\\users" {
				users := req.Param["users"].([]interface{})
				for _, v := range users {
					a := v.(map[string]interface{})
					fmt.Printf("\t%s\t\t%s\t%s\n", a["name"].(string), a["ip"].(string), a["port"].(string))
				}
			} else if req.Command == "\\version" {
				fmt.Printf("[ version is %s ]\n",req.Param["version"].(string))
			}
		}
	}

}
