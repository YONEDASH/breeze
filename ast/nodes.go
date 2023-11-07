package ast

import "breeze/scanner"

type NodeId uint8

const (
	BinaryId NodeId = iota
	IdentifierId
	IntegerId
	ErrId
	FloatingId
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
	VisitBinaryExpr(node *BinaryExpr) any
	VisitIdentifierExpr(node *IdentifierExpr) any
	VisitIntegerExpr(node *IntegerExpr) any
	VisitErrNode(node *ErrNode) any
	VisitFloatingExpr(node *FloatingExpr) any
}

type BinaryExpr struct {
	Node
	Left     Node
	Operator scanner.Token
	Right    Node
}

func (node *BinaryExpr) GetType() NodeType {
	return Expr
}

func (node *BinaryExpr) GetId() NodeId {
	return BinaryId
}

func (node *BinaryExpr) Stringify() string {
	return "(BinaryExpr Left=" + node.Left.Stringify() + " Operator=" + node.Operator.Stringify() + " Right=" + node.Right.Stringify() + ")"
}

func (node *BinaryExpr) GetToken() scanner.Token {
	return node.Operator
}

func (node *BinaryExpr) Visit(visitor Visitor) any {
	return visitor.VisitBinaryExpr(node)
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

type ErrNode struct {
	Node
	Token   scanner.Token
	Message string
	Hint    string
}

func (node *ErrNode) GetType() NodeType {
	return Err
}

func (node *ErrNode) GetId() NodeId {
	return ErrId
}

func (node *ErrNode) Stringify() string {
	return "(ErrNode Message=" + string(node.Message) + " Hint=" + string(node.Hint) + ")"
}

func (node *ErrNode) GetToken() scanner.Token {
	return node.Token
}

func (node *ErrNode) Visit(visitor Visitor) any {
	return visitor.VisitErrNode(node)
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
