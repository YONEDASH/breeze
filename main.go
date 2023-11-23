package main

import (
	"breeze/analyzer"
	"breeze/clang"
	"breeze/common"
	"breeze/out"
	"breeze/parser"
	"breeze/scanner"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"
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

	for _, n := range nodes {
		fmt.Println(n.String())
	}

	hadError = analyzer.Analyze(file, source, nodes)

	if hadError {
		out.PrintErrorMessage("Static analyzing phase failed")
		os.Exit(out.ExDataErr)
		return
	}

	executablePath := "test/out"
	compiled := clang.CompileClang(executablePath, file, nodes)
	fmt.Println("-- COMPILED CLANG SOURCE --")
	fmt.Println(compiled)

	err = common.WriteFile("test/breeze.c", compiled)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("-- RUN EXECUTABLE --")

	cmd := exec.Command(fmt.Sprintf("./%s", executablePath))

	var stdout = bytes.Buffer{}
	var stderr = bytes.Buffer{}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	timeStarted := time.Now().UnixMilli()

	exitCode := 0
	if err = cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}
	}

	_, _ = fmt.Println(stderr.String())
	_, _ = fmt.Println(stdout.String())

	timeDelta := time.Now().UnixMilli() - timeStarted

	fmt.Println("Exit Code:", exitCode)
	fmt.Println("Execution took:", timeDelta, "\bms")

}
