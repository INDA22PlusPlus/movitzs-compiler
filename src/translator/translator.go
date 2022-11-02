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

var vartable map[string]value.Value

var printfCall *ir.Func

func Translate(translatee parser.CodeBlock) {
	vartable = map[string]value.Value{}

	m := ir.NewModule()
	printfCall = m.NewFunc("printf", types.I32, &ir.Param{Typ: types.I8Ptr})
	printfCall.Sig.Variadic = true

	f := m.NewFunc("_start", types.Void)

	block := f.NewBlock("mainbl")

	codeblock(&block, translatee)

	block.NewRet(nil)

	fmt.Println(m.String())
}

func codeblock(block **ir.Block, x parser.CodeBlock) *ir.Block {
	//
	for _, stmt := range x.Statements {
		switch stmt := stmt.(type) {
		case parser.CodeBlock:
			codeblock(block, stmt)

		case parser.IfStatement:
			postIf := (*block).Parent.NewBlock("post_if_" + strconv.Itoa(len((*block).Parent.Blocks)))
			metaBlock := ifstmt((*block), postIf, stmt)
			(*block).NewBr(metaBlock)
			(*block) = postIf

		case parser.VarAss:
			varr := (*block).NewAdd(constant.NewInt(types.I64, 0), val((*block), stmt.Val))

			varr.SetName(stmt.Ident)
			vartable[stmt.Ident] = varr

		case parser.LoopStatement:
			/*
				lbb := (*block).Parent.NewBlock("loopBody")
				codeblock(&lbb, parser.CodeBlock(stmt))

				loopBlock := (*block).Parent.NewBlock("loop")
				loopBlock.NewBr(lbb)
				loopBlock.NewBr(loopBlock)

				(*block).NewBr(loopBlock)
			*/

		case parser.BreakKeyword:
			// todo
		case parser.PrintStmt:

			y := constant.NewCharArrayFromString("%d\n" + string([]byte{0}))
			typex := types.NewPointer(types.NewArray(uint64(5), types.I8))

			// zero := constant.NewInt(types.I32, 0)
			x := (*block).NewGetElementPtr(typex, y)

			// wtf do i do

			(*block).NewCall(printfCall, x, expr((*block), stmt.Printee))
		default:
			fmt.Printf("could not find type %T\n", stmt)
			panic("unknown block")
		}
	}

	return *block
}

func ifstmt(parent *ir.Block, postIf *ir.Block, x parser.IfStatement) *ir.Block {
	metaBlock := parent.Parent.NewBlock("if_meta_" + strconv.Itoa(len(parent.Parent.Blocks)))

	ifBody := parent.Parent.NewBlock("if_body_" + strconv.Itoa(len(parent.Parent.Blocks)))
	codeblock(&ifBody, x.Body)
	ifBody.NewBr(postIf)

	zero := constant.NewInt(types.I64, 0)
	cond := metaBlock.NewICmp(enum.IPredEQ, zero, expr(metaBlock, x.Condition))

	metaBlock.NewCondBr(cond, ifBody, postIf)

	return metaBlock
}

func val(block *ir.Block, x any) value.Value {
	switch l := x.(type) {
	case parser.ArithmaticStatement:
		return expr(block, l)
	case parser.Identifier:
		return vartable[string(l)]
	case parser.Integer:
		return constant.NewInt(types.I64, int64(l))

	}
	fmt.Printf("%T\n", x)
	panic("nah")
}

func expr(block *ir.Block, xx any) value.Value {
	switch t := xx.(type) {
	case parser.Integer:
		return block.NewAdd(constant.NewInt(types.I64, 0), constant.NewInt(types.I64, int64(t)))
	}

	x := xx.(parser.ArithmaticStatement)

	left := val(block, x.Left)
	right := val(block, x.Right)

	switch x.Op {
	case parser.AddOp:
		return block.NewAdd(left, right)
	case parser.SubOp:
		return block.NewSub(left, right)
	case parser.MulOp:
		return block.NewMul(left, right)
	case parser.DivOp:
		return block.NewSDiv(left, right) // is this correct?
	}
	panic("no")
}
