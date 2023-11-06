package out

import (
	"breeze/common"
	"fmt"
	"io"
	"os"
)

func PrintErrorMessage(message string) {
	_, err := fmt.Fprintf(os.Stderr, "%s%s%sERROR%s   %s%s%s%s\n", ColorBgRed.S(), ColorBlack.S(), ColorBold.S(), ColorReset.S(), ColorRed.S(), ColorBold.S(), message, ColorReset.S())
	if err != nil {
		os.Exit(ExIoErr)
		return
	}
}

func PrintErrorSource(path string, position common.Position) {
	// →
	_, err := fmt.Fprintf(os.Stderr, "%s      → %s:%d:%d%s\n", ColorWhite.S(), path, position.Line, position.Column, ColorReset.S())
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
	markedLine := markLexeme(lineString, length, position.Column, color, ColorBold)
	printLineString(writer, position.Line, markedLine)

	marker := getMarker(length, position.Column, icon)
	_, err := fmt.Fprintf(writer, "      | %s%s%s\n", color, marker, ColorReset.S())
	if err != nil {
		os.Exit(ExIoErr)
		return
	}
}
