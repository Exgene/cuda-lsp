package encode_test

import (
	"fmt"
	"testing"

	"github.com/exgene/cuda-autocompletes/internal/rpc"
)

type EncodingExample struct {
	Testing bool
}

func TestEncode(t *testing.T) {
	expected := "Content-Length: 16\r\n\r\n{\"Testing\":true}"
	actual := rpc.EncodeMessage(EncodingExample{Testing: true})

	if expected != actual {
		t.Fatalf("Encoded values dont match! Original ::: %s === Actual ::: %s", expected, actual)
	}
}

func TestDecode(t *testing.T) {
	incomingMessage := "Content-Length: 16\r\n\r\n{\"method\":\"hey\"}"
	method, content, err := rpc.DecodeMessage([]byte(incomingMessage))

	expectedLength := 16
	expectedMethodValue := "hey"

	contentLength := len(content)
	if err != nil {
		t.Fatalf("Error while decoding message....%v", err)
	}
	if expectedLength != contentLength {
		t.Fatalf("Decoded Length dont match! Original ::: %d === Actual ::: %d", expectedLength, contentLength)
	}
	fmt.Printf("Method:%s", method)
	if expectedMethodValue != method {
		t.Fatalf("Decoded Method Value dont match! Original ::: %s === Actual ::: %s", expectedMethodValue, method)
	}
}
