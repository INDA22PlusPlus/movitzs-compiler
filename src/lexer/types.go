package lexer

type Token struct {
	TT  TokenType
	Val []byte
}

type TokenType int

const (
	Integer TokenType = iota
	MulOp
	DivOp
	AddOp
	SubOp
	If
	Identifier
	Equal
	Semicol
	Var
	LParen
	RParen
	LCracket
	RCracket
	Break
	Loop
	Print
	EOF
)
