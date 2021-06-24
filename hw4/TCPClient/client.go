package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Request struct {
	Command string            `json:"command"`
	Param   map[string]string `json:"param"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"answer"`
	Items   map[string]interface{} `json:"items"`
}

func main() {
	ip := "127.0.0.1"
	port := "11227"

	network := "tcp"

	rand.Seed(time.Now().UnixNano())
	n, _ := uuid.NewV4()
	name := n.String()

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Make Connection
	conn, err := net.Dial(network, ip+":"+port)

	// SIGNAL detector (goroutine)
	go func(conn net.Conn) {
		<-sigs
		println("Signal detected, ", sigs)
		conn.Close()
		os.Exit(1)
	}(conn)

	if err != nil {
		log.Fatalln(err)
		return
	} else {
		Join(conn, name)
	}

	go Listener(conn)

	for {
		in := bufio.NewReader(os.Stdin)
		msg, err := in.ReadString('\n')
		msg = strings.ReplaceAll(msg, "\n","")

		if err != nil {
			// 에러처리
		}

		Say(conn, msg)
	}

}

func Join(conn net.Conn, name string) {
	var clientRequest = Request{
		Command: "\\join",
		Param: map[string]string{
			"name": name,
		},
	}
	b, _ := json.Marshal(clientRequest)

	conn.Write(b)
}

func Say(conn net.Conn, say string) {
	var clientRequest = Request{
		Command: "\\say",
		Param: map[string]string{
			"message": say,
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
			} else {
				log.Printf("Receive Failed: %v", err)
			}
			break
		}

		if count > 0 {
			data := buffer[:count]
			req := Request{}

			if err := json.Unmarshal(data, &req); err != nil {
				log.Fatalln(err)
				return
			}


			//fmt.Println("got packet: "+string(data))

			if req.Command == "\\say" {
				fmt.Printf("%s> %s\n", req.Param["from"], req.Param["message"])
			}
		}
	}

}