package lsp

type HoverRequest struct {
	Request
	Params HoverParams `json:"params"`
}

type HoverParams struct {
	TextDocumentPositionParams
}

type HoverResponse struct {
	Response
	Result HoverResult `json:"result"`
}

type HoverResult struct {
	Contents string `json:"contents"`
}

func NewIntializeHoverResponse(requestId int, content string) HoverResponse {
	return HoverResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &requestId,
		},
		Result: HoverResult{
			Contents: content,
		},
	}
}
