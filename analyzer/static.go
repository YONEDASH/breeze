package analyzer

import (
	"breeze/ast"
	"breeze/common"
	"breeze/out"
	"breeze/scanner"
	"fmt"
	"os"
)

type DeclarationType uint8

const (
	VariableDeclaration DeclarationType = iota
	FunctionDeclaration
)

type ReferenceType uint8

const (
	VariableReference ReferenceType = iota
	FunctionReference
	TypeReference
)

var initialNode = &ast.ErrNode{Token: scanner.Token{Id: scanner.EOF, Position: common.InitPosition()}, Message: "INITIAL", Hint: ""}

var (
	TypeNoReference    = &staticType{TypeName: "undef_type", DeclaredAt: initialNode}
	TypeVoidReference  = &staticType{TypeName: "void", DeclaredAt: initialNode}
	TypeIntReference   = &staticType{TypeName: "int", DeclaredAt: initialNode}
	TypeFloatReference = &staticType{TypeName: "float", DeclaredAt: initialNode}
	TypeBoolReference  = &staticType{TypeName: "bool", DeclaredAt: initialNode}
)

func declareTypes(context *Context) {
	context.declare(TypeNoReference, initialNode)
	context.declare(TypeVoidReference, initialNode)
	context.declare(TypeIntReference, initialNode)
	context.declare(TypeFloatReference, initialNode)
	context.declare(TypeBoolReference, initialNode)
}

func compareType(a staticType, b staticType) bool {
	return a.TypeName == b.TypeName
}

type staticDeclaration interface {
	RefType() ReferenceType
	Name() string
	Node() ast.Node
	Static() *staticType
}

type staticType struct {
	staticDeclaration
	DeclaredAt ast.Node
	TypeName   string
}

func (s *staticType) RefType() ReferenceType {
	return TypeReference
}

func (s *staticType) Name() string {
	return s.TypeName
}

func (s *staticType) Node() ast.Node {
	return s.DeclaredAt
}

func (s *staticType) Static() *staticType {
	return s
}

type variable struct {
	staticDeclaration
	DeclaredAt   ast.Node
	VariableName string
	VariableType *staticType
	Initialized  bool
}

func (v *variable) RefType() ReferenceType {
	return VariableReference
}

func (v *variable) Name() string {
	return v.VariableName
}

func (v *variable) Node() ast.Node {
	return v.DeclaredAt
}

func (v *variable) Static() *staticType {
	return v.VariableType
}

type function struct {
	staticDeclaration
	DeclaredAt     ast.Node
	FunctionName   string
	ReturnType     *staticType
	ParameterTypes []*staticType
}

func (f *function) RefType() ReferenceType {
	return FunctionReference
}

func (f *function) Name() string {
	return f.FunctionName
}

func (f *function) Node() ast.Node {
	return f.DeclaredAt
}

func (f *function) Static() *staticType {
	return f.ReturnType
}

type Scope struct {
	Declared map[string]staticDeclaration
}

func initScope() Scope {
	scope := Scope{
		Declared: make(map[string]staticDeclaration),
	}
	return scope
}

type Context struct {
	ast.Visitor
	File            common.SourceFile
	Source          string
	HadError        bool
	Stack           []Scope
	CurrentFunction *function
}

func Analyze(sourceFile common.SourceFile, source string, nodes []ast.Node) bool {
	context := &Context{Stack: make([]Scope, 0), HadError: false, CurrentFunction: nil, Source: source, File: sourceFile}
	context.begin()
	declareTypes(context)

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

func (c *Context) lookup(declName string) (staticDeclaration, bool) {
	size := len(c.Stack) - 1
	for i := size; i >= 0; i-- {
		decl, ok := c.Stack[i].Declared[declName]
		if ok {
			return decl, true
		}
	}
	return TypeVoidReference, false
}

func (c *Context) lookupType(node ast.Node, typeName string) (*staticType, bool) {
	if len(typeName) == 0 {
		return TypeNoReference, true
	}

	declType, ok := c.lookup(typeName)
	if !ok {
		c.nodeError(node, fmt.Sprintf("Undeclared type %s", typeName))
		return TypeVoidReference, false
	}
	if declType.RefType() != TypeReference {
		c.nodeError(node, "Invalid type")
		c.comparativeError(node, "Invalid type", declType.Node(), "This is not a type")
		return TypeVoidReference, false
	}
	staticDeclType := declType.(*staticType)
	return staticDeclType, true
}

func (c *Context) declare(staticDecl staticDeclaration, node ast.Node) {
	// Shadowed variables will be allowed (for now?)
	declName := staticDecl.Name()

	top := c.top()
	prev, ok := top.Declared[declName]
	if ok {
		c.comparativeError(node, "Already declared", prev.Node(), "Declared here")
		return
	}

	top.Declared[declName] = staticDecl
}

func (c *Context) define(name string, at ast.Node, value ast.Node) {
	decl, ok := c.lookup(name)

	if !ok {
		c.nodeError(at, "Cannot define undeclared identifier")
		return
	}

	if decl.RefType() == VariableReference {
		varDecl := decl.(*variable)
		varDecl.Initialized = true

		inferred := value.Visit(c).(staticDeclaration)
		if inferred.RefType() != TypeReference {
			c.nodeError(value, "Expected type reference")
			return
		}
		inferredType := inferred.(*staticType)
		if compareType(*varDecl.Static(), *TypeNoReference) {
			varDecl.VariableType = inferredType
		}

		if !compareType(*varDecl.Static(), *inferredType) {
			c.nodeError(value, "Unexpected type")
			out.PrintHintMessage(fmt.Sprintf("Expcted value of type %s", varDecl.VariableType.TypeName), out.ColorRed)
			return
		}

		// CONTEXT: Set type in node
		if decl.Node().GetId() == ast.LetId {
			letDecl := varDecl.DeclaredAt.(*ast.LetDecl)
			letDecl.Type = varDecl.VariableType.TypeName
		}

	}

	scope := c.top()
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

	decl, ok := c.lookup(name)

	if !ok {
		c.nodeError(node, "Undeclared identifier")
		return TypeVoidReference
	}

	if decl.RefType() == VariableReference {
		variable := decl.(*variable)
		if !variable.Initialized {
			c.nodeError(node, "Undefined variable")
			return TypeVoidReference
		}
		return variable
	}

	return decl
}
func (c *Context) VisitLetDecl(node *ast.LetDecl) any {
	declName := node.Identifier

	declType, ok := c.lookupType(node, node.Type)
	if !ok {
		return TypeVoidReference
	}
	decl := &variable{DeclaredAt: node, Initialized: false, VariableName: declName, VariableType: declType}
	c.declare(decl, node)
	return TypeVoidReference
}
func (c *Context) VisitFunctionDecl(node *ast.FunctionDecl) any {
	declName := node.Identifier
	declType, ok := c.lookupType(node, node.ReturnType)
	if !ok {
		return TypeVoidReference
	}

	paramCount := len(node.ParamType)
	parameterTypes := make([]*staticType, 0)

	for i := 0; i < paramCount; i++ {
		paramTypeName := node.ParamType[i]
		paramType, ok := c.lookupType(node, paramTypeName)
		if !ok {
			return TypeVoidReference
		}
		parameterTypes = append(parameterTypes, paramType)
	}

	fn := &function{DeclaredAt: node, FunctionName: declName, ReturnType: declType, ParameterTypes: parameterTypes}
	c.declare(fn, node)

	block := node.Closure.(*ast.ClosureStmt).Block

	c.begin()
	prev := c.CurrentFunction
	c.CurrentFunction = fn

	for i := 0; i < paramCount; i++ {
		paramName := node.ParamName[i]
		// Declare "initialized" variable
		decl := &variable{DeclaredAt: node, VariableType: parameterTypes[i], VariableName: paramName, Initialized: true}
		c.declare(decl, node)
	}

	_ = block.Visit(c)

	c.CurrentFunction = prev
	c.end()

	return TypeVoidReference
}
func (c *Context) VisitBinaryExpr(node *ast.BinaryExpr) any {
	leftType := node.Left.Visit(c).(staticDeclaration).Static()
	rightType := node.Right.Visit(c).(staticDeclaration).Static()
	combinedType := leftType

	if !compareType(*leftType, *rightType) {
		c.nodeError(node, "Type mismatch in binary expression")
		out.PrintHintMessage(fmt.Sprintf("type %s != type %s", leftType.Static().TypeName, rightType.Static().TypeName), out.ColorRed)
		return TypeVoidReference
	}

	switch node.Operator.Id {
	case scanner.Lower, scanner.Greater, scanner.LowerEquals, scanner.GreaterEquals, scanner.EqualsEquals, scanner.BangEquals:
		return TypeBoolReference
	}

	return combinedType
}

func (c *Context) VisitIntegerLitExpr(node *ast.IntegerLitExpr) any {
	return TypeIntReference
}

func (c *Context) VisitFloatingLitExpr(node *ast.FloatingLitExpr) any {
	return TypeFloatReference
}

func (c *Context) VisitBooleanLitExpr(node *ast.BooleanLitExpr) any {
	return TypeBoolReference
}

func (c *Context) VisitDebugStmt(node *ast.DebugStmt) any {
	_ = node.Expression.Visit(c)
	return TypeVoidReference
}

func (c *Context) VisitReturnStmt(node *ast.ReturnStmt) any {
	if c.CurrentFunction == nil {
		c.nodeError(node, "Cannot return outside of function")
		return TypeVoidReference
	}

	fn := c.CurrentFunction
	returnType := node.Expression.Visit(c).(staticDeclaration).Static()

	if !compareType(*returnType, *fn.ReturnType) {
		c.comparativeError(node, fmt.Sprintf("Invalid return type %s", returnType.TypeName), fn.Node(), fmt.Sprintf("Function expects return type of %s", fn.ReturnType.TypeName))
		return TypeVoidReference
	}

	// Stmt, return nothing
	return TypeVoidReference
}

func (c *Context) VisitContinueStmt(node *ast.ContinueStmt) any {
	return TypeVoidReference
}

func (c *Context) VisitBreakStmt(node *ast.BreakStmt) any {
	return TypeVoidReference
}

func (c *Context) VisitCallExpr(node *ast.CallExpr) any {
	exprDecl := node.Expression.Visit(c).(staticDeclaration)

	if exprDecl.RefType() != FunctionReference {
		c.nodeError(node.Expression, "Expected function")
		return TypeVoidReference
	}

	fn := exprDecl.(*function)

	paramCount := len(fn.ParameterTypes)
	argCount := len(node.Arguments)

	if argCount != paramCount {
		c.comparativeError(node, "Argument count mismatch", fn.Node(), fmt.Sprintf("Function has %d parameters", paramCount))
		return TypeVoidReference
	}

	for i := 0; i < paramCount; i++ {
		argType := node.Arguments[i].Visit(c).(staticDeclaration)
		expect := fn.ParameterTypes[i]
		if !compareType(*argType.Static(), *expect) {
			c.comparativeError(node, "Invalid argument type", fn.Node(), fmt.Sprintf("Function expects %s at position %d", expect.TypeName, i))
			return TypeVoidReference
		}
	}

	return fn.ReturnType
}

func (c *Context) VisitConditionalStmt(node *ast.ConditionalStmt) any {
	conditionType := node.Condition.Visit(c).(staticDeclaration)

	if !compareType(*conditionType.Static(), *TypeBoolReference) {
		c.comparativeError(node.Condition, "Unexpected condition type", node, fmt.Sprintf("Expected %s", TypeBoolReference.TypeName))
	}

	if node.Statement != nil {
		_ = node.Statement.Visit(c)
	}
	if node.ElseStatement != nil {
		_ = node.ElseStatement.Visit(c)
	}

	return TypeVoidReference
}

func (c *Context) VisitWhileStmt(node *ast.WhileStmt) any {
	conditionType := node.Condition.Visit(c)

	if conditionType != "bool" {
		c.comparativeError(node.Condition, "Unexpected condition type", node, "Expected bool type")
	}

	if node.Statement != nil {
		_ = node.Statement.Visit(c)
	}

	return TypeVoidReference
}

func (c *Context) VisitClosureStmt(node *ast.ClosureStmt) any {
	block := node.Block

	c.begin()
	_ = block.Visit(c)
	c.end()

	return TypeVoidReference
}

func (c *Context) VisitBlockStmt(node *ast.BlockStmt) any {
	for _, n := range node.Nodes {
		_ = n.Visit(c)
	}
	return TypeVoidReference
}

func (c *Context) VisitUnaryExpr(node *ast.UnaryExpr) any {
	exprType := node.Expression.Visit(c).(staticDeclaration).Static()
	switch node.Operator.Id {
	case scanner.Bang:
		if compareType(*exprType, *TypeBoolReference) {
			c.nodeError(node, "Unary operation possible on type bool")
		}

		break
	case scanner.Plus, scanner.Minus:
		if !compareType(*exprType, *TypeIntReference) && !compareType(*exprType, *TypeFloatReference) {
			c.nodeError(node, "Unary operation possible on types int and float")
		}

		break
	}
	return exprType
}

func (c *Context) VisitAssignExpr(node *ast.AssignExpr) any {
	defName := node.Name.Lexeme
	defValue := node.Value
	c.define(defName, node, defValue)
	return defValue.Visit(c)
}

func (c *Context) VisitExprStmt(node *ast.ExprStmt) any {
	_ = node.Expression.Visit(c)
	return TypeVoidReference
}

func (c *Context) VisitErrNode(node *ast.ErrNode) any {
	c.nodeError(node, fmt.Sprintf("Error node detected. %s", node.Message))
	return TypeVoidReference
}
