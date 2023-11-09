package ast

import "breeze/scanner"

type NodeId uint8

const (
	IdentifierLitId NodeId = iota
	IntegerLitId
	LetId
	ErrId
	BinaryId
	FloatingLitId
	DebugId
	BlockId
	UnaryId
	AssignId
	ExprId
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
	VisitIdentifierLitExpr(node *IdentifierLitExpr) any
	VisitIntegerLitExpr(node *IntegerLitExpr) any
	VisitLetDecl(node *LetDecl) any
	VisitErrNode(node *ErrNode) any
	VisitBinaryExpr(node *BinaryExpr) any
	VisitFloatingLitExpr(node *FloatingLitExpr) any
	VisitDebugStmt(node *DebugStmt) any
	VisitBlockStmt(node *BlockStmt) any
	VisitUnaryExpr(node *UnaryExpr) any
	VisitAssignExpr(node *AssignExpr) any
	VisitExprStmt(node *ExprStmt) any
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

func (node *IdentifierLitExpr) Stringify() string {
	return "(IdentifierLitExpr Name=" + string(node.Name) + ")"
}

func (node *IdentifierLitExpr) GetToken() scanner.Token {
	return node.Token
}

func (node *IdentifierLitExpr) Visit(visitor Visitor) any {
	return visitor.VisitIdentifierLitExpr(node)
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

func (node *IntegerLitExpr) Stringify() string {
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

func (node *LetDecl) Stringify() string {
	return "(LetDecl StaticType=" + string(node.Type) + " Identifier=" + string(node.Identifier) + ")"
}

func (node *LetDecl) GetToken() scanner.Token {
	return node.Token
}

func (node *LetDecl) Visit(visitor Visitor) any {
	return visitor.VisitLetDecl(node)
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

type BinaryExpr struct {
	Node
	Right    Node
	Operator scanner.Token
	Left     Node
}

func (node *BinaryExpr) GetType() NodeType {
	return Expr
}

func (node *BinaryExpr) GetId() NodeId {
	return BinaryId
}

func (node *BinaryExpr) Stringify() string {
	return "(BinaryExpr Right=" + node.Right.Stringify() + " Operator=" + node.Operator.Stringify() + " Left=" + node.Left.Stringify() + ")"
}

func (node *BinaryExpr) GetToken() scanner.Token {
	return node.Operator
}

func (node *BinaryExpr) Visit(visitor Visitor) any {
	return visitor.VisitBinaryExpr(node)
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

func (node *FloatingLitExpr) Stringify() string {
	return "(FloatingLitExpr Value=" + string(node.Value) + ")"
}

func (node *FloatingLitExpr) GetToken() scanner.Token {
	return node.Token
}

func (node *FloatingLitExpr) Visit(visitor Visitor) any {
	return visitor.VisitFloatingLitExpr(node)
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

func (node *UnaryExpr) Stringify() string {
	return "(UnaryExpr Operator=" + node.Operator.Stringify() + " Expression=" + node.Expression.Stringify() + ")"
}

func (node *UnaryExpr) GetToken() scanner.Token {
	return node.Operator
}

func (node *UnaryExpr) Visit(visitor Visitor) any {
	return visitor.VisitUnaryExpr(node)
}

type AssignExpr struct {
	Node
	Operator scanner.Token
	Name     scanner.Token
	Value    Node
}

func (node *AssignExpr) GetType() NodeType {
	return Expr
}

func (node *AssignExpr) GetId() NodeId {
	return AssignId
}

func (node *AssignExpr) Stringify() string {
	return "(AssignExpr Operator=" + node.Operator.Stringify() + " Name=" + node.Name.Stringify() + " Value=" + node.Value.Stringify() + ")"
}

func (node *AssignExpr) GetToken() scanner.Token {
	return node.Operator
}

func (node *AssignExpr) Visit(visitor Visitor) any {
	return visitor.VisitAssignExpr(node)
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
