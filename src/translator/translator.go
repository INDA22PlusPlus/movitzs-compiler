package translator

import (
	"fmt"
	"strconv"

	"github.com/InDA22PlusPlus/movitzs-hw3/src/parser"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Translator struct {
	curBlock *ir.Block
	breakTo  *ir.Block
	vartable map[string]value.Value
}

var printfCall *ir.Func

func (t *Translator) Translate(translatee parser.CodeBlock) {
	t.vartable = map[string]value.Value{}

	m := ir.NewModule()
	printfCall = m.NewFunc("printf", types.I32, &ir.Param{Typ: types.I8Ptr})
	printfCall.Sig.Variadic = true

	//	m.NewFunc("__libc_start_main", types.I32, ir.NewParam("main", &types.FuncType{}))

	f := m.NewFunc("_start", types.Void)

	t.curBlock = f.NewBlock("mainbl")
	t.codeblock(translatee)
	t.curBlock.NewRet(nil)

	fmt.Println(m.String())
}

func (t *Translator) codeblock(x parser.CodeBlock) {
	for _, stmt := range x.Statements {
		switch stmt := stmt.(type) {
		case parser.CodeBlock:
			t.codeblock(stmt)

		case parser.IfStatement:
			t.ifstmt(stmt)

		case parser.VarAss:
			varr := t.curBlock.NewAdd(constant.NewInt(types.I64, 0), t.val(stmt.Val))

			varr.SetName(stmt.Ident)
			t.vartable[stmt.Ident] = varr

		case parser.LoopStatement:

			loopBlock := t.curBlock.Parent.NewBlock(t.name("loop_body"))

			t.curBlock.NewBr(loopBlock)

			breakTo := t.curBlock.Parent.NewBlock(t.name("break_to"))
			t.breakTo = breakTo
			t.curBlock = loopBlock
			t.codeblock(parser.CodeBlock(stmt))

			if t.curBlock.Term == nil {
				t.curBlock.NewBr(loopBlock)
			}
			t.curBlock = breakTo

		case parser.BreakKeyword:
			t.curBlock.NewBr(t.breakTo)
			t.breakTo = nil

		case parser.PrintStmt:

			y := constant.NewCharArrayFromString("%d\n" + string([]byte{0}))
			typex := types.NewPointer(types.NewArray(uint64(5), types.I8))

			// zero := constant.NewInt(types.I32, 0)
			x := t.curBlock.NewGetElementPtr(typex, y)

			// wtf do i do

			t.curBlock.NewCall(printfCall, x, t.expr(stmt.Printee))
		default:
			fmt.Printf("could not find type %T\n", stmt)
			panic("unknown block")
		}
	}
}

func (t *Translator) ifstmt(x parser.IfStatement) {
	metaBlock := t.curBlock.Parent.NewBlock(t.name("if_meta"))
	t.curBlock.NewBr(metaBlock)

	ifBody := t.curBlock.Parent.NewBlock(t.name("if_body"))
	postIf := t.curBlock.Parent.NewBlock(t.name("post_if"))

	t.curBlock = metaBlock
	zero := constant.NewInt(types.I64, 0)
	cond := metaBlock.NewICmp(enum.IPredEQ, zero, t.expr(x.Condition))

	metaBlock.NewCondBr(cond, ifBody, postIf)

	t.curBlock = ifBody
	ifBody.NewBr(postIf) // might be overridden by break
	t.codeblock(x.Body)

	t.curBlock = postIf
}

func (t *Translator) val(x any) value.Value {
	switch l := x.(type) {
	case parser.ArithmaticStatement:
		return t.expr(l)
	case parser.Identifier:
		return t.vartable[string(l)]
	case parser.Integer:
		return constant.NewInt(types.I64, int64(l))

	}
	fmt.Printf("%T\n", x)
	panic("nah")
}

func (t *Translator) expr(xx any) value.Value {
	switch tt := xx.(type) {
	case parser.Integer:
		return t.curBlock.NewAdd(constant.NewInt(types.I64, 0), constant.NewInt(types.I64, int64(tt)))
	case parser.Identifier:
		return t.curBlock.NewAdd(constant.NewInt(types.I64, 0), t.vartable[string(tt)])
	}

	x := xx.(parser.ArithmaticStatement)

	left := t.val(x.Left)
	right := t.val(x.Right)

	switch x.Op {
	case parser.AddOp:
		return t.curBlock.NewAdd(left, right)
	case parser.SubOp:
		return t.curBlock.NewSub(left, right)
	case parser.MulOp:
		return t.curBlock.NewMul(left, right)
	case parser.DivOp:
		return t.curBlock.NewSDiv(left, right) // is this correct?
	}
	panic("no")
}

func (t *Translator) name(xx string) string {
	return xx + "_" + strconv.Itoa(len(t.curBlock.Parent.Blocks))
}
