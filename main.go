package main

import (
	"breeze/analyzer"
	"breeze/common"
	"breeze/out"
	"breeze/parser"
	"breeze/scanner"
	"breeze/slow"
	"fmt"
	"os"
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

	hadError = analyzer.Analyze(file, source, nodes)

	if hadError {
		out.PrintErrorMessage("Static analyzing phase failed")
		os.Exit(out.ExDataErr)
		return
	}

	runtime := &slow.GlobalRuntime
	for _, node := range nodes {
		fmt.Print(out.ColorWhite.S())
		fmt.Println(node.Stringify())
		fmt.Print(out.ColorReset.S())
		node.Visit(runtime)
	}

}
