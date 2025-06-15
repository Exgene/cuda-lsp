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
	scanner := bufio.NewScanner(os.Stdin)
	state := analysis.NewState()

	logFile, err := os.OpenFile("/home/kausthubh/GitHub/cuda-autocompletes/log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[lsp==cuda]")

	log.Println("Started the main file")

	scanner.Split(rpc.Split)
	for scanner.Scan() {
		msg := scanner.Bytes()
		handleMessage(msg, state)
	}
	defer logFile.Close()
}

func handleMessage(msg []byte, state analysis.State) {
	// logger.Printf("Received msg:%s", string(msg))
	method, content, err := rpc.DecodeMessage(msg)
	if err != nil {
		log.Fatalf("Error while Decoding Message:\n %v", err)
	}

	log.Printf("Recieved method: %s", method)

	switch method {
	case "initialize":
		var request lsp.IntializeRequest
		if err := json.Unmarshal(content, &request); err != nil {
			log.Fatalf("Error while Decoding IntializeRequest: \n%v", err)
		}
		log.Printf("Connected to %s, %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)
		writer := os.Stdout
		msg := lsp.NewIntializeResponse(request.ID)
		encodedMessage := rpc.EncodeMessage(msg)
		log.Printf("Encoded Msg: %s", encodedMessage)
		writer.Write([]byte(encodedMessage))
		log.Print("Written to StdOut")

	case "textDocument/didOpen":
		var document lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(content, &document); err != nil {
			log.Fatalf("Error while Decoding IntializeRequest: \n%v", err)
		}
		log.Printf("Document URI %s", document.Params.TextDocument.URI)
		// log.Printf("%s", document.Params.TextDocument.Text)
		state.OpenDocument(document.Params.TextDocument.URI, document.Params.TextDocument.Text)
		log.Printf("Showing State:::")
		for _, values := range state.Documents[document.Params.TextDocument.URI] {
			log.Printf("%s", values)
		}

	case "textDocument/didChange":
		var document lsp.DidChangeTextDocumentNotification
		if err := json.Unmarshal(content, &document); err != nil {
			log.Fatalf("Error while Decoding IntializeRequest: \n%v", err)
		}
		for _, change := range document.Params.ContentChanges {
			log.Printf("<----------------------->")
			log.Printf("Recieved Diffs to update Document, Diffed Text :: %s", change.Text)
			log.Printf("Recieved Diffs to update Document, Diffed Start Line :: %d", change.Range.Start.Line)
			log.Printf("Recieved Diffs to update Document, Diffed Start Character :: %d", change.Range.Start.Character)
			log.Printf("Recieved Diffs to update Document, Diffed End Line :: %d", change.Range.End.Line)
			log.Printf("Recieved Diffs to update Document, Diffed End Character :: %d", change.Range.End.Character)
			log.Printf("<----------------------->")
			diff := analysis.NewDiff()

			diff.StartRange.Line = change.Range.Start.Line
			diff.StartRange.Character = change.Range.Start.Character
			diff.EndRange.Line = change.Range.End.Line
			diff.EndRange.Character = change.Range.End.Character
			diff.Text = change.Text

			// log.Printf("Starting to apply diffs")
			// err := state.ApplyDiffs(document.Params.TextDocument.URI, diff)
			// if err != nil {
			// 	log.Printf("Error while Applying Diffs %v", err)
			// 	break
			// }
			// for _, values := range state.Documents[document.Params.TextDocument.URI] {
			// 	fmt.Printf("%s\n", values)
			// }
		}
		// for _, v := range state.Documents[document.Params.TextDocument.URI] {
		// 	log.Printf("Diffs Updated, The new document state looks like: %s", v)
		// }
	}
}
