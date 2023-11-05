package main

import (
	"breeze/out"
	"breeze/scanner"
	"fmt"
	"os"
)

func main() {
	tokens, hadError := scanner.Scan("#(( ) identifier Another1 .4 if 10 1.5 ä \"hello\"\nhello ö")

	if hadError {
		out.PrintErrorMessage("Scanning phase failed")
		os.Exit(out.ExDataErr)
		return
	}

	for _, tk := range tokens {
		fmt.Println(tk)
	}

}
