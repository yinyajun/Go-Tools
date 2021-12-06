package front_end

import (
	"fmt"
	. "simple_complier/astNode"
)

type SimpleASTNode struct {
	parent   *SimpleASTNode
	children []*SimpleASTNode
	nodeType ASTNodeType
	text     string
}

func NewSimpleAstNode(nodeType ASTNodeType, text string) *SimpleASTNode {
	n := &SimpleASTNode{}
	n.children = []*SimpleASTNode{}
	n.nodeType = nodeType
	n.text = text
	return n
}

func (c *SimpleASTNode) GetParent() ASTNode { return c.parent }

func (c *SimpleASTNode) GetChildren() []ASTNode {
	var ret []ASTNode
	for _, c := range c.children {
		ret = append(ret, ASTNode(c))
	}

	return ret
}

func (c *SimpleASTNode) GetType() ASTNodeType { return c.nodeType }

func (c *SimpleASTNode) GetText() string { return c.text }

func (c *SimpleASTNode) AddChild(child *SimpleASTNode) {
	c.children = append(c.children, child)
	child.parent = c
}

func DumpAST(node ASTNode, indent string) {
	fmt.Println(indent + GetAstNodeTypeName(node.GetType()) + " " + node.GetText())
	for _, child := range node.GetChildren() {
		DumpAST(child, " | "+indent)
	}
}
