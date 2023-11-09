package ast

import "breeze/scanner"

type NodeId uint8

const (
	AssignId NodeId = iota
	IntegerId
	ExprId
	FloatingId
	ErrId
	LetId
	BinaryId
	DebugId
	BlockId
	IdentifierId
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
	Stringify() string
	GetToken() scanner.Token
	Visit(visitor Visitor) any
}

type Visitor interface {
	VisitAssignExpr(node *AssignExpr) any
	VisitIntegerExpr(node *IntegerExpr) any
	VisitExprStmt(node *ExprStmt) any
	VisitFloatingExpr(node *FloatingExpr) any
	VisitErrNode(node *ErrNode) any
	VisitLetDecl(node *LetDecl) any
	VisitBinaryExpr(node *BinaryExpr) any
	VisitDebugStmt(node *DebugStmt) any
	VisitBlockStmt(node *BlockStmt) any
	VisitIdentifierExpr(node *IdentifierExpr) any
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

func (node *AssignExpr) Stringify() string {
	return "(AssignExpr Name=" + node.Name.Stringify() + " Value=" + node.Value.Stringify() + " Operator=" + node.Operator.Stringify() + ")"
}

func (node *AssignExpr) GetToken() scanner.Token {
	return node.Name
}

func (node *AssignExpr) Visit(visitor Visitor) any {
	return visitor.VisitAssignExpr(node)
}

type IntegerExpr struct {
	Node
	Token scanner.Token
	Value string
}

func (node *IntegerExpr) GetType() NodeType {
	return Expr
}

func (node *IntegerExpr) GetId() NodeId {
	return IntegerId
}

func (node *IntegerExpr) Stringify() string {
	return "(IntegerExpr Value=" + string(node.Value) + ")"
}

func (node *IntegerExpr) GetToken() scanner.Token {
	return node.Token
}

func (node *IntegerExpr) Visit(visitor Visitor) any {
	return visitor.VisitIntegerExpr(node)
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

func (node *ExprStmt) Stringify() string {
	return "(ExprStmt Expression=" + node.Expression.Stringify() + ")"
}

func (node *ExprStmt) GetToken() scanner.Token {
	return node.Token
}

func (node *ExprStmt) Visit(visitor Visitor) any {
	return visitor.VisitExprStmt(node)
}

type FloatingExpr struct {
	Node
	Token scanner.Token
	Value string
}

func (node *FloatingExpr) GetType() NodeType {
	return Expr
}

func (node *FloatingExpr) GetId() NodeId {
	return FloatingId
}

func (node *FloatingExpr) Stringify() string {
	return "(FloatingExpr Value=" + string(node.Value) + ")"
}

func (node *FloatingExpr) GetToken() scanner.Token {
	return node.Token
}

func (node *FloatingExpr) Visit(visitor Visitor) any {
	return visitor.VisitFloatingExpr(node)
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

func (node *ErrNode) Stringify() string {
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

func (node *LetDecl) Stringify() string {
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
	Left     Node
	Right    Node
	Operator scanner.Token
}

func (node *BinaryExpr) GetType() NodeType {
	return Expr
}

func (node *BinaryExpr) GetId() NodeId {
	return BinaryId
}

func (node *BinaryExpr) Stringify() string {
	return "(BinaryExpr Left=" + node.Left.Stringify() + " Right=" + node.Right.Stringify() + " Operator=" + node.Operator.Stringify() + ")"
}

func (node *BinaryExpr) GetToken() scanner.Token {
	return node.Operator
}

func (node *BinaryExpr) Visit(visitor Visitor) any {
	return visitor.VisitBinaryExpr(node)
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

func (node *DebugStmt) Stringify() string {
	return "(DebugStmt Expression=" + node.Expression.Stringify() + ")"
}

func (node *DebugStmt) GetToken() scanner.Token {
	return node.Token
}

func (node *DebugStmt) Visit(visitor Visitor) any {
	return visitor.VisitDebugStmt(node)
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

func (node *BlockStmt) Stringify() string {
	strNodes := "{"
	for i, n := range node.Nodes {
		strNodes += n.Stringify()
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

type IdentifierExpr struct {
	Node
	Token scanner.Token
	Name  string
}

func (node *IdentifierExpr) GetType() NodeType {
	return Expr
}

func (node *IdentifierExpr) GetId() NodeId {
	return IdentifierId
}

func (node *IdentifierExpr) Stringify() string {
	return "(IdentifierExpr Name=" + string(node.Name) + ")"
}

func (node *IdentifierExpr) GetToken() scanner.Token {
	return node.Token
}

func (node *IdentifierExpr) Visit(visitor Visitor) any {
	return visitor.VisitIdentifierExpr(node)
}
