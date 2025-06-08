package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type BaseMessage struct {
	Method string `json:"method"`
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

	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return "", nil, errors.New("Invalid format :: Invalid contentLength value, Message:" + string(msg))
	}

	var baseMessage BaseMessage
	if err := json.Unmarshal(content[:contentLength], &baseMessage); err != nil {
		return "", nil, errors.New("Invalid format :: Invalid JSON Format, Message:" + string(msg))
	}

	return baseMessage.Method, content[:contentLength], nil
}

// type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)
func Split(msg []byte, _ bool) (int, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return 0, nil, nil
	}

	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return 0, nil, errors.New("Invalid format :: Invalid contentLength value, Message:" + string(msg))
	}

	if len(content) < contentLength {
		return 0, nil, nil
	}

	// +4 for the seperator
	totalBytesLength := len(header) + 4 + contentLength
	return totalBytesLength, msg[:totalBytesLength], nil
}
