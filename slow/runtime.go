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
	Current *Environment
}

type Environment struct {
	Parent    *Environment
	Variables map[string]any
}

func (e *Environment) get(name string) any {
	val, ok := e.Variables[name]
	if !ok {
		if e.Parent != nil {
			return e.Parent.get(name)
		} else {
			fmt.Println("Runtime Error:", name, "not found")
			return nil
		}
	}
	return val
}

func (e *Environment) set(name string, value any) {
	e.Variables[name] = value
}

func initEnv(parent *Environment) *Environment {
	return &Environment{Parent: parent, Variables: make(map[string]any)}
}

var GlobalRuntime = Runtime{Current: initEnv(nil)}

func (r *Runtime) VisitLetDecl(node *ast.LetDecl) any {
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

func (r *Runtime) VisitClosureStmt(node *ast.ClosureStmt) any {
	block := node.Block

	before := r.Current
	r.Current = initEnv(r.Current)
	_ = block.Visit(r)
	r.Current = before

	return nil
}

func (r *Runtime) VisitBlockStmt(node *ast.BlockStmt) any {
	nodes := node.Nodes

	for _, n := range nodes {
		_ = n.Visit(r)
	}

	return nil
}

func (r *Runtime) VisitAssignExpr(node *ast.AssignExpr) any {
	name := node.Name.Lexeme

	switch node.Operator.Id {
	case scanner.Equals:
		val := node.Value.Visit(r)
		GlobalRuntime.Current.set(name, val)
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

func (r *Runtime) VisitUnaryExpr(node *ast.UnaryExpr) any {
	value := node.Expression.Visit(r)

	switch node.Operator.Id {
	case scanner.Plus:
		return +value.(int)
	case scanner.Minus:
		return -value.(int)
	}

	return nil
}

func (r *Runtime) VisitIdentifierLitExpr(node *ast.IdentifierLitExpr) any {
	return r.Current.get(node.Name)
}

func (r *Runtime) VisitIntegerLitExpr(node *ast.IntegerLitExpr) any {
	i, _ := strconv.Atoi(node.Value)
	return i
}

func (r *Runtime) VisitErrNode(node *ast.ErrNode) any {
	return nil
}

func (r *Runtime) VisitFloatingLitExpr(node *ast.FloatingLitExpr) any {
	return nil
}
