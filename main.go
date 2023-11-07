package main

import (
	"breeze/ast"
	"breeze/common"
	"breeze/out"
	"breeze/parser"
	"breeze/scanner"
	"fmt"
	"os"
	"strconv"
)

func main() {
	file := common.InitSource("test/breeze.bz")

	err := file.Validate()
	if err != nil {
		out.PrintErrorMessage(fmt.Sprintf("Could not validate path %s: %s", file.Path, err.Error()))
		os.Exit(out.ExOsFile)
	}

	source, err := file.GetContent()

	if err != nil {
		out.PrintErrorMessage(fmt.Sprintf("Could not read %s", file.Path))
		os.Exit(out.ExOsFile)
	}

	tokens, hadError := scanner.Scan(&file, source)

	if hadError {
		out.PrintErrorMessage("Scanning phase failed")
		os.Exit(out.ExDataErr)
		return
	}

	for _, tk := range tokens {
		fmt.Println(tk.Stringify())
	}

	nodes, hadError := parser.ParseTokens(file, source, tokens)

	if hadError {
		out.PrintErrorMessage("Parsing phase failed")
		os.Exit(out.ExDataErr)
		return
	}

	for _, node := range nodes {
		fmt.Println(node.Stringify())
	}

	fmt.Println("Result:", debugInterpret(nodes[0]))

}

func debugInterpret(node ast.Node) int {
	if node.GetId() == ast.IntegerId {
		i, _ := strconv.Atoi(node.(*ast.IntegerExpr).Value)
		fmt.Println("ret int", i)
		return i
	} else if node.GetId() == ast.BinaryId {
		binary := node.(*ast.BinaryExpr)

		fmt.Println("binary l,r op " + binary.Operator.Lexeme)
		left := debugInterpret(binary.Left)
		fmt.Println("l =", left)
		right := debugInterpret(binary.Right)
		fmt.Println("r =", right)

		switch binary.Operator.Id {
		case scanner.Plus:
			return left + right
		case scanner.Minus:
			return left - right
		case scanner.Star:
			return left * right
		case scanner.Slash:
			return left / right
		}

		fmt.Println("unknown op")
		return 0
	}

	fmt.Println("unknown expr")
	return 0
}
