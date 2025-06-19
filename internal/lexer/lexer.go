package lexer

import (
	"fmt"
)

func NewTokenizer(source string) *Tokenizer {
	return &Tokenizer{
		src_code: source,
		idx:      0,
		tokens:   make([]Token, 0),
		buf:      "",
	}
}

func (t TokenType) String() string {
	switch t {
	case Identifier:
		return "Identifier"
	default:
		return fmt.Sprintf("Unknown(%d)", int(t))
	}
}

func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %v, Value: %v}", t.Kind, t.Value)
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
