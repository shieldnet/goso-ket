package TCPServer

import "time"

type TCPServerInfo struct {
	StartTime time.Time
	IP        string
	Port      string
}

type InputPayload struct {
	Menu    string `json:"menu"`
	Message string `json:"message"`
}

type ReturnPayload struct {
	Answer string `json:"answer"`
	Error  string `json:"error"`
}
