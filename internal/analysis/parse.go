package analysis

import (
	"log"

	"github.com/exgene/cuda-autocompletes/internal/lexer"
)

func (s *State) NewParseJob() []lexer.Token {
	log.Printf("Current Buffer State: %s", s.CurrentBuffer)
	parser := lexer.NewTokenizer(s.CurrentBuffer)
	log.Println("Scanning Toknes!")
	tokens := parser.ScanTokens()
	log.Printf("Length of tokens recieved: %d", len(tokens))
	for _, token := range tokens {
		log.Printf("Token Type : %v", token.Kind)
		log.Printf("Token Value : %s", token.Value)
	}
	return tokens
}
