package main

import (
	"breeze/common"
	"breeze/out"
	"breeze/scanner"
	"fmt"
	"os"
)

func main() {
	file := common.InitSource("test/breeze.bz")
	tokens, hadError := scanner.Scan(&file)

	if hadError {
		out.PrintErrorMessage("Scanning phase failed")
		os.Exit(out.ExDataErr)
		return
	}

	for _, tk := range tokens {
		fmt.Println(tk)
	}

}
