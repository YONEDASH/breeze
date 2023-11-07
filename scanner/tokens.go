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

	// Keywords
	Let

	// Literals
	Identifier
	Integer
	Float
	String

	// Others
	OpenParen
	CloseParen
)

type Token struct {
	Id       TokenId
	Lexeme   string
	Position common.Position
}

func (t *Token) Stringify() string {
	return fmt.Sprintf("#%2d: %s", t.Id, t.Lexeme)
}
