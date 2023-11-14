package ast

import (
	"breeze/scanner"
	"fmt"
)

type NodeId uint8

const (
	ConditionalId NodeId = iota
	WhileId
	ClosureId
	ExprId
	AssignId
	ErrId
	LetId
	BinaryId
	UnaryId
	FunctionId
	CallId
	DebugId
	ReturnId
	IdentifierLitId
	ContinueId
	BreakId
	IntegerLitId
	BlockId
	FloatingLitId
	BooleanLitId
)

type NodeType uint8

const (
	Err NodeType = iota
	Expr
	Stmt
	Decl
)

type Node interface {
	GetId() NodeId
	GetType() NodeType
	String() string
	GetToken() scanner.Token
	Visit(visitor Visitor) any
}

type Visitor interface {
	VisitConditionalStmt(node *ConditionalStmt) any
	VisitWhileStmt(node *WhileStmt) any
	VisitClosureStmt(node *ClosureStmt) any
	VisitExprStmt(node *ExprStmt) any
	VisitAssignExpr(node *AssignExpr) any
	VisitErrNode(node *ErrNode) any
	VisitLetDecl(node *LetDecl) any
	VisitBinaryExpr(node *BinaryExpr) any
	VisitUnaryExpr(node *UnaryExpr) any
	VisitFunctionDecl(node *FunctionDecl) any
	VisitCallExpr(node *CallExpr) any
	VisitDebugStmt(node *DebugStmt) any
	VisitReturnStmt(node *ReturnStmt) any
	VisitIdentifierLitExpr(node *IdentifierLitExpr) any
	VisitContinueStmt(node *ContinueStmt) any
	VisitBreakStmt(node *BreakStmt) any
	VisitIntegerLitExpr(node *IntegerLitExpr) any
	VisitBlockStmt(node *BlockStmt) any
	VisitFloatingLitExpr(node *FloatingLitExpr) any
	VisitBooleanLitExpr(node *BooleanLitExpr) any
}

type ConditionalStmt struct {
	Node
	Token         scanner.Token
	Condition     Node
	Statement     Node
	ElseStatement Node
}

func (node *ConditionalStmt) GetType() NodeType {
	return Stmt
}

func (node *ConditionalStmt) GetId() NodeId {
	return ConditionalId
}

func (node *ConditionalStmt) String() string {
	return "(ConditionalStmt Condition=" + fmt.Sprintf("%s", node.Condition) + " Statement=" + fmt.Sprintf("%s", node.Statement) + " ElseStatement=" + fmt.Sprintf("%s", node.ElseStatement) + ")"
}

func (node *ConditionalStmt) GetToken() scanner.Token {
	return node.Token
}

func (node *ConditionalStmt) Visit(visitor Visitor) any {
	return visitor.VisitConditionalStmt(node)
}

type WhileStmt struct {
	Node
	Token     scanner.Token
	Condition Node
	Statement Node
}

func (node *WhileStmt) GetType() NodeType {
	return Stmt
}

func (node *WhileStmt) GetId() NodeId {
	return WhileId
}

func (node *WhileStmt) String() string {
	return "(WhileStmt Condition=" + fmt.Sprintf("%s", node.Condition) + " Statement=" + fmt.Sprintf("%s", node.Statement) + ")"
}

func (node *WhileStmt) GetToken() scanner.Token {
	return node.Token
}

func (node *WhileStmt) Visit(visitor Visitor) any {
	return visitor.VisitWhileStmt(node)
}

type ClosureStmt struct {
	Node
	Token scanner.Token
	Block Node
}

func (node *ClosureStmt) GetType() NodeType {
	return Stmt
}

func (node *ClosureStmt) GetId() NodeId {
	return ClosureId
}

func (node *ClosureStmt) String() string {
	return "(ClosureStmt Block=" + fmt.Sprintf("%s", node.Block) + ")"
}

func (node *ClosureStmt) GetToken() scanner.Token {
	return node.Token
}

func (node *ClosureStmt) Visit(visitor Visitor) any {
	return visitor.VisitClosureStmt(node)
}

type ExprStmt struct {
	Node
	Token      scanner.Token
	Expression Node
}

func (node *ExprStmt) GetType() NodeType {
	return Stmt
}

func (node *ExprStmt) GetId() NodeId {
	return ExprId
}

func (node *ExprStmt) String() string {
	return "(ExprStmt Expression=" + fmt.Sprintf("%s", node.Expression) + ")"
}

func (node *ExprStmt) GetToken() scanner.Token {
	return node.Token
}

func (node *ExprStmt) Visit(visitor Visitor) any {
	return visitor.VisitExprStmt(node)
}

type AssignExpr struct {
	Node
	Name     scanner.Token
	Value    Node
	Operator scanner.Token
}

func (node *AssignExpr) GetType() NodeType {
	return Expr
}

func (node *AssignExpr) GetId() NodeId {
	return AssignId
}

func (node *AssignExpr) String() string {
	return "(AssignExpr Name=" + fmt.Sprintf("%s", node.Name) + " Value=" + fmt.Sprintf("%s", node.Value) + " Operator=" + fmt.Sprintf("%s", node.Operator) + ")"
}

func (node *AssignExpr) GetToken() scanner.Token {
	return node.Name
}

func (node *AssignExpr) Visit(visitor Visitor) any {
	return visitor.VisitAssignExpr(node)
}

type ErrNode struct {
	Node
	Token   scanner.Token
	Hint    string
	Message string
}

func (node *ErrNode) GetType() NodeType {
	return Err
}

func (node *ErrNode) GetId() NodeId {
	return ErrId
}

func (node *ErrNode) String() string {
	return "(ErrNode Hint=" + string(node.Hint) + " Message=" + string(node.Message) + ")"
}

func (node *ErrNode) GetToken() scanner.Token {
	return node.Token
}

func (node *ErrNode) Visit(visitor Visitor) any {
	return visitor.VisitErrNode(node)
}

type LetDecl struct {
	Node
	Token      scanner.Token
	Identifier string
	Type       string
}

func (node *LetDecl) GetType() NodeType {
	return Decl
}

func (node *LetDecl) GetId() NodeId {
	return LetId
}

func (node *LetDecl) String() string {
	return "(LetDecl Identifier=" + string(node.Identifier) + " Type=" + string(node.Type) + ")"
}

func (node *LetDecl) GetToken() scanner.Token {
	return node.Token
}

func (node *LetDecl) Visit(visitor Visitor) any {
	return visitor.VisitLetDecl(node)
}

type BinaryExpr struct {
	Node
	Right    Node
	Left     Node
	Operator scanner.Token
}

func (node *BinaryExpr) GetType() NodeType {
	return Expr
}

func (node *BinaryExpr) GetId() NodeId {
	return BinaryId
}

func (node *BinaryExpr) String() string {
	return "(BinaryExpr Right=" + fmt.Sprintf("%s", node.Right) + " Left=" + fmt.Sprintf("%s", node.Left) + " Operator=" + fmt.Sprintf("%s", node.Operator) + ")"
}

func (node *BinaryExpr) GetToken() scanner.Token {
	return node.Operator
}

func (node *BinaryExpr) Visit(visitor Visitor) any {
	return visitor.VisitBinaryExpr(node)
}

type UnaryExpr struct {
	Node
	Operator   scanner.Token
	Expression Node
}

func (node *UnaryExpr) GetType() NodeType {
	return Expr
}

func (node *UnaryExpr) GetId() NodeId {
	return UnaryId
}

func (node *UnaryExpr) String() string {
	return "(UnaryExpr Operator=" + fmt.Sprintf("%s", node.Operator) + " Expression=" + fmt.Sprintf("%s", node.Expression) + ")"
}

func (node *UnaryExpr) GetToken() scanner.Token {
	return node.Operator
}

func (node *UnaryExpr) Visit(visitor Visitor) any {
	return visitor.VisitUnaryExpr(node)
}

type FunctionDecl struct {
	Node
	Token      scanner.Token
	Closure    Node
	ReturnType string
	ParamType  []string
	Identifier string
	ParamName  []string
}

func (node *FunctionDecl) GetType() NodeType {
	return Decl
}

func (node *FunctionDecl) GetId() NodeId {
	return FunctionId
}

func (node *FunctionDecl) String() string {
	str_ParamType := "{"
	for i, n := range node.ParamType {
		str_ParamType += fmt.Sprintf("%s", n)
		if i <= len(node.ParamType)-1 {
			str_ParamType += ", "
		}
	}
	str_ParamType += "}"
	str_ParamName := "{"
	for i, n := range node.ParamName {
		str_ParamName += fmt.Sprintf("%s", n)
		if i <= len(node.ParamName)-1 {
			str_ParamName += ", "
		}
	}
	str_ParamName += "}"
	return "(FunctionDecl Closure=" + fmt.Sprintf("%s", node.Closure) + " ReturnType=" + string(node.ReturnType) + " ParamType=" + str_ParamType + " Identifier=" + string(node.Identifier) + " ParamName=" + str_ParamName + ")"
}

func (node *FunctionDecl) GetToken() scanner.Token {
	return node.Token
}

func (node *FunctionDecl) Visit(visitor Visitor) any {
	return visitor.VisitFunctionDecl(node)
}

type CallExpr struct {
	Node
	Token      scanner.Token
	Arguments  []Node
	Expression Node
}

func (node *CallExpr) GetType() NodeType {
	return Expr
}

func (node *CallExpr) GetId() NodeId {
	return CallId
}

func (node *CallExpr) String() string {
	str_Arguments := "{"
	for i, n := range node.Arguments {
		str_Arguments += fmt.Sprintf("%s", n)
		if i <= len(node.Arguments)-1 {
			str_Arguments += ", "
		}
	}
	str_Arguments += "}"
	return "(CallExpr Arguments=" + str_Arguments + " Expression=" + fmt.Sprintf("%s", node.Expression) + ")"
}

func (node *CallExpr) GetToken() scanner.Token {
	return node.Token
}

func (node *CallExpr) Visit(visitor Visitor) any {
	return visitor.VisitCallExpr(node)
}

type DebugStmt struct {
	Node
	Token      scanner.Token
	Expression Node
}

func (node *DebugStmt) GetType() NodeType {
	return Stmt
}

func (node *DebugStmt) GetId() NodeId {
	return DebugId
}

func (node *DebugStmt) String() string {
	return "(DebugStmt Expression=" + fmt.Sprintf("%s", node.Expression) + ")"
}

func (node *DebugStmt) GetToken() scanner.Token {
	return node.Token
}

func (node *DebugStmt) Visit(visitor Visitor) any {
	return visitor.VisitDebugStmt(node)
}

type ReturnStmt struct {
	Node
	Token      scanner.Token
	Expression Node
}

func (node *ReturnStmt) GetType() NodeType {
	return Stmt
}

func (node *ReturnStmt) GetId() NodeId {
	return ReturnId
}

func (node *ReturnStmt) String() string {
	return "(ReturnStmt Expression=" + fmt.Sprintf("%s", node.Expression) + ")"
}

func (node *ReturnStmt) GetToken() scanner.Token {
	return node.Token
}

func (node *ReturnStmt) Visit(visitor Visitor) any {
	return visitor.VisitReturnStmt(node)
}

type IdentifierLitExpr struct {
	Node
	Token scanner.Token
	Name  string
}

func (node *IdentifierLitExpr) GetType() NodeType {
	return Expr
}

func (node *IdentifierLitExpr) GetId() NodeId {
	return IdentifierLitId
}

func (node *IdentifierLitExpr) String() string {
	return "(IdentifierLitExpr Name=" + string(node.Name) + ")"
}

func (node *IdentifierLitExpr) GetToken() scanner.Token {
	return node.Token
}

func (node *IdentifierLitExpr) Visit(visitor Visitor) any {
	return visitor.VisitIdentifierLitExpr(node)
}

type ContinueStmt struct {
	Node
	Token scanner.Token
}

func (node *ContinueStmt) GetType() NodeType {
	return Stmt
}

func (node *ContinueStmt) GetId() NodeId {
	return ContinueId
}

func (node *ContinueStmt) String() string {
	return "(ContinueStmt)"
}

func (node *ContinueStmt) GetToken() scanner.Token {
	return node.Token
}

func (node *ContinueStmt) Visit(visitor Visitor) any {
	return visitor.VisitContinueStmt(node)
}

type BreakStmt struct {
	Node
	Token scanner.Token
}

func (node *BreakStmt) GetType() NodeType {
	return Stmt
}

func (node *BreakStmt) GetId() NodeId {
	return BreakId
}

func (node *BreakStmt) String() string {
	return "(BreakStmt)"
}

func (node *BreakStmt) GetToken() scanner.Token {
	return node.Token
}

func (node *BreakStmt) Visit(visitor Visitor) any {
	return visitor.VisitBreakStmt(node)
}

type IntegerLitExpr struct {
	Node
	Token scanner.Token
	Value string
}

func (node *IntegerLitExpr) GetType() NodeType {
	return Expr
}

func (node *IntegerLitExpr) GetId() NodeId {
	return IntegerLitId
}

func (node *IntegerLitExpr) String() string {
	return "(IntegerLitExpr Value=" + string(node.Value) + ")"
}

func (node *IntegerLitExpr) GetToken() scanner.Token {
	return node.Token
}

func (node *IntegerLitExpr) Visit(visitor Visitor) any {
	return visitor.VisitIntegerLitExpr(node)
}

type BlockStmt struct {
	Node
	Token scanner.Token
	Nodes []Node
}

func (node *BlockStmt) GetType() NodeType {
	return Stmt
}

func (node *BlockStmt) GetId() NodeId {
	return BlockId
}

func (node *BlockStmt) String() string {
	str_Nodes := "{"
	for i, n := range node.Nodes {
		str_Nodes += fmt.Sprintf("%s", n)
		if i <= len(node.Nodes)-1 {
			str_Nodes += ", "
		}
	}
	str_Nodes += "}"
	return "(BlockStmt Nodes=" + str_Nodes + ")"
}

func (node *BlockStmt) GetToken() scanner.Token {
	return node.Token
}

func (node *BlockStmt) Visit(visitor Visitor) any {
	return visitor.VisitBlockStmt(node)
}

type FloatingLitExpr struct {
	Node
	Token scanner.Token
	Value string
}

func (node *FloatingLitExpr) GetType() NodeType {
	return Expr
}

func (node *FloatingLitExpr) GetId() NodeId {
	return FloatingLitId
}

func (node *FloatingLitExpr) String() string {
	return "(FloatingLitExpr Value=" + string(node.Value) + ")"
}

func (node *FloatingLitExpr) GetToken() scanner.Token {
	return node.Token
}

func (node *FloatingLitExpr) Visit(visitor Visitor) any {
	return visitor.VisitFloatingLitExpr(node)
}

type BooleanLitExpr struct {
	Node
	Token scanner.Token
	Value string
}

func (node *BooleanLitExpr) GetType() NodeType {
	return Expr
}

func (node *BooleanLitExpr) GetId() NodeId {
	return BooleanLitId
}

func (node *BooleanLitExpr) String() string {
	return "(BooleanLitExpr Value=" + string(node.Value) + ")"
}

func (node *BooleanLitExpr) GetToken() scanner.Token {
	return node.Token
}

func (node *BooleanLitExpr) Visit(visitor Visitor) any {
	return visitor.VisitBooleanLitExpr(node)
}
