package scanner

import (
	"breeze/common"
	"fmt"
)

type TokenId int8

//goland:noinspection GoCommentStart
const (
	Invalid = iota
	Eof

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
