package front_end

import (
	"fmt"
	"simple_complier/astNode"
	tok "simple_complier/token"
	"strconv"
)

type SimpleCalculator struct{}

func (c *SimpleCalculator) Evaluate(script string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	tree := c.Parse(script)
	DumpAST(tree, " |---")
	c.evaluate(tree, " |---")
}

//解析脚本，并返回根节点
func (c *SimpleCalculator) Parse(code string) astNode.ASTNode {
	lexer := NewSimpleLexer()
	tokens := lexer.Tokenize(code)

	rootNode := c.Prog(tokens)
	return rootNode
}

func (c *SimpleCalculator) evaluate(node astNode.ASTNode, indent string) int {
	result := 0
	fmt.Println(indent+"Calculating:", astNode.GetAstNodeTypeName(node.GetType()))
	switch node.GetType() {
	case astNode.Programm:
		for _, child := range node.GetChildren() {
			result = c.evaluate(child, " | "+indent)
		}
		break
	case astNode.Additive:
		child1 := node.GetChildren()[0]
		value1 := c.evaluate(child1, " | "+indent)
		child2 := node.GetChildren()[1]
		value2 := c.evaluate(child2, " | "+indent)
		if node.GetText() == "+" {
			result = value1 + value2
		} else {
			result = value1 - value2
		}
		break
	case astNode.Multiplicative:
		child1 := node.GetChildren()[0]
		value1 := c.evaluate(child1, " | "+indent)
		child2 := node.GetChildren()[1]
		value2 := c.evaluate(child2, " | "+indent)
		if node.GetText() == "*" {
			result = value1 * value2
		} else {
			result = value1 / value2
		}
		break
	case astNode.IntLiteral:
		result, _ = strconv.Atoi(node.GetText())
		break
	default:
	}
	fmt.Println(indent+"Result:", result)
	return result
}

// 语法解析：根节点
func (c *SimpleCalculator) Prog(tokens tok.TokenReader) *SimpleASTNode {
	node := NewSimpleAstNode(astNode.Programm, "Calculator")

	child := c.Additive(tokens)

	if child != nil {
		node.AddChild(child)
	}
	return node
}

/**
* 整型变量声明语句，如：
* int a;
* int b = 2*3;
 */
func (c *SimpleCalculator) IntDeclare(tokens tok.TokenReader) *SimpleASTNode {
	var node *SimpleASTNode
	token := tokens.Peek()
	if token != nil && token.GetType() == tok.Int { //匹配Int
		tokens.Read()
		if tokens.Peek().GetType() == tok.Identifier { //匹配标识符
			token = tokens.Read()
			node = NewSimpleAstNode(astNode.IntDeclaration, token.GetText())
			token = tokens.Peek()
			if token != nil && token.GetType() == tok.Assignment {
				tokens.Read()
				child := c.Additive(tokens) //匹配一个表达式
				if child == nil {
					panic("invalid variable initialization, expecting an expression")
				} else {
					node.AddChild(child)
				}
			}
		} else {
			panic("variable name expected")
		}

		if node != nil {
			token := tokens.Peek()
			if token != nil && token.GetType() == tok.SemiColon {
				tokens.Read()
			} else {
				panic("invalid statement, expecting semicolon")
			}
		}
	}
	return node
}

//语法解析：加法表达式
func (c *SimpleCalculator) Additive(tokens tok.TokenReader) *SimpleASTNode {
	child1 := c.Multiplicative(tokens)
	node := child1

	token := tokens.Peek()
	if child1 != nil && token != nil {
		if token.GetType() == tok.Plus || token.GetType() == tok.Minus {
			token = tokens.Read()
			child2 := c.Additive(tokens)
			if child2 != nil {
				node = NewSimpleAstNode(astNode.Additive, token.GetText())
				node.AddChild(child1)
				node.AddChild(child2)
			} else {
				panic("invalid additive expression, expecting the right part.")
			}
		}
	}
	return node
}

func (c *SimpleCalculator) Multiplicative(tokens tok.TokenReader) *SimpleASTNode {
	child1 := c.Primary(tokens)
	node := child1

	token := tokens.Peek()
	if child1 != nil && token != nil {
		if token.GetType() == tok.Star || token.GetType() == tok.Slash {
			token = tokens.Read()
			child2 := c.Multiplicative(tokens)
			if child2 != nil {
				node = NewSimpleAstNode(astNode.Multiplicative, token.GetText())
				node.AddChild(child1)
				node.AddChild(child2)
			} else {
				panic("invalid multiplicative expression, expecting the right part.")
			}
		}
	}
	return node
}

func (c *SimpleCalculator) Primary(tokens tok.TokenReader) *SimpleASTNode {
	var node *SimpleASTNode
	token := tokens.Peek()

	if token != nil {
		if token.GetType() == tok.IntLiteral {
			token = tokens.Read()
			node = NewSimpleAstNode(astNode.IntLiteral, token.GetText())
		} else if token.GetType() == tok.Identifier {
			token = tokens.Read()
			node = NewSimpleAstNode(astNode.Identifier, token.GetText())
		} else if token.GetType() == tok.LeftParen {
			tokens.Read()
			node = c.Additive(tokens)
			if node != nil {
				token = tokens.Peek()
				if token != nil && token.GetType() == tok.RightParen {
					tokens.Read()
				} else {
					panic("expecting right parenthesis")
				}
			} else {
				panic("expecting an additive expression inside parenthesis")
			}
		}
	}
	return node
}
