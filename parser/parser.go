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

		node := declaration(&parser)

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

func expectSemicolon(parser *tokenParser, result ast.Node) ast.Node {
	// Check for ;
	if parser.advance().Id != scanner.Semicolon {
		return err(parser.peekPrevious(), "Unfinished statement", "Add ; to end of statement")
	}
	return result
}

func declaration(parser *tokenParser) ast.Node {
	current := parser.peek()

	switch current.Id {
	case scanner.Let:
		return let(parser)
	}

	return statement(parser)
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
		_ = parser.advance()

		if parser.peek().Id != scanner.Identifier {
			return err(identifierToken, "Expected identifier in let type declaration", "")
		}

		lexeme := parser.advance().Lexeme
		varType = lexeme
	}

	letDecl := &ast.LetDecl{Token: keyword, Identifier: varName, Type: varType}

	// Check for assign
	if parser.peek().Id != scanner.Equals {
		return expectSemicolon(parser, letDecl)
	}

	// Assign desugaring

	// Consume =
	operator := parser.advance()

	// Expression for value
	expr := expression(parser)
	if expr.GetId() == ast.ErrId {
		return expr
	}

	assignExpr := &ast.AssignExpr{Operator: operator, Name: identifierToken, Value: expr}

	blockStmt := &ast.BlockStmt{Token: keyword, Nodes: []ast.Node{letDecl, assignExpr}}
	return expectSemicolon(parser, blockStmt)
}

func statement(parser *tokenParser) ast.Node {
	current := parser.peek()

	switch current.Id {
	case scanner.If:
		return conditional(parser)
	case scanner.OpenBrace:
		return closure(parser)
	case scanner.Debug:
		return debug(parser)
	}

	// Parse expression statement
	expr := expression(parser)
	if expr.GetId() == ast.ErrId {
		return expr
	}

	result := &ast.ExprStmt{Token: expr.GetToken(), Expression: expr}
	return expectSemicolon(parser, result)
}

func conditional(parser *tokenParser) ast.Node {
	// Consume if
	keyword := parser.advance()

	condition := expression(parser)
	if condition.GetId() == ast.ErrId {
		return condition
	}

	statement := declaration(parser)
	if statement.GetId() == ast.ErrId {
		return statement
	}

	var elseStatement ast.Node
	if parser.peek().Id == scanner.Else {
		_ = parser.advance()
		elseStatement = declaration(parser)
		if elseStatement.GetId() == ast.ErrId {
			return elseStatement
		}
	}

	return &ast.ConditionalStmt{Token: keyword, Statement: statement, ElseStatement: elseStatement, Condition: condition}
}

func closure(parser *tokenParser) ast.Node {
	// Consume {
	keyword := parser.advance()

	nodes := make([]ast.Node, 0)

	for {
		if parser.isDone() {
			break
		}

		if parser.peek().Id == scanner.CloseBrace {
			parser.advance()
			break
		}

		node := declaration(parser)
		nodes = append(nodes, node)
	}

	block := &ast.BlockStmt{Token: keyword, Nodes: nodes}
	return &ast.ClosureStmt{Token: keyword, Block: block}
}

func debug(parser *tokenParser) ast.Node {
	keyword := parser.advance()
	expr := expression(parser)
	if expr.GetId() == ast.ErrId {
		return expr
	}
	return expectSemicolon(parser, &ast.DebugStmt{Token: keyword, Expression: expr})
}

func expression(parser *tokenParser) ast.Node {
	return assign(parser)
}

func assign(parser *tokenParser) ast.Node {
	expr := equality(parser)
	if expr.GetId() == ast.ErrId {
		return expr
	}

	var current = parser.peek()

	if current.Id != scanner.Equals && current.Id != scanner.PlusEquals && current.Id != scanner.MinusEquals && current.Id != scanner.StarEquals && current.Id != scanner.SlashEquals {
		return expr
	}

	// consume operator
	operator := parser.advance()

	right := assign(parser)
	if right.GetId() == ast.ErrId {
		return right
	}

	if expr.GetId() == ast.IdentifierLitId {
		identifier := expr.(*ast.IdentifierLitExpr)
		return &ast.AssignExpr{
			Operator: operator,
			Name:     identifier.GetToken(),
			Value:    right,
		}
	}

	return err(operator, "Unsupported assign operation on token", "Expected identifier TODO: or call")
}

func equality(parser *tokenParser) ast.Node {
	left := comparison(parser)

	for {
		if parser.isDone() {
			break
		}

		current := parser.peek()
		if current.Id != scanner.EqualsEquals && current.Id != scanner.BangEquals {
			break
		}

		operator := parser.advance()
		right := comparison(parser)

		left = &ast.BinaryExpr{Operator: operator, Left: left, Right: right}
	}

	return left
}

func comparison(parser *tokenParser) ast.Node {
	left := add(parser)

	for {
		if parser.isDone() {
			break
		}

		current := parser.peek()
		if current.Id != scanner.Lower && current.Id != scanner.Greater && current.Id != scanner.LowerEquals && current.Id != scanner.GreaterEquals {
			break
		}

		operator := parser.advance()
		right := add(parser)

		left = &ast.BinaryExpr{Operator: operator, Left: left, Right: right}
	}

	return left
}

func add(parser *tokenParser) ast.Node {
	left := multiply(parser)
	if left.GetId() == ast.ErrId {
		return left
	}

	for {
		if parser.isDone() {
			break
		}

		operator := parser.peek()

		if operator.Id != scanner.Plus && operator.Id != scanner.Minus {
			break
		}

		// consume operator
		_ = parser.advance()

		right := multiply(parser)
		if right.GetId() == ast.ErrId {
			return right
		}

		left = &ast.BinaryExpr{Operator: operator, Left: left, Right: right}
	}

	return left
}

func multiply(parser *tokenParser) ast.Node {
	left := unary(parser)
	if left.GetId() == ast.ErrId {
		return left
	}

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

		right := unary(parser)
		if right.GetId() == ast.ErrId {
			return right
		}

		left = &ast.BinaryExpr{Operator: operator, Left: left, Right: right}
	}

	return left
}

func unary(parser *tokenParser) ast.Node {
	current := parser.peek()

	switch current.Id {
	case scanner.Plus, scanner.Minus, scanner.Bang:
		_ = parser.advance()

		// No recursive parsing for unary. We don't want something like +-+--10 or !!!!!boolVal
		expr := primary(parser)

		if expr.GetId() == ast.ErrId {
			return expr
		}

		return &ast.UnaryExpr{Operator: current, Expression: expr}
	}

	return primary(parser)
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
		return &ast.IdentifierLitExpr{Token: current, Name: current.Lexeme}

	case scanner.Integer:
		return &ast.IntegerLitExpr{Token: current, Value: current.Lexeme}

	case scanner.Float:
		return &ast.FloatingLitExpr{Token: current, Value: current.Lexeme}

	case scanner.True, scanner.False:
		return &ast.BooleanLitExpr{Token: current, Value: current.Lexeme}

	}

	return err(current, "Unexpected token", "")
}

func err(token scanner.Token, message string, hint string) ast.Node {
	return &ast.ErrNode{
		Token:   token,
		Message: message,
		Hint:    hint,
	}
}
