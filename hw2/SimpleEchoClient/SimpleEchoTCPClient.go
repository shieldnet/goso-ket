package main

import (
	"bufio"
	"fmt"
	"github.com/shieldnet/goso-ket/hw2/utils"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	ip := "0.0.0.0"
	port := "8001"

	network := "tcp"
	address := utils.GetListenerIpAddress(ip, port)

	conn, err := net.Dial(network, address)
	if err != nil {
		log.Fatalln(err)
		return
	} else {
		// Print Greeting message
		fmt.Printf("The Client is running on port %s.\n", port)
		fmt.Printf("Input lowercase sentence: ")

		// Get input from user
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')

		// Convert input string to bytes, and send it to server
		_, err = conn.Write([]byte(input))
		if err != nil {
			log.Fatalln(err)
			return
		}

		// Try to receive answer from server
		buffer := make([]byte, 1024)
		count, err := conn.Read(buffer)
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
			fmt.Printf("Reply from server: %s\n", string(data))
		}
		conn.Close()
	}
}
