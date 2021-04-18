package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/shieldnet/goso-ket/hw2/utils"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	type (
		MenuPayload struct {
			Menu    string `json:"menu"`
			Message string `json:"message"`
		}

		AnswerPayload struct {
			Answer string `json:"answer"`
			Error  string `json:"error"`
		}
	)

	sigs := make(chan os.Signal)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ip := "0.0.0.0"
	port := "8001"

	network := "tcp"
	address := utils.GetListenerIpAddress(ip, port)

	// Make Connection
	conn, err := net.Dial(network, address)

	// SIGNAL detector goroutine
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
		// Print Greeting message
		fmt.Printf("The Client is running on port %s.\n", port)
		menu := "99"

		// Get Menu code from user
		for menu != "5" {
			// Print Menu
			PrintMenu()
			fmt.Printf("Input Option: ")

			// Get input from user
			reader := bufio.NewReader(os.Stdin)
			menu, _ = reader.ReadString('\n')
			menu = strings.Replace(menu, "\n", "", -1)

			// Convert input string to bytes, and send it to server
			m := MenuPayload{}
			m.Menu = menu

			// Process Menu
			switch menu {
			case "1":
				fmt.Printf("Input sentence: ")
				is, _ := reader.ReadString('\n')
				is = strings.Replace(is, "\n", "", -1)
				m.Message = is
				break
			case "2":
				break
			case "3":
				break
			case "4":
				break
			case "5":
				return
			default:
				break
			}

			// Make payload with json
			mBytes, _ := json.Marshal(m)

			// Send the payload to server
			st := time.Now()
			_, err = conn.Write(mBytes)
			if err != nil {
				log.Fatalln(err)
				return
			}

			// Try to receive answer from server
			buffer := make([]byte, 1024)
			count, err := conn.Read(buffer)
			fin := time.Now()
			if err != nil {
				if io.EOF == err {
					log.Printf("Connection End: %v", conn.RemoteAddr().String())
				} else {
					log.Printf("Receive Failed: %v", err)
				}
			}

			// Received bytes > 0
			if count > 0 {
				data := buffer[:count]
				res := AnswerPayload{}
				json.Unmarshal(data, &res)
				fmt.Printf("Reply from server: %s\n", res.Answer)

				diff := fin.Sub(st).Microseconds()

				fmt.Printf("Respionse Time: %.3f ms\n", float64(diff)/1000.0)
			}

		}
		conn.Close()
	}
}

func PrintMenu() {
	println("<Menu>")
	println("1) convert text to UPPER-case")
	println("2) get my IP address and port number")
	println("3) get server time")
	println("4) get server running time")
	println("5) exit")
}
