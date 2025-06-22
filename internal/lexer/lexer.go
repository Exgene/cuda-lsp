package lexer

import (
	"fmt"
)

func (t TokenType) String() string {
	switch t {
	case Identifier:
		return "Identifier"
	case EOF:
		return "EOF"
	case CUDA:
		return "Cuda"
	default:
		return fmt.Sprintf("Unknown(%d)", int(t))
	}
}

func Tokenize(code string) {
	tokenizer := NewTokenizer(code)
	tokens := tokenizer.ScanTokens()
	for _, token := range tokens {
		fmt.Println(token.String())
	}
}

func (t *Tokenizer) ScanTokens() []Token {
	for !t.isAtEnd() {
		err := t.scanToken()
		if err != nil {
			fmt.Printf("Error while parsing tokens...%v", err)
		}
	}

	t.addToTokensArray(Token{TokenType(EOF), ""})
	return t.tokens
}
