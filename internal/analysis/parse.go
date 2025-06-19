package analysis

import (
	"log"

	"github.com/exgene/cuda-autocompletes/internal/lexer"
)

func (s *State) NewParseJob() {
	parser := lexer.NewTokenizer(s.CurrentBuffer)
	tokens := parser.ScanTokens()
	for _, token := range tokens {
		log.Printf("Token:%v", token)
	}
}
