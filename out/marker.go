package out

import (
	"fmt"
	"os"
)

func getMarker(length int, column int, icon rune) string {
	marker := ""
	for i := 0; i < column-1; i++ {
		marker += " "
	}
	for i := 0; i < length; i++ {
		marker += string(icon)
	}

	return marker
}

func markLexeme(line string, length int, column int, colors ...Color) string {
	runes := []rune(line)
	lineLen := len(runes)

	if column < 0 || length < 0 || column+length > lineLen+1 {
		fmt.Println("Lexeme marker out of bounds")
		os.Exit(70)
	}

	beforeEnd := column
	if beforeEnd > 0 {
		beforeEnd--
	}

	lexemeEnd := beforeEnd + length
	if lexemeEnd >= lineLen+1 {
		lexemeEnd = lineLen - 1
	}

	before := runes[0:beforeEnd]
	lexeme := runes[beforeEnd:lexemeEnd]
	after := runes[lexemeEnd:lineLen]

	colorString := ""
	for _, c := range colors {
		colorString += string(c)
	}

	return string(before) + colorString + string(lexeme) + string(ColorReset) + string(after)
}

func getLineBounds(source string, index int) (int, int) {
	sourceLen := len(source)

	lineStart := index
	lineEnd := index

	for {
		if lineStart <= 0 || lineStart >= sourceLen {
			break
		}

		if source[lineStart] == '\n' {
			lineStart++ // Offset newline char
			break
		}

		lineStart--
	}

	for {
		if lineEnd < 0 || lineEnd >= sourceLen {
			break
		}

		if source[lineEnd] == '\n' {
			break
		}

		lineEnd++
	}

	return lineStart, lineEnd
}

func getSourceLine(source string, index int) string {
	lineStart, lineEnd := getLineBounds(source, index)
	lineString := source[lineStart:lineEnd]
	return lineString
}
