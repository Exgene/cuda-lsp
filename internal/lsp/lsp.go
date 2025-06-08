package lsp

type IntializeRequest struct {
	Request RequestMessage
	Params  IntializeParams `json:"params"`
}

type IntializeParams struct {
	ClientInfo *ClientInfo `json:"clientInfo,omitempty"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
