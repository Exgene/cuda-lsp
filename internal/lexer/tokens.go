package lexer

type CudaTokenType int

const (
	CUDA CudaTokenType = iota
	Identifier
	EOF
	Global
	Device
	ThreadIdx
	BlockIdx
	BlockDim
	GridDim
)

var keywords = map[string]CudaTokenType{
	// Here i define my cuda tokens
	"cuda":       CUDA,
	"__global__": Global,
	"__device__": Device,
	"threadIdx":  ThreadIdx,
	"blockIdx":   BlockIdx,
	"blockDim":   BlockDim,
	"gridDim":    GridDim,
}

type Token struct {
	Kind  CudaTokenType
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
