package ast

import "breeze/scanner"

type NodeId uint8

const (
	BinaryId NodeId = iota
	IntegerId
	FloatingId
	IdentifierId
)

type NodeType uint8

const (
	Expr NodeType = iota
	Stmt
	Decl
)

type Node interface {
	GetId() NodeId
	GetType() NodeType
	Stringify() string
	GetToken() scanner.Token
}

type BinaryExpr struct {
	Node
	Operator scanner.Token
}

func (node *BinaryExpr) GetType() NodeType {
	return Expr
}

func (node *BinaryExpr) GetId() NodeId {
	return BinaryId
}

func (node *BinaryExpr) Stringify() string {
	return "(BinaryExpr Operator=" + node.Operator.Stringify() + ")"
}

func (node *BinaryExpr) GetToken() scanner.Token {
	return node.Operator
}

type IntegerExpr struct {
	Node
	token scanner.Token
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
	return node.token
}

type FloatingExpr struct {
	Node
	token scanner.Token
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
	return node.token
}

type IdentifierExpr struct {
	Node
	token scanner.Token
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
	return node.token
}
