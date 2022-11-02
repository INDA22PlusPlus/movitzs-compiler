package lexer

import "regexp"

type Lexer struct {
	str    string
	cursor int
}

func New(str string) *Lexer {
	return &Lexer{
		str:    str,
		cursor: 0,
	}
}

type rule struct {
	tt  TokenType
	reg *regexp.Regexp
}

var rules = []rule{
	{Integer, regexp.MustCompile(`^\d+`)},
	{AddOp, regexp.MustCompile(`^\+`)},
	{SubOp, regexp.MustCompile(`^\-`)},
	{MulOp, regexp.MustCompile(`^\*`)},
	{DivOp, regexp.MustCompile(`^\/`)},
	{If, regexp.MustCompile(`^if`)},
	{Equal, regexp.MustCompile(`^=`)},
	{Semicol, regexp.MustCompile(`^;`)},
	{Var, regexp.MustCompile(`^var`)},
	{LParen, regexp.MustCompile(`^\(`)},
	{RParen, regexp.MustCompile(`^\)`)},
	{LCracket, regexp.MustCompile(`^\{`)},
	{RCracket, regexp.MustCompile(`^\}`)},
	{Print, regexp.MustCompile((`^print`))},

	{Loop, regexp.MustCompile(`^loop`)},
	{Break, regexp.MustCompile(`^break`)},
	{Identifier, regexp.MustCompile(`^[a-zA-Z]\w*`)},

	{-1, regexp.MustCompile(`^\s`)},
}

func (l *Lexer) HasNext() bool {
	return l.cursor < len(l.str)
}

func (l *Lexer) NextToken() Token {
begin:
	if l.cursor == len(l.str) {
		return Token{TT: EOF}
	}

	for _, rul := range rules {

		match := rul.reg.Find([]byte(l.str[l.cursor:]))
		if match == nil {
			continue
		}

		l.cursor += len(match)

		if rul.tt == -1 { // skip whitespace
			goto begin
		}

		return Token{
			TT:  rul.tt,
			Val: match,
		}
	}

	panic("no token recognized, val: (" + l.str[l.cursor:] + ")")
}
