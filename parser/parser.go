package parser

import (
	"breeze/ast"
	"breeze/common"
	"breeze/out"
	"breeze/scanner"
	"os"
)

type tokenParser struct {
	tokens []scanner.Token
	length int
	cursor int
}

var emptyToken = scanner.Token{Id: scanner.EOF, Lexeme: "empty token", Position: common.InitPosition()}

func (p *tokenParser) isDone() bool {
	return p.peek().Id == scanner.EOF || p.cursor >= p.length
}

func (p *tokenParser) advance() scanner.Token {
	p.cursor++
	return p.tokens[p.cursor-1]
}

func (p *tokenParser) peek() scanner.Token {
	return p.tokens[p.cursor]
}

func (p *tokenParser) peekNext() scanner.Token {
	if p.cursor >= p.length-1 {
		return emptyToken
	}
	return p.tokens[p.cursor+1]
}

func (p *tokenParser) peekPrevious() scanner.Token {
	if p.cursor <= 0 {
		return emptyToken
	}
	return p.tokens[p.cursor-1]
}

func (p *tokenParser) expect(id scanner.TokenId) bool {
	return p.peek().Id == id
}

func initParser(tokens []scanner.Token) tokenParser {
	return tokenParser{
		tokens: tokens,
		length: len(tokens),
		cursor: 0,
	}
}

func ParseTokens(file common.SourceFile, source string, tokens []scanner.Token) ([]ast.Node, bool) {
	parser := initParser(tokens)
	hadError := false
	var nodes []ast.Node

	for {
		if parser.isDone() {
			break
		}

		node := expression(&parser)

		if node.GetType() == ast.Err {
			hadError = true
			err := node.(*ast.ErrNode)
			token := node.GetToken()
			pos := node.GetToken().Position

			out.PrintErrorMessage(err.Message)
			out.PrintErrorSource(file.Path, pos)
			out.PrintMarkedLine(os.Stderr, source, len(token.Lexeme), pos, out.ColorRed, '^')
			// Synchronize to ;
			continue
		}

		nodes = append(nodes, node)
	}

	return nodes, hadError
}

func declaration(parser *tokenParser) ast.Node {
	return err(emptyToken, "Unexpected token")
}

func statement(parser *tokenParser) ast.Node {
	return err(emptyToken, "Unexpected token")
}

func expression(parser *tokenParser) ast.Node {
	return add(parser)
}

func add(parser *tokenParser) ast.Node {
	left := multiply(parser)

	for {
		if parser.isDone() {
			break
		}

		operator := parser.peek()

		if operator.Id != scanner.Plus && operator.Id != scanner.Minus {
			break
		}

		// consume operator
		parser.advance()

		right := multiply(parser)

		left = &ast.BinaryExpr{Operator: operator, Left: left, Right: right}
	}

	return left
}

func multiply(parser *tokenParser) ast.Node {
	left := primary(parser)

	for {
		if parser.isDone() {
			break
		}

		operator := parser.peek()

		if operator.Id != scanner.Star && operator.Id != scanner.Slash {
			break
		}

		// consume operator
		parser.advance()

		right := primary(parser)

		left = &ast.BinaryExpr{Operator: operator, Left: left, Right: right}
	}

	return left
}

func primary(parser *tokenParser) ast.Node {
	current := parser.advance()

	switch current.Id {
	case scanner.Integer:
		return &ast.IntegerExpr{Token: current, Value: current.Lexeme}
	}

	return err(current, "Expected expression")
}

func err(token scanner.Token, message string) ast.Node {
	return &ast.ErrNode{
		Token:   token,
		Message: message,
	}
}
