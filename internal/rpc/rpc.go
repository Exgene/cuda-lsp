package rpc

import (
	"encoding/json"
	"fmt"
)

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic("Cant encode message!" + err.Error())
	}

	formattedStr := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
	return formattedStr
}
