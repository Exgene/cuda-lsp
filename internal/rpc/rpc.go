package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type BaseMessage struct {
	method string `json:"string"`
}

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic("Cant encode message!" + err.Error())
	}

	formattedStr := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
	return formattedStr
}

func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return "", nil, errors.New("Invalid format :: Recieved no separator, Message:" + string(msg))
	}

	// Try to decode just for validation
	contentLengthBytes := header[len("Content-Length: "):]
	_, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return "", nil, errors.New("Invalid format :: Invalid contentLength value, Message:" + string(msg))
	}

	var baseMessage BaseMessage
	if err := json.Unmarshal(content, &baseMessage); err != nil {
		return "", nil, errors.New("Invalid format :: Invalid JSON Format, Message:" + string(msg))
	}

	return baseMessage.method, contentLengthBytes, nil
}
