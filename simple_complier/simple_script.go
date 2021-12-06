package front_end

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"simple_complier/astNode"
)

/**
 * 一个简单的脚本解释器。
 * 所支持的语法，请参见SimpleParser.java
 *
 * 运行脚本：
 * 在命令行下，键入：java SimpleScript
 * 则进入一个REPL界面。你可以依次敲入命令。比如：
 * > 2+3;
 * > int age = 10;
 * > int b;
 * > b = 10*2;
 * > age = age + b;
 *
 * 你还可以使用一个参数 -v，让每次执行脚本的时候，都输出AST和整个计算过程。
 *
 */

type SimpleScript struct {
	variables map[string]interface{}
	verbose   bool
}

func NewSimpleScript() *SimpleScript {
	s := &SimpleScript{}
	s.variables = make(map[string]interface{})
	return s
}

type Reader struct {
	*bufio.Scanner
}

func NewReader(r io.Reader) Reader {
	m := Reader{bufio.NewScanner(r)}
	m.Scanner.Split(bufio.ScanLines)
	return m
}

func (r Reader) ReadLine() (string, error) {
	if r.Scan() {
		return strings.Trim(r.Text(), " "), nil
	} else {
		return "", io.EOF
	}
}

func (s *SimpleScript) REPL() {
	flag.BoolVar(&s.verbose, "v", false, "verbose mode")
	if s.verbose {
		fmt.Println("verbose mode")
	}
	fmt.Println("Simple script language!")

	parser := &SimpleParser{}
	reader := NewReader(os.Stdin)

	scriptText := ""
	fmt.Print("\n>")

	for {
		line, err := reader.ReadLine()
		if err != nil {
			break
		}
		scriptText += line + "\n"
		if len(line) > 0 && line[len(line)-1:] == ";" {
			tree := parser.Parse(scriptText)
			if s.verbose {
				DumpAST(tree, " |---")
			}
			s.process(tree, "")
			fmt.Print("\n>")
			scriptText = ""
		}
	}
}

func (s *SimpleScript) process(node astNode.ASTNode, indent string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	s.evaluate(node, indent)
}

func (s *SimpleScript) evaluate(node astNode.ASTNode, indent string) interface{} {
	var result interface{}
	if s.verbose {
		fmt.Println(indent+"Calculating:", astNode.GetAstNodeTypeName(node.GetType()))
	}
	switch node.GetType() {
	case astNode.Programm:
		for _, child := range node.GetChildren() {
			result = s.evaluate(child, indent)
		}
		break
	case astNode.Additive:
		child1 := node.GetChildren()[0]
		value1 := s.evaluate(child1, indent+"\t").(int)
		child2 := node.GetChildren()[1]
		value2 := s.evaluate(child2, indent+"\t").(int)
		if node.GetText() == "+" {
			result = value1 + value2
		} else {
			result = value1 - value2
		}
		break
	case astNode.Multiplicative:
		child1 := node.GetChildren()[0]
		value1 := s.evaluate(child1, indent+"\t").(int)
		child2 := node.GetChildren()[1]
		value2 := s.evaluate(child2, indent+"\t").(int)
		if node.GetText() == "*" {
			result = value1 * value2
		} else {
			result = value1 / value2
		}
		break
	case astNode.IntLiteral:
		result, _ = strconv.Atoi(node.GetText())
		break
	case astNode.Identifier:
		varName := node.GetText()
		value, ok := s.variables[varName]
		if ok {
			if value != nil {
				result = value.(int)
			} else {
				panic("variable " + varName + " has not been set any value.")
			}
		} else {
			panic("unknown variable:" + varName)
		}
		break
	case astNode.AssignmentStmt:
		varName := node.GetText()
		_, ok := s.variables[varName]
		if !ok {
			panic("unknown variable: " + varName)
		}
		var varValue interface{}
		if len(node.GetChildren()) > 0 {
			child := node.GetChildren()[0]
			result = s.evaluate(child, indent+"\t")
			varValue = result
		}
		s.variables[varName] = varValue
		break
	case astNode.IntDeclaration:
		varName := node.GetText()
		var varValue interface{}
		if len(node.GetChildren()) > 0 {
			child := node.GetChildren()[0]
			result = s.evaluate(child, indent+"\t")
			varValue = result
		}
		s.variables[varName] = varValue
		break
	default:
	}
	if s.verbose {
		fmt.Println(indent+"Result:", result)
	} else if indent == "" {
		if node.GetType() == astNode.IntDeclaration || node.GetType() == astNode.AssignmentStmt {
			fmt.Println(node.GetText()+":", result)
		} else if node.GetType() != astNode.Programm {
			fmt.Println(result)
		}
	}
	return result
}
