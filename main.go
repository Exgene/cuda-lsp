package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/exgene/cuda-autocompletes/internal/analysis"
	"github.com/exgene/cuda-autocompletes/internal/lsp"
	"github.com/exgene/cuda-autocompletes/internal/rpc"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	state := analysis.NewState()
	writer := os.Stdout

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
		handleMessage(msg, state, writer)
	}
	defer logFile.Close()
}

func writeMessage(msg any, writer io.Writer) {
	encodedMessage := rpc.EncodeMessage(msg)
	writer.Write([]byte(encodedMessage))
	log.Printf("Encoded Msg: %s", encodedMessage)
}

func handleMessage(msg []byte, state analysis.State, writer io.Writer) {
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
			log.Fatalf("Error unmarshalling intialize: %v", err)
		}
		log.Printf("Connected to %s, %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)
		msg := lsp.NewIntializeResponse(request.ID)
		writeMessage(msg, writer)
		log.Print("Written to StdOut")

	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(content, &request); err != nil {
			log.Fatalf("Error unmarshalling textDocument/hover: %v", err)
		}
		log.Printf("Recieved Character Position: %d", request.Params.Position.Character)
		log.Printf("Recieved Line Position: %d", request.Params.Position.Line)
		token := state.GetToken(request.Params.TextDocument.URI, request.Params.Position)
		log.Printf("Token Hovered: %s", token)
		// log.Print("THIS IS SO FUCKING HARD THAN IT NEEDS TO BE ARGH")
		// Done ToT

		/*
			Ideally you would have some sort of dicitionary of all the available tokens parsed through parseJob()
			Then you have its relative position encoded. Then when hovered you round it to the relative position
			If its an interested token you directly return its information. But I Aint doin allat type shit. Tbh i can
			But i wanna learn CUDA. Thats the whole point of this thing.
		*/

		t := state.ValidToken(token, request.Params.TextDocument.URI)
		if t == nil {
			log.Printf("No Token Found so no hover for you")
			return
		}

		msg := lsp.NewIntializeHoverResponse(request.Request.ID, t.Kind.String())
		writeMessage(msg, writer)
		log.Print("Written to StdOut")

	case "textDocument/didOpen":
		var document lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(content, &document); err != nil {
			log.Fatalf("Error unmarshalling textDocument/didOpen: %v", err)
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
			log.Fatalf("Error unmarshalling textDocument/didChange: %v", err)
		}
		for _, change := range document.Params.ContentChanges {
			log.Printf("<----------------------->")
			log.Printf("Recieved Diffs to update Document, Diffed Text :: %s", change.Text)
			log.Printf("Recieved Diffs to update Document, Diffed Start Line :: %d", change.Range.Start.Line)
			log.Printf("Recieved Diffs to update Document, Diffed Start Character :: %d", change.Range.Start.Character)
			log.Printf("Recieved Diffs to update Document, Diffed End Line :: %d", change.Range.End.Line)
			log.Printf("Recieved Diffs to update Document, Diffed End Character :: %d", change.Range.End.Character)
			log.Printf("<----------------------->")
			diff := analysis.NewDiff(change)

			log.Printf("Starting to apply diffs")
			err := state.ApplyDiffs(document.Params.TextDocument.URI, diff)
			if err != nil {
				log.Printf("Error while Applying Diffs %v", err)
				break
			}

			state.Tokens[document.Params.TextDocument.URI] = state.NewParseJob()
		}
		for _, values := range state.Documents[document.Params.TextDocument.URI] {
			log.Printf("%s\n", values)
		}
	}
}
