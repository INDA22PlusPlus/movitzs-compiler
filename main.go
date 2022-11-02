package main

import (
	"github.com/InDA22PlusPlus/movitzs-hw3/src/parser"
	"github.com/InDA22PlusPlus/movitzs-hw3/src/translator"
)

func main() {
	l := parser.New("print 1;")

	c := l.Program()

	translator.Translate(c)
}
