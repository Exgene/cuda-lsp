package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/exgene/cuda-autocompletes/internal/analysis"
	"github.com/exgene/cuda-autocompletes/internal/lsp"
	"github.com/exgene/cuda-autocompletes/internal/rpc"
)

func main() {
	logger := getLogger("/home/kausthubh/GitHub/cuda-autocompletes/log.txt")
	scanner := bufio.NewScanner(os.Stdin)
	state := analysis.NewState()

	logger.Println("Started the main file")

	scanner.Split(rpc.Split)
	for scanner.Scan() {
		msg := scanner.Bytes()
		handleMessage(msg, logger, state)
	}
}

func handleMessage(msg []byte, logger *log.Logger, state analysis.State) {
	logger.Printf("Received msg:%s", string(msg))
	method, content, err := rpc.DecodeMessage(msg)
	if err != nil {
		logger.Fatalf("Error while Decoding Message:\n %v", err)
	}

	logger.Printf("Recieved method: %s", method)

	switch method {
	case "initialize":
		var request lsp.IntializeRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Fatalf("Error while Decoding IntializeRequest: \n%v", err)
		}
		logger.Printf("Connected to %s, %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)
		writer := os.Stdout
		msg := lsp.NewIntializeResponse(request.ID)
		encodedMessage := rpc.EncodeMessage(msg)
		logger.Printf("Encoded Msg: %s", encodedMessage)
		writer.Write([]byte(encodedMessage))
		logger.Print("Written to StdOut")

	case "textDocument/didOpen":
		var document lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(content, &document); err != nil {
			logger.Fatalf("Error while Decoding IntializeRequest: \n%v", err)
		}
		logger.Printf("Document URI %s", document.Params.TextDocument.URI)
		logger.Printf("%s", document.Params.TextDocument.Text)
		state.OpenDocument(document.Params.TextDocument.URI, document.Params.TextDocument.Text)
	case "textDocument/didChange":
		var document lsp.DidChangeTextDocumentNotification
		if err := json.Unmarshal(content, &document); err != nil {
			logger.Fatalf("Error while Decoding IntializeRequest: \n%v", err)
		}
		for _, change := range document.Params.ContentChanges {
			logger.Printf("Recieved Diffs to update Document, Diffed Text :: %s", change.Text)
			diff := analysis.NewDiff()

			diff.StartRange.Character = change.Range.Start.Character
			diff.EndRange.Line = change.Range.Start.Line
			diff.StartRange.Character = change.Range.End.Character
			diff.EndRange.Line = change.Range.End.Line
			diff.Text = change.Text

			state.ApplyDiffs(document.Params.TextDocument.URI, diff)
		}
		for _, v := range state.Documents {
			logger.Printf("Diffs Updated, The new document state looks like: %s", v)
		}
	}
}

func getLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("You fucked up")
	}
	return log.New(logFile, "[lsp==cuda]", log.Ldate|log.Ltime|log.Lshortfile)
}
