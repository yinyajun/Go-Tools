package main

import (
	"fmt"
	"simple_complier"
)

func main() {
	lexer := front_end.NewSimpleLexer()

	script := "int age =  45;"
	fmt.Println("Tokenize:", script)
	tokenReader := lexer.Tokenize(script)
	lexer.Dump(tokenReader)
	fmt.Println()

	script = "inta age =  45;"
	fmt.Println("Tokenize:", script)
	tokenReader = lexer.Tokenize(script)
	lexer.Dump(tokenReader)
	fmt.Println()

	script = "in age = 45;"
	fmt.Println("Tokenize:", script)
	tokenReader = lexer.Tokenize(script)
	lexer.Dump(tokenReader)
	fmt.Println()

	script = "age >= 45;"
	fmt.Println("Tokenize:", script)
	tokenReader = lexer.Tokenize(script)
	lexer.Dump(tokenReader)
	fmt.Println()

	script = "age > 45;"
	fmt.Println("Tokenize:", script)
	tokenReader = lexer.Tokenize(script)
	lexer.Dump(tokenReader)
	fmt.Println()
}
