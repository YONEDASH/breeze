package common

type Position struct {
	Index  int
	Line   int
	Column int
}

func InitPosition() Position {
	return Position{Line: 1, Column: 1, Index: 0}
}
