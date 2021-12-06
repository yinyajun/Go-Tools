package front_end

import . "simple_complier/token"

/**
* Token的一个简单实现。只有类型和文本值两个属性。
 */
type SimpleToken struct {
	typ  TokenType
	text string
}

func (s *SimpleToken) GetType() TokenType { return s.typ }

func (s *SimpleToken) GetText() string { return s.text }

func (s *SimpleToken) SetType(t TokenType) { s.typ = t }

func (s *SimpleToken) SetText(c string) { s.text = c }

/**
* 一个简单的Token流。是把一个Token列表进行了封装。
 */
type SimpleTokenReader struct {
	Tokens []Token
	Pos    int
}

func NewSimpleTokenReader(tokens []Token) *SimpleTokenReader {
	return &SimpleTokenReader{Tokens: tokens}
}

func (r *SimpleTokenReader) Read() Token {
	if r.Pos < len(r.Tokens) {
		ret := r.Tokens[r.Pos]
		r.Pos++
		return ret
	}
	return nil
}

func (r *SimpleTokenReader) Peek() Token {
	if r.Pos < len(r.Tokens) {
		ret := r.Tokens[r.Pos]
		return ret
	}
	return nil
}

func (r *SimpleTokenReader) Unread() {
	if r.Pos > 0 {
		r.Pos--
	}
}

func (r *SimpleTokenReader) GetPosition() int { return r.Pos }

func (r *SimpleTokenReader) SetPosition(position int) {
	if position > 0 && position < len(r.Tokens) {
		r.Pos = position
	}
}
