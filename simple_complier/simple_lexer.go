package front_end

import (
	st "simple_complier/state"
	. "simple_complier/token"
	"strings"

	"bytes"
	"fmt"
)

type SimpleLexer struct {
	tokenText strings.Builder
	tokens    []Token
	token     Token
}

func NewSimpleLexer() *SimpleLexer { return &SimpleLexer{} }

func (m *SimpleLexer) isAlpha(ch byte) bool { return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' }

func (m *SimpleLexer) isDigit(ch byte) bool { return ch >= '0' && ch <= '9' }

func (m *SimpleLexer) isBlank(ch byte) bool { return ch == ' ' || ch == '\t' || ch == '\n' }

/**
 * 有限状态机进入初始状态。
 * 这个初始状态其实并不做停留，它马上进入其他状态。
 * 开始解析的时候，进入初始状态；某个Token解析完毕，也进入初始状态，在这里把Token记下来，然后建立一个新的Token。
 */
func (m *SimpleLexer) initToken(ch byte) st.DfaState {
	if m.tokenText.Len() > 0 {
		m.token.SetText(m.tokenText.String())
		m.tokens = append(m.tokens, m.token)

		m.tokenText = strings.Builder{}
		m.token = &SimpleToken{}
	}

	var newState st.DfaState

	if m.isAlpha(ch) {
		if ch == 'i' {
			newState = st.Id_int1
		} else {
			newState = st.Id
		}
		m.token.SetType(Identifier)
		m.tokenText.WriteByte(ch)
	} else if m.isDigit(ch) {
		newState = st.IntLiteral
		m.token.SetType(IntLiteral)
		m.tokenText.WriteByte(ch)
	} else if ch == '>' {
		newState = st.GT
		m.token.SetType(GT)
		m.tokenText.WriteByte(ch)
	} else if ch == '+' {
		newState = st.Plus
		m.token.SetType(Plus)
		m.tokenText.WriteByte(ch)
	} else if ch == '-' {
		newState = st.Minus
		m.token.SetType(Minus)
		m.tokenText.WriteByte(ch)
	} else if ch == '*' {
		newState = st.Star
		m.token.SetType(Star)
		m.tokenText.WriteByte(ch)
	} else if ch == '/' {
		newState = st.Slash
		m.token.SetType(Slash)
		m.tokenText.WriteByte(ch)
	} else if ch == ';' {
		newState = st.SemiColon
		m.token.SetType(SemiColon)
		m.tokenText.WriteByte(ch)
	} else if ch == '(' {
		newState = st.LeftParen
		m.token.SetType(LeftParen)
		m.tokenText.WriteByte(ch)
	} else if ch == ')' {
		newState = st.RightParen
		m.token.SetType(RightParen)
		m.tokenText.WriteByte(ch)
	} else if ch == '=' {
		newState = st.Assignment
		m.token.SetType(Assignment)
		m.tokenText.WriteByte(ch)
	} else {
		newState = st.Initial // skip all unknown patterns
	}
	return newState
}

/**
* 解析字符串，形成Token。
* 这是一个有限状态自动机，在不同的状态中迁移。
 */
func (m *SimpleLexer) Tokenize(code string) *SimpleTokenReader {
	var tokSlice []Token
	m.tokens = tokSlice
	reader := bytes.NewReader([]byte(code))

	m.tokenText = strings.Builder{}
	m.token = &SimpleToken{}
	ch := byte(0)

	state := st.Initial

	for {
		ch, err := reader.ReadByte()
		if err != nil {
			break
		}

		switch state {
		case st.Initial:
			state = m.initToken(ch)
			break
		case st.Id:
			if m.isAlpha(ch) || m.isDigit(ch) {
				m.tokenText.WriteByte(ch)
			} else {
				state = m.initToken(ch)
			}
			break
		case st.GT:
			if ch == '=' {
				m.token.SetType(GE)
				state = st.GE
				m.tokenText.WriteByte(ch)
			} else {
				state = m.initToken(ch)
			}
			break
		case st.GE:
			state = m.initToken(ch)
			break
		case st.Assignment:
			state = m.initToken(ch)
			break
		case st.Plus:
			state = m.initToken(ch)
			break
		case st.Minus:
			state = m.initToken(ch)
			break
		case st.Star:
			state = m.initToken(ch)
			break
		case st.Slash:
			state = m.initToken(ch)
			break
		case st.SemiColon:
			state = m.initToken(ch)
			break
		case st.LeftParen:
			state = m.initToken(ch)
			break
		case st.RightParen:
			state = m.initToken(ch)
			break
		case st.IntLiteral:
			if m.isDigit(ch) {
				m.tokenText.WriteByte(ch)
			} else {
				state = m.initToken(ch)
			}
			break
		case st.Id_int1:
			if ch == 'n' {
				state = st.Id_int2
				m.tokenText.WriteByte(ch)
			} else if m.isDigit(ch) || m.isAlpha(ch) {
				state = st.Id
				m.tokenText.WriteByte(ch)
			} else {
				state = m.initToken(ch)
			}
			break
		case st.Id_int2:
			if ch == 't' {
				state = st.Id_int3
				m.tokenText.WriteByte(ch)
			} else if m.isAlpha(ch) || m.isDigit(ch) {
				state = st.Id
				m.tokenText.WriteByte(ch)
			} else {
				state = m.initToken(ch)
			}
			break
		case st.Id_int3:
			if m.isBlank(ch) {
				m.token.SetType(Int)
				state = m.initToken(ch)
			} else {
				state = st.Id
				m.tokenText.WriteByte(ch)
			}
			break
		default:
		}
		//fmt.Println("ch:", string(ch), "state:", state, "token:", m.token, "tokenText:", m.tokenText.String(), "tokens:", m.tokens)
	}
	if m.tokenText.Len() > 0 {
		m.initToken(ch)
	}
	return NewSimpleTokenReader(m.tokens)
}

func (m *SimpleLexer) Dump(reader TokenReader) {
	fmt.Println("text\ttype")
	for {
		token := reader.Read()
		if token == nil {
			break
		}
		fmt.Printf("%s\t\t%s\n", token.GetText(), GetTokenTypeName(token.GetType()))
	}
}
