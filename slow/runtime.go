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
			_, _ = fmt.Println("Runtime Error:", name, "not found")
			return nil
		}
	}
	return val
}

func (e *Environment) set(name string, value any) {
	_, ok := e.Variables[name]
	if !ok {
		if e.Parent != nil {
			e.Parent.set(name, value)
		} else {
			_, _ = fmt.Println("Runtime Error:", name, "not found")
			panic(name)
		}
		return
	}
	e.Variables[name] = value
}

func initEnv(parent *Environment) *Environment {
	return &Environment{Parent: parent, Variables: make(map[string]any)}
}

var GlobalRuntime = Runtime{Current: initEnv(nil)}

func (r *Runtime) VisitLetDecl(node *ast.LetDecl) any {
	r.Current.Variables[node.Identifier] = nil
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

func (r *Runtime) VisitWhileStmt(node *ast.WhileStmt) any {
	for {
		result := node.Condition.Visit(r)

		if !isTrue(result) {
			break
		}

		_ = node.Statement.Visit(r)
	}

	return nil
}

func (r *Runtime) VisitConditionalStmt(node *ast.ConditionalStmt) any {
	result := node.Condition.Visit(r)

	if isTrue(result) {
		_ = node.Statement.Visit(r)
	} else if node.ElseStatement != nil {
		_ = node.ElseStatement.Visit(r)
	}

	return nil
}

func isTrue(a any) bool {
	return a != 0 && a != false
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

func getType(v interface{}) string {
	switch v.(type) {
	case float32:
		return "float32"
	case int:
		return "int"
	case bool:
		return "bool"
	}
	return ""
}

func (r *Runtime) VisitBinaryExpr(node *ast.BinaryExpr) any {
	left := node.Left.Visit(r)
	right := node.Right.Visit(r)

	leftType := getType(left)

	if leftType == "int" {
		switch node.Operator.Id {
		case scanner.Plus:
			return left.(int) + right.(int)
		case scanner.Minus:
			return left.(int) - right.(int)
		case scanner.Star:
			return left.(int) * right.(int)
		case scanner.Slash:
			return left.(int) / right.(int)
		case scanner.Lower:
			return left.(int) < right.(int)
		case scanner.Greater:
			return left.(int) > right.(int)
		case scanner.LowerEquals:
			return left.(int) <= right.(int)
		case scanner.GreaterEquals:
			return left.(int) >= right.(int)
		}
	} else if leftType == "float32" {
		switch node.Operator.Id {
		case scanner.Plus:
			return left.(float32) + right.(float32)
		case scanner.Minus:
			return left.(float32) - right.(float32)
		case scanner.Star:
			return left.(float32) * right.(float32)
		case scanner.Slash:
			return left.(float32) / right.(float32)
		case scanner.Lower:
			return left.(float32) < right.(float32)
		case scanner.Greater:
			return left.(float32) > right.(float32)
		case scanner.LowerEquals:
			return left.(float32) <= right.(float32)
		case scanner.GreaterEquals:
			return left.(float32) >= right.(float32)
		}
	}

	switch node.Operator.Id {
	case scanner.EqualsEquals:
		return left == right
	case scanner.BangEquals:
		return left != right
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

	valueType := getType(value)

	if valueType == "int" {
		switch node.Operator.Id {
		case scanner.Plus:
			return +value.(int)
		case scanner.Minus:
			return -value.(int)
		}
	} else if valueType == "float32" {
		switch node.Operator.Id {
		case scanner.Plus:
			return +value.(float32)
		case scanner.Minus:
			return -value.(float32)
		}
	} else if valueType == "bool" {
		switch node.Operator.Id {
		case scanner.Bang:
			return !value.(bool)
		}
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

func (r *Runtime) VisitBooleanLitExpr(node *ast.BooleanLitExpr) any {
	return node.Value == "true"
}

func (r *Runtime) VisitErrNode(node *ast.ErrNode) any {
	return nil
}

func (r *Runtime) VisitFloatingLitExpr(node *ast.FloatingLitExpr) any {
	f, _ := strconv.ParseFloat(node.Value, 32)
	return f
}
