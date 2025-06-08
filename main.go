package main

import (
	"bufio"
	"fmt"
	"github.com/exgene/cuda-autocompletes/internal/rpc"
	"os"
)

func main() {
	fmt.Println("Started the main file")
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Split(rpc.Split)
	for scanner.Scan() {
		msg := scanner.Text()
		handleMessage(msg)
	}
}

func handleMessage(msg string) {
	fmt.Println(msg)
}
