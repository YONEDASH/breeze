package slow

import (
	"breeze/ast"
	"breeze/scanner"
	"strconv"
)

// Language is in development. Many features may change.
// Would be too much work to implement a VM right now. Temporary runtime written in Go.

type Runtime struct {
	ast.Visitor
}

var runtime Runtime

func (r *Runtime) VisitBinaryExpr(node *ast.BinaryExpr) any {
	left := node.Left.Visit(&runtime)
	right := node.Left.Visit(&runtime)

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
	return nil
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
