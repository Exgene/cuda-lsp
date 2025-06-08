package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic("Cant encode message!" + err.Error())
	}

	formattedStr := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
	return formattedStr
}

func DecodeMessage(msg []byte) (int, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return 0, errors.New("Invalid format :: Recieved no separator, Message:" + string(msg))
	}

	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return 0, errors.New("Invalid format :: Invalid contentLength value, Message:" + string(msg))
	}

	_ = content
	return contentLength, nil
}
