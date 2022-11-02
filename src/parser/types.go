package parser

type CodeBlock struct {
	Statements []any
}

type LoopStatement struct {
	Statements []any
}

type BreakKeyword struct{}

type VarAss struct {
	Ident string
	Val   any
}

type IfStatement struct {
	Condition any
	Body      CodeBlock
}

type PrintStmt struct {
	Printee any
}

type OpType int

const (
	MulOp OpType = iota
	DivOp
	AddOp
	SubOp
)

type ArithmaticStatement struct {
	Op    OpType
	Left  any
	Right any
}

type (
	Identifier string
	Integer    int
)
