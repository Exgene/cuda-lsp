package lexer

type TokenType int

const (
	Identifier TokenType = iota
	EOF
)

var keywords = map[string]TokenType{
	// Here i define my cuda tokens
}

type Token struct {
	Kind  TokenType
	Value string
}

type Tokenizer struct {
	idx      int
	buf      string
	tokens   []Token
	src_code string
}

func NewTokenizer(source string) *Tokenizer {
	return &Tokenizer{
		src_code: source,
		idx:      0,
		tokens:   make([]Token, 0),
		buf:      "",
	}
}
