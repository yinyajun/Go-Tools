package astNode

type ASTNodeType uint8

const (
	Programm ASTNodeType = iota //程序入口，根节点

	IntDeclaration //整型变量声明
	ExpressionStmt //表达式语句，即表达式后面跟个分号
	AssignmentStmt //赋值语句

	Primary        //基础表达式
	Multiplicative //乘法表达式
	Additive       //加法表达式

	Identifier //标识符
	IntLiteral //整型字面量
)

type ASTNode interface {
	GetParent() ASTNode
	GetChildren() []ASTNode
	GetType() ASTNodeType
	GetText() string
}

func GetAstNodeTypeName(t ASTNodeType) string {
	dict := map[ASTNodeType]string{
		Programm:       "Programm",
		IntDeclaration: "IntDeclaration",
		ExpressionStmt: "ExpressionStmt",
		AssignmentStmt: "AssignmentStmt",
		Primary:        "Primary",
		Multiplicative: "Multiplicative",
		Additive:       "Additive",
		Identifier:     "Identifier",
		IntLiteral:     "IntLiteral",
	}
	return dict[t]
}
