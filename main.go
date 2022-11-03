package main

import (
	"github.com/InDA22PlusPlus/movitzs-hw3/src/parser"
	"github.com/InDA22PlusPlus/movitzs-hw3/src/translator"
)

// exempel på hur man gör broken kod: loop {var x = 1; break; var xx = 2;}

func main() {
	l := parser.New("var x = 666; loop { var x = x + x; }")

	c := l.Program()

	t := &translator.Translator{}
	t.Translate(c)
}
