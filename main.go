package main

import (
	"encoding/json"
	"fmt"

	"github.com/InDA22PlusPlus/movitzs-hw3/src/parser"
)

func main() {
	println("here")
	l := parser.New("loop break; loop { var x = 1; break; }")

	c := l.Program()

	x, _ := json.MarshalIndent(c, " ", " ")
	fmt.Printf("%s\n", x)
	fmt.Printf("%+#v\n", c.Statements)
}
