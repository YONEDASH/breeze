package scanner

import (
	"breeze/common"
	"breeze/out"
	"os"
)

type sourceScanner struct {
	source []rune
	length int
	start  common.Position
	cursor common.Position
	file   *common.SourceFile
}

func (s *sourceScanner) isDone() bool {
	return s.cursor.Index >= s.length
}

func (s *sourceScanner) advance() rune {
	if s.isDone() {
		return 0
	}
	s.cursor.Index++
	s.cursor.Column++
	return s.source[s.cursor.Index-1]
}

func (s *sourceScanner) peekPrevious() rune {
	if s.cursor.Index == 0 {
		return 0
	}
	return s.source[s.cursor.Index-1]
}

func (s *sourceScanner) peek() rune {
	if s.isDone() {
		return 0
	}
	return s.source[s.cursor.Index]
}

func (s *sourceScanner) peekNext() rune {
	if s.cursor.Index >= s.length-1 {
		return 0
	}
	return s.source[s.cursor.Index+1]
}

func (s *sourceScanner) match(r rune) bool {
	if !s.isDone() && s.source[s.cursor.Index] == r {
		s.advance()
		return true
	}
	return false
}

func initScanner(file *common.SourceFile, source string) sourceScanner {
	runes := []rune(source)
	runesLen := len(runes)
	return sourceScanner{
		source: runes,
		length: runesLen,
		start:  common.InitPosition(),
		cursor: common.InitPosition(),
		file:   file,
	}
}

func makeToken(scanner *sourceScanner, id TokenId) Token {
	lexeme := string(scanner.source[scanner.start.Index:scanner.cursor.Index])
	position := scanner.start
	scanner.start = scanner.cursor
	return Token{
		Id:       id,
		Lexeme:   lexeme,
		Position: position,
	}
}

func errorToken(scanner *sourceScanner, message string) Token {
	position := scanner.start
	scanner.start = scanner.cursor
	return Token{
		Id:       Invalid,
		Lexeme:   message,
		Position: position,
	}
}

func skipWhitespace(scanner *sourceScanner) {
	for {
		if scanner.isDone() {
			return
		}

		switch scanner.peek() {
		case '\n':
			scanner.cursor.Line++
			scanner.cursor.Column = 0
			scanner.advance()
			scanner.start = scanner.cursor
			break
		case ' ', '\t':
			scanner.advance()
			scanner.start = scanner.cursor
			break
		default:
			return
		}
	}
}

func isAlpha(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_'
}

func isNumber(r rune) bool {
	return r >= '0' && r <= '9'
}

func identifier(scanner *sourceScanner) Token {
	for {
		if scanner.isDone() {
			break
		}

		current := scanner.peek()
		if !isAlpha(current) && !isNumber(current) {
			break
		}

		scanner.advance()
	}

	// Keywords
	lexeme := string(scanner.source[scanner.start.Index:scanner.cursor.Index])
	switch lexeme {
	case "let":
		return makeToken(scanner, Let)
	}

	return makeToken(scanner, Identifier)
}

func number(scanner *sourceScanner) Token {
	hadDot := scanner.peekPrevious() == '.'

	for {
		if scanner.isDone() {
			break
		}

		current := scanner.peek()

		if !hadDot && current == '.' {
			hadDot = true
			scanner.advance()
			continue
		}

		if !isNumber(current) {
			break
		}

		scanner.advance()
	}

	if hadDot {
		return makeToken(scanner, Float)
	}

	return makeToken(scanner, Integer)
}

func text(scanner *sourceScanner) Token {
	for {
		if scanner.isDone() {
			return errorToken(scanner, "Expected closing \"")
		}

		current := scanner.peek()
		if current == '"' && scanner.peekPrevious() != '\\' {
			scanner.advance()
			break
		}

		scanner.advance()
	}

	// ignore opening "
	scanner.start.Index++

	// ignore closing "
	scanner.cursor.Index--

	token := makeToken(scanner, String)

	// revert our ignoring magic
	scanner.cursor.Index++
	scanner.start = scanner.cursor

	return token
}

func scanToken(scanner *sourceScanner) Token {
	skipWhitespace(scanner)

	current := scanner.advance()

	// Identifier
	if isAlpha(current) {
		return identifier(scanner)
	}

	// Number
	if isNumber(current) || current == '.' {
		return number(scanner)
	}

	switch current {
	case '"':
		return text(scanner)
	case '=':
		return makeToken(scanner, Equals)
	case '+':
		if scanner.peek() == '=' {
			scanner.advance()
			return makeToken(scanner, PlusEquals)
		}
		return makeToken(scanner, Plus)
	case '-':
		if scanner.peek() == '=' {
			scanner.advance()
			return makeToken(scanner, MinusEquals)
		}
		return makeToken(scanner, Minus)
	case '*':
		if scanner.peek() == '=' {
			scanner.advance()
			return makeToken(scanner, StarEquals)
		}
		return makeToken(scanner, Star)
	case '/':
		if scanner.peek() == '=' {
			scanner.advance()
			return makeToken(scanner, SlashEquals)
		}
		return makeToken(scanner, Slash)
	case ';':
		return makeToken(scanner, Semicolon)
	case ':':
		return makeToken(scanner, Colon)
	case '(':
		return makeToken(scanner, OpenParen)
	case ')':
		return makeToken(scanner, CloseParen)
	}

	return errorToken(scanner, "Unexpected token")
}

func Scan(file *common.SourceFile, source string) ([]Token, bool) {
	scanner := initScanner(file, source)
	var tokens []Token
	var hadError = false

	for {
		if scanner.isDone() {
			break
		}

		token := scanToken(&scanner)
		if token.Id == Invalid {
			hadError = true

			if scanner.peekPrevious() != '💨' {
				out.PrintErrorMessage(token.Lexeme)
			} else {
				out.PrintErrorMessage("This breeze is unfortunately an unexpected token")
			}

			out.PrintErrorSource(file.Path, token.Position)
			out.PrintMarkedLine(os.Stderr, source, 1, token.Position, out.ColorRed, '^')

			continue
		}
		tokens = append(tokens, token)
	}

	tokens = append(tokens, makeToken(&scanner, EOF))

	return tokens, hadError
}
