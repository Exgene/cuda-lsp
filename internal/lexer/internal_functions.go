package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %v, Value: %v}", t.Kind, t.Value)
}

func isSkippable(char rune) bool {
	return unicode.IsSpace(char)
}

func isKeyword(word string) (TokenType, bool) {
	keyword, exists := keywords[word]
	return keyword, exists
}

func (tokenizer *Tokenizer) peek(ahead int) rune {
	if tokenizer.idx+ahead > len(tokenizer.src_code) {
		// IM USING THIS AS AN ASSERT UNTIL I FIGURE OUT GO!!!!
		fmt.Println("Shouldnt happen")
	}
	c, _ := utf8.DecodeRuneInString(tokenizer.src_code[tokenizer.idx+ahead-1:])
	return c
}

func (t *Tokenizer) isAtEnd() bool {
	if t.idx >= len(t.src_code) {
		return true
	}
	return false
}

func (t *Tokenizer) scanToken() error {
	for !t.isAtEnd() && isSkippable(t.peek(1)) {
		t.next()
	}

	if t.isAtEnd() {
		return nil
	}

	c := t.next()

	switch {
	case unicode.IsLetter(c):
		t.scanLiteral(c)
	}

	return nil
}

func (t *Tokenizer) addToTokensArray(token Token) {
	t.tokens = append(t.tokens, token)
}

func (t *Tokenizer) scanLiteral(c rune) {
	t.buf += string(c)
	for !t.isAtEnd() && (unicode.IsLetter(t.peek(1)) || unicode.IsNumber(t.peek(1)) || t.peek(1) == '_') {
		t.buf += string(t.next())
	}

	tokenType, isBufAKeyword := isKeyword(t.buf)

	if isBufAKeyword {
		t.addToTokensArray(Token{Kind: tokenType, Value: t.buf})
	}
	t.buf = ""
}

func (t *Tokenizer) next() rune {
	if t.isAtEnd() {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(t.src_code[t.idx:])
	t.idx += 1
	return r
}
