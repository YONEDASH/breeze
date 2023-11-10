package scanner

import (
	"breeze/common"
	"fmt"
)

type TokenId uint8

//goland:noinspection GoCommentStart
const (
	Invalid = iota
	EOF

	// Operators
	Plus
	Minus
	Star
	Slash
	Equals
	Bang

	// 2-char operators
	PlusEquals
	MinusEquals
	StarEquals
	SlashEquals

	// Comparative
	Lower
	Greater
	LowerEquals
	GreaterEquals
	EqualsEquals
	BangEquals

	// Keywords
	Debug // prints node value
	Let
	If
	Else

	// Literals
	Identifier
	Integer
	Float
	String
	True
	False

	// Others
	OpenParen
	CloseParen
	OpenBrace
	CloseBrace
	OpenBracket
	CloseBracket
	Semicolon
	Colon
)

type Token struct {
	Id       TokenId
	Lexeme   string
	Position common.Position
}

func (t *Token) Stringify() string {
	return fmt.Sprintf("#%2d: %s", t.Id, t.Lexeme)
}

func (t *Token) LexemeLength() int {
	runes := []rune(t.Lexeme)
	return len(runes)
}
