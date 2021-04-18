package BasicTCPServer

import (
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func convertTextToUpperCase(msg string) string {
	return strings.ToUpper(msg)
}

func getMyIPAddrAndPort() (ip string, port string) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatalln(err, "Are you connected to internet??")
	}

	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = localAddr.IP.String()
	port = strconv.Itoa(localAddr.Port)
	return ip, port
}

func getServerTime() time.Time {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		panic(err)
	}
	now := time.Now()
	t := now.In(loc)
	return t
}

func getServerRunningTime() time.Duration {
	st := getServerTime()
	diff := st.Sub(tcpServerInfo.StartTime)
	return diff
}
