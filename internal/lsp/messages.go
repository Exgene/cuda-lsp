package lsp

type RequestMessage struct {
	RPC    string `json:"jsonrpc"`
	ID     int    `json:"id"`
	Method string `json:"method"`
}

type ResponseMessage struct {
	RPC string `json:"jsonrpc"`
	ID  *int    `json:"id,omitempty"`
}

type Notification struct {
	RPC    string `json:"jsonrpc"`
	Method string `json:"method"`
}
