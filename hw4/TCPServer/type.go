package TCPServer

type Request struct {
	Command string            `json:"command"`
	Param   map[string]interface{} `json:"param"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"answer"`
	Items   map[string]interface{} `json:"items"`
}
