package parser

import (
	"fmt"
	"strconv"

	"github.com/InDA22PlusPlus/movitzs-hw3/src/lexer"
)

type Parser struct {
	l    *lexer.Lexer
	peek lexer.Token
}

func New(str string) *Parser {
	x := lexer.New(str)
	return &Parser{
		l:    x,
		peek: x.NextToken(),
	}
}

func (p *Parser) nextToken() lexer.Token {
	old := p.peek
	p.peek = p.l.NextToken()
	return old
}

func (p *Parser) Program() CodeBlock {
	stmts := make([]any, 0, 420)
	for p.peek.TT != lexer.EOF {
		stmts = append(stmts, p.codeblock())
	}

	return CodeBlock{
		Statements: stmts,
	}
}

func (p *Parser) codeblock() CodeBlock {
	pp := CodeBlock{
		Statements: make([]any, 0),
	}

	if p.peek.TT == lexer.LCracket {
		p.expect(lexer.LCracket)
		for p.l.HasNext() && p.peek.TT != lexer.RCracket {
			pp.Statements = append(pp.Statements, p.stmt())
		}
		p.expect(lexer.RCracket)
	} else {
		pp.Statements = append(pp.Statements, p.stmt())
	}

	return pp
}

func (p *Parser) stmt() any {
	switch p.peek.TT {
	case lexer.Print:
		return p.print()
	case lexer.Var:
		return p.varass()
	case lexer.If:
		return p.if_stmt()
	case lexer.Loop:
		p.expect(lexer.Loop)
		return LoopStatement{
			Statements: p.codeblock().Statements,
		}
	case lexer.Break:
		p.expect(lexer.Break)
		p.expect(lexer.Semicol)

		return BreakKeyword{}

	default:
		panic("unknown token: " + fmt.Sprint(p.peek.TT))
	}
}

func (p *Parser) print() PrintStmt {
	p.expect(lexer.Print)
	x := PrintStmt{Printee: p.expr()}
	p.expect(lexer.Semicol)
	return x
}

func (p *Parser) if_stmt() IfStatement {
	p.expect(lexer.If)
	p.expect(lexer.LParen)
	e := p.expr()
	p.expect(lexer.RParen)
	return IfStatement{
		Condition: e,
		Body:      p.codeblock(),
	}
}

func (p *Parser) varass() VarAss {
	p.expect(lexer.Var)

	id := p.expect(lexer.Identifier)

	p.expect(lexer.Equal)

	expr := p.expr()

	p.expect(lexer.Semicol)
	return VarAss{
		Ident: string(id.Val),
		Val:   expr,
	}
}

func (p *Parser) expr() any {
	if p.peek.TT == lexer.LParen {
		p.expect(lexer.LParen)
		defer p.expect(lexer.RParen)
	}

	return p.addt_expr()
}

func (p *Parser) addt_expr() any {
	left := p.mult_expr()

	if p.peek.TT != lexer.AddOp && p.peek.TT != lexer.SubOp {
		return left
	}

	op := p.nextToken()

	right := p.addt_expr()

	return ArithmaticStatement{
		Op:    OpType(op.TT - 1),
		Left:  left,
		Right: right,
	}
}

func (p *Parser) mult_expr() any {
	var left any = nil
	var right any = nil

	left = p.unry_expr()

	if p.peek.TT != lexer.DivOp && p.peek.TT != lexer.MulOp {
		return left
	}

	op := OpType(p.nextToken().TT - 1) // lol

	right = p.unry_expr()

	return ArithmaticStatement{
		Op:    op,
		Left:  left,
		Right: right,
	}
}

func (p *Parser) unry_expr() any {
	if p.peek.TT == lexer.LParen {
		p.expect(lexer.LParen)
		res := p.expr()
		p.expect(lexer.RParen)
		return res
	}

	if p.peek.TT == lexer.Identifier {
		return Identifier(p.nextToken().Val)
	}

	if p.peek.TT == lexer.Integer {
		x, _ := strconv.Atoi(string(p.nextToken().Val))
		return Integer(x)
	}

	if p.peek.TT == lexer.SubOp {
		p.nextToken()
		return ArithmaticStatement{
			Op:    SubOp,
			Left:  Integer(0),
			Right: p.expr(),
		}
	}

	panic("eh")
}

func (p *Parser) expect(tt lexer.TokenType) lexer.Token {
	t := p.nextToken()
	if t.TT != tt {
		panic(fmt.Sprintf("ayowtf, expected %+v got %+v\n", tt, t.TT))
	}
	return t
}
