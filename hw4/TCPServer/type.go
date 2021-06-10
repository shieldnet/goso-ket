package TCPServer

type Request struct {
	Command string            `json:"command"`
	Param   map[string]string `json:"param"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"answer"`
	Items   map[string]interface{} `json:"items"`
}
