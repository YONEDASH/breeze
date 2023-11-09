package slow

import (
	"breeze/ast"
	"breeze/scanner"
	"fmt"
	"strconv"
)

// Language is in development. Many features may change.
// Would be too much work to implement a VM right now. Temporary GlobalRuntime written in Go.

type Runtime struct {
	ast.Visitor
	Declared  map[string]bool
	Variables map[string]any
}

var GlobalRuntime = Runtime{Declared: make(map[string]bool), Variables: make(map[string]any)}

func (r *Runtime) VisitLetDecl(node *ast.LetDecl) any {
	r.Declared[node.Identifier] = true
	return nil
}

func (r *Runtime) VisitDebugStmt(node *ast.DebugStmt) any {
	result := node.Expression.Visit(r)
	fmt.Println(result)
	return nil
}

func (r *Runtime) VisitExprStmt(node *ast.ExprStmt) any {
	node.Expression.Visit(r)
	return nil
}

func (r *Runtime) VisitBlockStmt(node *ast.BlockStmt) any {
	nodes := node.Nodes

	for _, n := range nodes {
		n.Visit(r)
	}

	return nil
}

func (r *Runtime) VisitAssignExpr(node *ast.AssignExpr) any {
	name := node.Name.Lexeme

	switch node.Operator.Id {
	case scanner.Equals:
		val := node.Value.Visit(r)
		GlobalRuntime.Variables[name] = val
		return val
	}

	return nil
}

func (r *Runtime) VisitBinaryExpr(node *ast.BinaryExpr) any {
	left := node.Left.Visit(r)
	right := node.Right.Visit(r)

	switch node.Operator.Id {
	case scanner.Plus:
		return left.(int) + right.(int)
	case scanner.Minus:
		return left.(int) - right.(int)
	case scanner.Star:
		return left.(int) * right.(int)
	case scanner.Slash:
		return left.(int) / right.(int)
	}

	return nil
}

func (r *Runtime) VisitIdentifierExpr(node *ast.IdentifierExpr) any {
	return r.Variables[node.Name]
}

func (r *Runtime) VisitIntegerExpr(node *ast.IntegerExpr) any {
	i, _ := strconv.Atoi(node.Value)
	return i
}

func (r *Runtime) VisitErrNode(node *ast.ErrNode) any {
	return nil
}

func (r *Runtime) VisitFloatingExpr(node *ast.FloatingExpr) any {
	return nil
}
