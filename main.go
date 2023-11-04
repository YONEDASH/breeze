package main

import (
	"breeze/scanner"
	"fmt"
)

func main() {
	tokens, hadError := scanner.Scan("(( ) identifier Another1 .4 if 10 1.5 \"hello\"")

	for _, tk := range tokens {
		fmt.Println(tk)
	}

	fmt.Println(hadError)
}
