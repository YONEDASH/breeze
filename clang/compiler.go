package clang

import (
	"breeze/ast"
	"breeze/common"
	"breeze/scanner"
	"bytes"
	"fmt"
	"os/exec"
)

func CompileClang(executablePath string, file common.SourceFile, nodes []ast.Node) string {
	source := CompileToSource(nodes)

	filePath := fmt.Sprintf("%s.c", file.Path[:len(file.Path)-3]) // substring .bz

	err := common.WriteFile(filePath, source)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("clang", "-o", executablePath, filePath)
	fmt.Println("%", cmd.String())

	var stdout = bytes.Buffer{}
	var stderr = bytes.Buffer{}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	fmt.Println(stderr.String())
	fmt.Println(stdout.String())

	if err != nil {
		fmt.Println(err)
		panic(err)
		return source
	}

	return source
}

func CompileToSource(nodes []ast.Node) string {
	c := &compiler{
		header: "",
		body:   "",
	}
	for _, node := range nodes {
		_ = node.Visit(c)
	}
	return fmt.Sprintf("%s\n%s", c.header, c.body)
}

type compiler struct {
	ast.Visitor
	header string
	body   string
}

func clangTypeName(name string) string {
	if len(name) == 0 {
		return "void"
	}

	return name
}

func (c *compiler) VisitDebugStmt(node *ast.DebugStmt) any {
	return nil
}
func (c *compiler) VisitFunctionDecl(node *ast.FunctionDecl) any {
	typeName := clangTypeName(node.ReturnType)

	c.body += typeName
	c.body += " "
	c.body += node.Identifier
	c.body += "("
	paramCount := len(node.ParamType)
	for i := 0; i < paramCount; i++ {
		c.body += clangTypeName(node.ParamType[i])
		c.body += " "
		c.body += node.ParamName[i]
		if i != paramCount-1 {
			c.body += ", "
		}
	}
	c.body += ")\n"

	_ = node.Closure.Visit(c)

	return nil
}
func (c *compiler) VisitConditionalStmt(node *ast.ConditionalStmt) any {
	c.body += "if ("
	_ = node.Condition.Visit(c)
	c.body += ")\n"
	_ = node.Statement.Visit(c)
	if node.ElseStatement != nil {
		c.body += "\nelse\n"
		_ = node.ElseStatement.Visit(c)
	}

	return nil
}
func (c *compiler) VisitLetDecl(node *ast.LetDecl) any {
	c.body += node.Type + " " + node.Identifier + ";\n"
	return nil
}
func (c *compiler) VisitWhileStmt(node *ast.WhileStmt) any { return nil }
func (c *compiler) VisitAssignExpr(node *ast.AssignExpr) any {
	c.body += "(" + node.Name.Lexeme + " = "
	node.Value.Visit(c)
	c.body += ")"

	return nil
}
func (c *compiler) VisitBinaryExpr(node *ast.BinaryExpr) any {
	c.body += "("
	_ = node.Left.Visit(c)

	switch node.Operator.Id {
	case scanner.Plus:
		c.body += "+"
	case scanner.Minus:
		c.body += "-"
	case scanner.Star:
		c.body += "*"
	case scanner.Slash:
		c.body += "/"
	case scanner.Lower:
		c.body += "<"
	case scanner.Greater:
		c.body += ">"
	case scanner.LowerEquals:
		c.body += "<="
	case scanner.GreaterEquals:
		c.body += ">="
	case scanner.EqualsEquals:
		c.body += "=="
	case scanner.BangEquals:
		c.body += "!="
	case scanner.AndAnd:
		c.body += "&&"
	case scanner.PipePipe:
		c.body += "||"
	default:
		panic(fmt.Sprintf("Missing binary operation translation for Clang: %d ", node.Operator.Id))
	}

	_ = node.Right.Visit(c)
	c.body += ")"

	return nil
}
func (c *compiler) VisitBlockStmt(node *ast.BlockStmt) any {
	for _, node := range node.Nodes {
		_ = node.Visit(c)
	}
	return nil
}
func (c *compiler) VisitReturnStmt(node *ast.ReturnStmt) any {
	c.body += "return"
	if node.Expression != nil {
		c.body += " "
		_ = node.Expression.Visit(c)
	}
	c.body += ";\n"
	return nil
}
func (c *compiler) VisitContinueStmt(node *ast.ContinueStmt) any {
	c.body += "continue;\n"
	return nil
}
func (c *compiler) VisitBreakStmt(node *ast.BreakStmt) any {
	c.body += "break;\n"
	return nil
}
func (c *compiler) VisitUnaryExpr(node *ast.UnaryExpr) any {
	c.body += "("
	switch node.Operator.Id {
	case scanner.Minus:
		c.body += "-"
	}
	node.Expression.Visit(c)
	c.body += ")"

	return nil
}
func (c *compiler) VisitFloatingLitExpr(node *ast.FloatingLitExpr) any {
	c.body += node.Value
	return nil
}
func (c *compiler) VisitClosureStmt(node *ast.ClosureStmt) any {
	c.body += "{\n"
	node.Block.Visit(c)
	c.body += "}\n"
	return nil
}
func (c *compiler) VisitCallExpr(node *ast.CallExpr) any {
	node.Expression.Visit(c)
	c.body += "("
	argCount := len(node.Arguments)
	for i, arg := range node.Arguments {
		_ = arg.Visit(c)
		if i != argCount-1 {
			c.body += ", "
		}
	}
	c.body += ")"
	return nil
}
func (c *compiler) VisitBooleanLitExpr(node *ast.BooleanLitExpr) any {
	c.body += node.Value
	return nil
}
func (c *compiler) VisitExprStmt(node *ast.ExprStmt) any {
	_ = node.Expression.Visit(c)
	c.body += ";\n"
	return nil
}
func (c *compiler) VisitIdentifierLitExpr(node *ast.IdentifierLitExpr) any {
	c.body += node.Name
	return nil
}
func (c *compiler) VisitErrNode(node *ast.ErrNode) any {
	// Should NEVER be called, maybe analysis stage missed?
	panic(node)
	return nil
}
func (c *compiler) VisitIntegerLitExpr(node *ast.IntegerLitExpr) any {
	c.body += node.Value
	return nil
}
