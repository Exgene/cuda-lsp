package encode_test

import (
	"github.com/exgene/cuda-autocompletes/internal/rpc"
	"testing"
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
