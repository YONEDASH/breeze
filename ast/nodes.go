package ast

import (
	"breeze/scanner"
	"fmt"
)

type NodeId uint8

const (
	BlockId NodeId = iota
	ExprId
	UnaryId
	IdentifierLitId
	BooleanLitId
	IntegerLitId
	LetId
	BinaryId
	ConditionalId
	FloatingLitId
	ErrId
	AssignId
	ClosureId
	DebugId
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
	VisitBlockStmt(node *BlockStmt) any
	VisitExprStmt(node *ExprStmt) any
	VisitUnaryExpr(node *UnaryExpr) any
	VisitIdentifierLitExpr(node *IdentifierLitExpr) any
	VisitBooleanLitExpr(node *BooleanLitExpr) any
	VisitIntegerLitExpr(node *IntegerLitExpr) any
	VisitLetDecl(node *LetDecl) any
	VisitBinaryExpr(node *BinaryExpr) any
	VisitConditionalStmt(node *ConditionalStmt) any
	VisitFloatingLitExpr(node *FloatingLitExpr) any
	VisitErrNode(node *ErrNode) any
	VisitAssignExpr(node *AssignExpr) any
	VisitClosureStmt(node *ClosureStmt) any
	VisitDebugStmt(node *DebugStmt) any
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
	strNodes := "{"
	for i, n := range node.Nodes {
		strNodes += fmt.Sprintf("%s", n)
		if i <= len(node.Nodes)-1 {
			strNodes += ", "
		}
	}
	strNodes += "}"
	return "(BlockStmt Nodes=" + strNodes + ")"
}

func (node *BlockStmt) GetToken() scanner.Token {
	return node.Token
}

func (node *BlockStmt) Visit(visitor Visitor) any {
	return visitor.VisitBlockStmt(node)
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

type LetDecl struct {
	Node
	Token      scanner.Token
	Type       string
	Identifier string
}

func (node *LetDecl) GetType() NodeType {
	return Decl
}

func (node *LetDecl) GetId() NodeId {
	return LetId
}

func (node *LetDecl) String() string {
	return "(LetDecl Type=" + string(node.Type) + " Identifier=" + string(node.Identifier) + ")"
}

func (node *LetDecl) GetToken() scanner.Token {
	return node.Token
}

func (node *LetDecl) Visit(visitor Visitor) any {
	return visitor.VisitLetDecl(node)
}

type BinaryExpr struct {
	Node
	Operator scanner.Token
	Right    Node
	Left     Node
}

func (node *BinaryExpr) GetType() NodeType {
	return Expr
}

func (node *BinaryExpr) GetId() NodeId {
	return BinaryId
}

func (node *BinaryExpr) String() string {
	return "(BinaryExpr Operator=" + fmt.Sprintf("%s", node.Operator) + " Right=" + fmt.Sprintf("%s", node.Right) + " Left=" + fmt.Sprintf("%s", node.Left) + ")"
}

func (node *BinaryExpr) GetToken() scanner.Token {
	return node.Operator
}

func (node *BinaryExpr) Visit(visitor Visitor) any {
	return visitor.VisitBinaryExpr(node)
}

type ConditionalStmt struct {
	Node
	Token         scanner.Token
	Statement     Node
	ElseStatement Node
	Condition     Node
}

func (node *ConditionalStmt) GetType() NodeType {
	return Stmt
}

func (node *ConditionalStmt) GetId() NodeId {
	return ConditionalId
}

func (node *ConditionalStmt) String() string {
	return "(ConditionalStmt Statement=" + fmt.Sprintf("%s", node.Statement) + " ElseStatement=" + fmt.Sprintf("%s", node.ElseStatement) + " Condition=" + fmt.Sprintf("%s", node.Condition) + ")"
}

func (node *ConditionalStmt) GetToken() scanner.Token {
	return node.Token
}

func (node *ConditionalStmt) Visit(visitor Visitor) any {
	return visitor.VisitConditionalStmt(node)
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

type AssignExpr struct {
	Node
	Operator scanner.Token
	Value    Node
	Name     scanner.Token
}

func (node *AssignExpr) GetType() NodeType {
	return Expr
}

func (node *AssignExpr) GetId() NodeId {
	return AssignId
}

func (node *AssignExpr) String() string {
	return "(AssignExpr Operator=" + fmt.Sprintf("%s", node.Operator) + " Value=" + fmt.Sprintf("%s", node.Value) + " Name=" + fmt.Sprintf("%s", node.Name) + ")"
}

func (node *AssignExpr) GetToken() scanner.Token {
	return node.Operator
}

func (node *AssignExpr) Visit(visitor Visitor) any {
	return visitor.VisitAssignExpr(node)
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
