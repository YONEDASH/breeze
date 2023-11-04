package scanner

type TokenId int8

const (
	Invalid = iota
	Eof

	OpenParen
	CloseParen

	// Operators
	Dot

	Identifier

	Integer
	Float
	String
)

type Position struct {
	Index  int
	Line   int
	Column int
}

func initPosition() Position {
	return Position{Line: 1, Column: 1, Index: 0}
}

type Token struct {
	Id       TokenId
	Lexeme   string
	Position Position
}
