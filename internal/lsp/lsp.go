package lsp

type IntializeRequest struct {
	Request
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
	Response
	Result IntializeResult `json:"result"`
}

type IntializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

type ServerCapabilities struct {
	TextDocumentSync int  `json:"textDocumentSync"`
	HoverProvider    bool `json:"hoverProvider"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewIntializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: IntializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync: 2,
				HoverProvider:    true,
			},
			ServerInfo: ServerInfo{
				Name:    "cuda",
				Version: "0.0.1",
			},
		},
	}

}
