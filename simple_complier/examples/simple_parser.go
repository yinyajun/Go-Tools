package main

import (
	"fmt"
	"simple_complier"
)

func main() {
	parser := &front_end.SimpleParser{}
	script := "int age = 45+2; age= 20; age+10*2;"
	fmt.Println("解析：", script)
	tree := parser.Parse(script)
	front_end.DumpAST(tree, " |---")
	fmt.Println()

	script = "2+3+;"
	fmt.Println("解析：", script)
	tree = parser.Parse(script)
	fmt.Println()

	script = "2+3*;"
	fmt.Println("解析：", script)
	tree = parser.Parse(script)
	fmt.Println()
}
