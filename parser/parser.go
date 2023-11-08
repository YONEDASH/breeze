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
	return p.cursor >= p.length || p.peek().Id == scanner.EOF
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

func synchronize(parser *tokenParser, id scanner.TokenId) {
	parser.cursor--

	// Throw away any token until the one to synchronize to is found
	for {
		if parser.isDone() {
			break
		}

		if parser.advance().Id == id {
			break
		}
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

		node := statement(&parser)

		if node.GetType() == ast.Err {
			hadError = true
			errNode := node.(*ast.ErrNode)
			token := node.GetToken()
			pos := node.GetToken().Position

			out.PrintErrorMessage(errNode.Message)
			out.PrintErrorSource(file.Path, pos)
			out.PrintMarkedLine(os.Stderr, source, token.LexemeLength(), pos, out.ColorRed, '^')

			hintLen := len(errNode.Hint)
			if hintLen > 0 {
				out.PrintHintMessage(errNode.Hint, out.ColorRed)
			}

			// Synchronize to ;
			synchronize(&parser, scanner.Semicolon)

			continue
		}

		nodes = append(nodes, node)
	}

	return nodes, hadError
}

func declaration(parser *tokenParser) ast.Node {
	current := parser.advance()

	return err(current, "Unexpected token", "Expecting let, fn")
}

func statement(parser *tokenParser) ast.Node {
	var result ast.Node

	switch parser.peek().Id {
	case scanner.Let:
		result = let(parser)
		break
	}

	if result == nil {
		// Parse expression statement
		expr := expression(parser)
		result = &ast.ExprStmt{Token: expr.GetToken(), Expression: expr}
	}

	// Check for ;
	if parser.advance().Id != scanner.Semicolon {
		return err(parser.peekPrevious(), "Unfinished statement", "Add ; to end of statement")
	}

	return result
}

func let(parser *tokenParser) ast.Node {
	keyword := parser.advance()

	identifierToken := parser.advance()
	if identifierToken.Id != scanner.Identifier {
		return err(identifierToken, "Expected identifier in let name declaration", "")
	}
	varName := identifierToken.Lexeme

	var varType = ""
	if parser.peek().Id == scanner.Colon {
		// Consume :
		parser.advance()

		if parser.peek().Id != scanner.Identifier {
			return err(identifierToken, "Expected identifier in let type declaration", "")
		}

		lexeme := parser.advance().Lexeme
		varType = lexeme
	}

	letDecl := &ast.LetDecl{Token: keyword, Identifier: varName, Type: varType}

	// Check for assign
	if parser.peek().Id != scanner.Equals {
		return letDecl
	}

	// Assign desugaring

	// Consume =
	operator := parser.advance()

	// Expression for value
	expr := expression(parser)
	identifierExpr := &ast.IdentifierExpr{Token: identifierToken, Name: varName}
	assignExpr := &ast.BinaryExpr{Operator: operator, Left: identifierExpr, Right: expr}

	blockStmt := &ast.BlockStmt{Token: keyword, Nodes: []ast.Node{letDecl, assignExpr}}
	return blockStmt
}

func expression(parser *tokenParser) ast.Node {
	return assign(parser)
}

func assign(parser *tokenParser) ast.Node {
	left := add(parser)

	var current scanner.Token

	for {
		if parser.isDone() {
			break
		}

		current = parser.peek()

		if current.Id != scanner.Equals && current.Id != scanner.PlusEquals && current.Id != scanner.MinusEquals && current.Id != scanner.StarEquals && current.Id != scanner.SlashEquals {
			break
		}

		// consume =
		operator := parser.advance()

		right := expression(parser)
		left = &ast.BinaryExpr{
			Operator: operator,
			Left:     left,
			Right:    right,
		}
	}

	return left
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
	case scanner.OpenParen:
		node := expression(parser)

		if parser.advance().Id != scanner.CloseParen {
			return err(current, "Unclosed grouping expression", "Add missing ) to close group")
		}

		return node

	case scanner.Identifier:
		return &ast.IdentifierExpr{Token: current, Name: current.Lexeme}

	case scanner.Integer:
		return &ast.IntegerExpr{Token: current, Value: current.Lexeme}
	}

	return err(current, "Expected expression", "")
}

func err(token scanner.Token, message string, hint string) ast.Node {
	return &ast.ErrNode{
		Token:   token,
		Message: message,
		Hint:    hint,
	}
}
