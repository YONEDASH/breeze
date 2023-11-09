package analyzer

import (
	"breeze/ast"
	"breeze/common"
	"breeze/out"
	"fmt"
	"os"
)

type DeclarationType uint8

const (
	VariableDeclaration DeclarationType = iota
)

type declaration struct {
	Type        DeclarationType
	DeclaredAt  ast.Node
	Initialized bool
	StaticType  string
}

type Scope struct {
	Declared map[string]declaration
}

func initScope() Scope {
	scope := Scope{
		Declared: make(map[string]declaration),
	}
	return scope
}

type FunctionType uint8

const (
	None FunctionType = iota
)

type Context struct {
	ast.Visitor
	File            common.SourceFile
	Source          string
	HadError        bool
	Stack           []Scope
	CurrentFunction FunctionType
}

func Analyze(sourceFile common.SourceFile, source string, nodes []ast.Node) bool {
	context := &Context{Stack: make([]Scope, 0), HadError: false, CurrentFunction: None, Source: source, File: sourceFile}
	context.begin()
	for _, node := range nodes {
		_ = node.Visit(context)
	}
	context.end()
	return context.HadError
}

func (c *Context) push(scope Scope) {
	c.Stack = append(c.Stack, scope)
}

func (c *Context) top() Scope {
	size := len(c.Stack)
	return c.Stack[size-1]
}

func (c *Context) pop() Scope {
	size := len(c.Stack)
	top := c.top()
	c.Stack = c.Stack[:size-1]
	return top
}

func (c *Context) empty() bool {
	return len(c.Stack) == 0
}

func (c *Context) nodeError(node ast.Node, message string) {
	c.HadError = true
	out.PrintErrorMessage(message)
	token := node.GetToken()
	out.PrintErrorSource(c.File.Path, token.Position)
	out.PrintMarkedLine(os.Stderr, c.Source, token.LexemeLength(), token.Position, out.ColorRed, '^')
}

func (c *Context) nodeHelpHint(node ast.Node, message string) {
	token := node.GetToken()
	out.PrintErrorSource(c.File.Path, token.Position)
	out.PrintMarkedLine(os.Stderr, c.Source, token.LexemeLength(), token.Position, out.ColorBlue, '-')
	out.PrintHintMessage(message, out.ColorBlue)
}

func (c *Context) comparativeError(cause ast.Node, causeMessage string, where ast.Node, whereMessage string) {
	c.HadError = true

	c.nodeError(cause, causeMessage)
	c.nodeHelpHint(where, whereMessage)
}

func (c *Context) declare(declName string, declStaticType string, declType DeclarationType, node ast.Node) {
	prev, ok := c.top().Declared[declName]

	if ok {
		c.comparativeError(node, "Identifier already declared", prev.DeclaredAt, "Declared here")
		return
	}

	decl := declaration{DeclaredAt: node, Initialized: false, StaticType: declStaticType, Type: declType}
	c.top().Declared[declName] = decl
}

func (c *Context) define(name string, at ast.Node, value ast.Node) {
	scope := c.top()
	decl, ok := scope.Declared[name]

	if !ok {
		c.nodeError(at, "Cannot define undeclared identifier")
		return
	}

	decl.Initialized = true

	inferredValueType := fmt.Sprint(value.Visit(c))
	if len(decl.StaticType) == 0 {
		decl.StaticType = inferredValueType
	}

	if decl.StaticType != inferredValueType {
		c.nodeError(value, "Unexpected type")
		out.PrintHintMessage(fmt.Sprintf("Expcted value of type %s", inferredValueType), out.ColorRed)
		return
	}

	scope.Declared[name] = decl
}

func (c *Context) begin() {
	c.push(initScope())
}

func (c *Context) end() {
	_ = c.pop()
}

func (c *Context) VisitIdentifierLitExpr(node *ast.IdentifierLitExpr) any {
	name := node.Name

	scope := c.top()
	decl, ok := scope.Declared[name]

	if !ok {
		c.nodeError(node, "Undeclared identifier")
		return nil
	}

	if !decl.Initialized {
		c.nodeError(node, "Undefined identifier")
		return nil
	}

	return decl.StaticType
}
func (c *Context) VisitIntegerLitExpr(node *ast.IntegerLitExpr) any {
	return "int"
}
func (c *Context) VisitLetDecl(node *ast.LetDecl) any {
	declName := node.Identifier
	declType := node.Type
	c.declare(declName, declType, VariableDeclaration, node)
	return nil
}
func (c *Context) VisitBinaryExpr(node *ast.BinaryExpr) any {
	leftType := node.Left.Visit(c)
	rightType := node.Right.Visit(c)
	combinedType := leftType

	if leftType != rightType {
		c.nodeError(node, "Cannot combine types")
		out.PrintHintMessage(fmt.Sprintf("type %s != type %s", leftType, rightType), out.ColorRed)
		return nil
	}

	return combinedType
}
func (c *Context) VisitFloatingLitExpr(node *ast.FloatingLitExpr) any {
	return "float"
}
func (c *Context) VisitDebugStmt(node *ast.DebugStmt) any {
	return nil
}
func (c *Context) VisitBlockStmt(node *ast.BlockStmt) any {
	for _, n := range node.Nodes {
		n.Visit(c)
	}
	return nil
}
func (c *Context) VisitUnaryExpr(node *ast.UnaryExpr) any {
	if node.GetId() != ast.IntegerLitId && node.GetId() != ast.FloatingLitId {
		out.PrintErrorMessage("Unary not possible")
		return nil
	}
	return node.Expression.Visit(c)
}
func (c *Context) VisitAssignExpr(node *ast.AssignExpr) any {
	defName := node.Name.Lexeme
	defValue := node.Value
	c.define(defName, node, defValue)
	return defValue.Visit(c)
}
func (c *Context) VisitExprStmt(node *ast.ExprStmt) any {
	node.Expression.Visit(c)
	return nil
}
