package token

type TokenType uint8

// enum
const (
	PlaceHolder TokenType = iota

	Plus
	Minus // -
	Star  // *
	Slash // /

	GE // >=
	GT // >
	EQ // ==
	LE // <=
	LT // <

	SemiColon  // ;
	LeftParen  // (
	RightParen // )

	Assignment // =

	If
	Else

	Int

	Identifier //标识符

	IntLiteral    //整型字面量
	StringLiteral //字符串字面量
)

func GetTokenTypeName(t TokenType) string {
	dict := map[TokenType]string{
		PlaceHolder:   "unk",
		Plus:          "Plus",
		Minus:         "Minus",
		Star:          "Star",
		Slash:         "Slash",
		GE:            "GE",
		GT:            "GT",
		EQ:            "EQ",
		LE:            "LE",
		LT:            "LT",
		SemiColon:     "SemiColon",
		LeftParen:     "LeftParen",
		RightParen:    "RightParen",
		Assignment:    "Assignment",
		If:            "If",
		Else:          "Else",
		Int:           "Int",
		Identifier:    "Identifier",
		IntLiteral:    "IntLiteral",
		StringLiteral: "StringLiteral",
	}
	return dict[t]
}

/**
* 一个简单的Token。
* 只有类型和文本值两个属性。
 */
type Token interface {
	GetType() TokenType //Token的类型
	GetText() string    // Token的文本值
	SetType(t TokenType)
	SetText(c string)
}

/**
* 一个Token流。由Lexer生成。Parser可以从中获取Token。
 */
type TokenReader interface {
	Read() Token              // 返回Token流中下一个Token，并从流中取出。 如果流已经为空，返回nil
	Peek() Token              // 返回Token流中下一个Token，但不从流中取出。 如果流已经为空，返回nil
	Unread()                  // Token流回退一步。恢复原来的Token
	GetPosition() int         // 获取Token流当前的读取位置
	SetPosition(position int) // 设置Token流当前的读取位置
}
