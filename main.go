package main

import (
	"bufio"
	"github.com/exgene/cuda-autocompletes/internal/rpc"
	"log"
	"os"
)

func main() {
	logger := getLogger("/home/kausthubh/GitHub/cuda-autocompletes/log.txt")
	scanner := bufio.NewScanner(os.Stdin)
	logger.Println("Started the main file")

	scanner.Split(rpc.Split)
	for scanner.Scan() {
		msg := scanner.Bytes()
		handleMessage(msg, logger)
	}
}

func handleMessage(msg []byte, logger *log.Logger) {
	logger.Printf("Received msg:%v", msg)
	method, content, err := rpc.DecodeMessage(msg)
	if err != nil {
		logger.Fatalf("Error while Decoding Message:\n %v", err)
		return
	}
	logger.Printf("Parsed Successfully! method: %s, content: %s", method, string(content))
}

func getLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("You fucked up")
	}
	return log.New(logFile, "[lsp==cuda]", log.Ldate|log.Ltime|log.Lshortfile)
}
