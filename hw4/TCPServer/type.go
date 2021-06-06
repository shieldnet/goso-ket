package TCPServer

type Request struct {
	Command string `json:"menu"`
	Parameter []string `json:"message"`
}

type Response struct {
	Status string `json:"status"`
	Message string `json:"answer"`
}
