package main

import (
	"fmt"
	"simple_complier"
)

func main() {
	calculator := &front_end.SimpleCalculator{}

	script := "int a = b+3;"
	fmt.Println("解析变量声明语句: " + script)
	lexer := front_end.NewSimpleLexer()
	tokens := lexer.Tokenize(script)
	node := calculator.IntDeclare(tokens)
	front_end.DumpAST(node, " |---")

	//测试表达式
	script = "2+3*5;"
	fmt.Println("\n计算: " + script + "，看上去一切正常。")
	calculator.Evaluate(script)

	//测试语法错误
	script = "2+"
	fmt.Println("\n计算: " + script + "，应该有语法错误。")
	calculator.Evaluate(script)

	script = "2+3+4;"
	fmt.Println("\n计算: " + script + "，结合性出现错误。")
	calculator.Evaluate(script)

	script = "2*3*4;"
	fmt.Println("\n计算: " + script + "，结合性出现错误。")
	calculator.Evaluate(script)
}
