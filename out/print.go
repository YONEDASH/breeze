package out

import (
	"breeze/common"
	"fmt"
	"io"
	"os"
)

func PrintErrorMessage(message string) {
	_, err := fmt.Fprintf(os.Stderr, "%s%s%sERROR%s   %s%s%s%s\n", BgRed, Black, Bold, Reset, Red, Bold, message, Reset)
	if err != nil {
		os.Exit(ExIoErr)
		return
	}
}

func PrintErrorSource(path string, position common.Position) {
	// →
	_, err := fmt.Fprintf(os.Stderr, "%s      → %s:%d:%d%s\n", White, path, position.Line, position.Column, Reset)
	if err != nil {
		os.Exit(ExIoErr)
		return
	}
}

func printLineString(writer io.Writer, line int, lineString string) {
	_, err := fmt.Fprintf(writer, "%5d | %s\n", line, lineString)
	if err != nil {
		os.Exit(ExIoErr)
		return
	}
}

func PrintLine(writer io.Writer, source string, position common.Position) {
	lineString := getSourceLine(source, position.Index)
	printLineString(writer, position.Line, lineString)
}

func PrintMarkedLine(writer io.Writer, source string, length int, position common.Position, color Color, icon rune) {
	lineString := getSourceLine(source, position.Index)
	markedLine := markLexeme(lineString, length, position.Column, color, Bold)
	printLineString(writer, position.Line, markedLine)

	marker := getMarker(length, position.Column, icon)
	_, err := fmt.Fprintf(writer, "      | %s%s%s\n", color, marker, Reset)
	if err != nil {
		os.Exit(ExIoErr)
		return
	}
}
