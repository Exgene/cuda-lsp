package lsp

type IntializeRequest struct {
	RequestMessage
	Params IntializeParams `json:"params"`
}

type IntializeParams struct {
	ClientInfo *ClientInfo `json:"clientInfo,omitempty"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResponse struct {
	ResponseMessage
	Result IntializeResult `json:"result"`
}

type IntializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

type ServerCapabilities struct {
	TextDocumentSync int `json:"textDocumentSync"`
}
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewIntializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		ResponseMessage: ResponseMessage{
			RPC: "2.0",
			ID:  &id,
		},
		Result: IntializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync: 1,
			},
			ServerInfo: ServerInfo{
				Name:    "cuda",
				Version: "0.0.1",
			},
		},
	}

}
