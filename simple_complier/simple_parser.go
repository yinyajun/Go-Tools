package front_end

import (
	"fmt"
	"simple_complier/astNode"
	tok "simple_complier/token"
)

/**
 * 一个简单的语法解析器。
 * 能够解析简单的表达式、变量声明和初始化语句、赋值语句。
 * 它支持的语法规则为：
 *
 * programm -> intDeclare | expressionStatement | assignmentStatement
 * intDeclare -> 'int' Id ( = additive) ';'
 * expressionStatement -> addtive ';'
 * additive -> multiplicative ( (+ | -) multiplicative)*
 * multiplicative -> primary ( (* | /) primary)*
 * primary -> IntLiteral | Id | (additive)
 */
type SimpleParser struct{}

func (p *SimpleParser) Parse(script string) astNode.ASTNode {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	lexer := NewSimpleLexer()
	tokens := lexer.Tokenize(script)
	rootNode := p.prog(tokens)
	return rootNode
}

func (p *SimpleParser) prog(tokens tok.TokenReader) *SimpleASTNode {
	node := NewSimpleAstNode(astNode.Programm, "pwc")

	for tokens.Peek() != nil {
		child := p.intDeclare(tokens)

		if child == nil {
			child = p.expressionStatement(tokens)
		}

		if child == nil {
			child = p.assignmentStatement(tokens)
		}

		if child != nil {
			node.AddChild(child)
		} else {
			panic("unknown statement")
		}
	}
	return node

}

func (p *SimpleParser) intDeclare(tokens tok.TokenReader) *SimpleASTNode {
	var node *SimpleASTNode
	token := tokens.Peek()
	if token != nil && token.GetType() == tok.Int {
		tokens.Read()
		if tokens.Peek().GetType() == tok.Identifier {
			token := tokens.Read()
			node = NewSimpleAstNode(astNode.IntDeclaration, token.GetText())
			token = tokens.Peek()
			if token != nil && token.GetType() == tok.Assignment {
				tokens.Read()
				child := p.additive(tokens)
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
			token = tokens.Peek()
			if token != nil && token.GetType() == tok.SemiColon {
				tokens.Read()
			} else {
				panic("invalid statement, expecting semicolon.")
			}
		}
	}
	return node
}

func (p *SimpleParser) expressionStatement(tokens tok.TokenReader) *SimpleASTNode {
	pos := tokens.GetPosition()
	node := p.additive(tokens)

	if node != nil {
		token := tokens.Peek()
		if token != nil && token.GetType() == tok.SemiColon {
			tokens.Read()
		} else {
			node = nil
			tokens.SetPosition(pos) // 回溯
		}
	}
	return node
}

func (p *SimpleParser) assignmentStatement(tokens tok.TokenReader) *SimpleASTNode {
	var node *SimpleASTNode
	token := tokens.Peek()
	if token != nil && token.GetType() == tok.Identifier {
		token = tokens.Read()
		node = NewSimpleAstNode(astNode.AssignmentStmt, token.GetText())
		token = tokens.Peek()
		if token != nil && token.GetType() == tok.Assignment {
			tokens.Read()
			child := p.additive(tokens)
			if child == nil {
				panic("invalid assignment statement, expecting an expression.")
			} else {
				node.AddChild(child)
				token = tokens.Peek()
				if token != nil && token.GetType() == tok.SemiColon {
					tokens.Read()
				} else {
					panic("invalid statement, expecting semicolon.")
				}
			}
		} else {
			tokens.Unread()
			node = nil
		}
	}
	return node
}

func (p *SimpleParser) additive(tokens tok.TokenReader) *SimpleASTNode {
	child1 := p.multiplicative(tokens)
	node := child1

	if child1 != nil {
		for {
			token := tokens.Peek()
			if token != nil && (token.GetType() == tok.Plus || token.GetType() == tok.Minus) {
				token = tokens.Read()
				child2 := p.multiplicative(tokens)
				if child2 != nil {
					node = NewSimpleAstNode(astNode.Additive, token.GetText())
					node.AddChild(child1)
					node.AddChild(child2)
				} else {
					panic("invalid additive expression, expecting the right part.")
				}
			} else {
				break
			}
		}
	}
	return node
}

func (p *SimpleParser) multiplicative(tokens tok.TokenReader) *SimpleASTNode {
	child1 := p.primary(tokens)
	node := child1

	//if child1 == nil {
	//	return node
	//}
	for {
		token := tokens.Peek()
		if token != nil && (token.GetType() == tok.Star || token.GetType() == tok.Slash) {
			token = tokens.Read()
			child2 := p.primary(tokens)
			if child2 != nil {
				node = NewSimpleAstNode(astNode.Multiplicative, token.GetText())
				node.AddChild(child1)
				node.AddChild(child2)
				child1 = node
			} else {
				panic("invalid multiplicative expression, expecting the right part.")
			}
		} else {
			break
		}
	}
	return node
}

func (p *SimpleParser) primary(tokens tok.TokenReader) *SimpleASTNode {
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
			node = p.additive(tokens)
			if node != nil {
				token = tokens.Peek()
				if token != nil && token.GetType() == tok.RightParen {
					tokens.Read()
				} else {
					panic("expecting right parenthesis.")
				}
			} else {
				panic("expecting an additive expression inside parenthesis.")
			}
		}
	}
	return node
}
