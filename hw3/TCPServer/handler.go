package TCPServer

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func handle(connection net.Conn) {
	buffer := make([]byte, 1024)
	for {
		count, err := connection.Read(buffer)
		if err != nil {
			if io.EOF == err {
				log.Printf("Connection End: %v", connection.RemoteAddr().String())
			} else {
				log.Printf("Receive Failed: %v", err)
			}
			break
		}

		if count > 0 {
			data := buffer[:count]
			msg := InputPayload{}
			ans := ReturnPayload{}
			if err := json.Unmarshal(data, &msg); err != nil {
				log.Fatalln(err)
				return
			}
			switch msg.Menu {

			case "1":
				ans.Answer = convertTextToUpperCase(msg.Message)
			case "2":
				ans.Answer = fmt.Sprintf("IP= %s port=%s", tcpServerInfo.IP, tcpServerInfo.Port)
			case "3":
				t := getServerTime()
				ans.Answer = fmt.Sprintf("time = %s", t.Format("15:04:05"))
			case "4":
				d := getServerRunningTime()
				dt := d.Round(time.Second)

				h := dt / time.Hour
				dt -= h * time.Hour

				m := dt / time.Minute
				dt -= m * time.Minute

				s := dt / time.Second
				dt -= m * time.Second
				ans.Answer = fmt.Sprintf("run time = %02d:%02d:%02d", h, m, s)
			default:
				ans.Answer = ""
				ans.Error = "The menu is wrong"
			}
			b, _ := json.Marshal(ans)
			connection.Write(b)
		}
	}
}
