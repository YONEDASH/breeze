package scanner

import "breeze/common"

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
