package out

import "fmt"

type Color string

var (
	ansiEnabled = true
)

const (
	Black         Color = "\033[30m" // Foreground colors
	Red           Color = "\033[31m"
	Green         Color = "\033[32m"
	Yellow        Color = "\033[33m"
	Blue          Color = "\033[34m"
	Magenta       Color = "\033[35m"
	Cyan          Color = "\033[36m"
	White         Color = "\033[37m"
	Reset         Color = "\033[0m"
	BgBlack       Color = "\033[40m" // Background colors
	BgRed         Color = "\033[41m"
	BgGreen       Color = "\033[42m"
	BgYellow      Color = "\033[43m"
	BgBlue        Color = "\033[44m"
	BgMagenta     Color = "\033[45m"
	BgCyan        Color = "\033[46m"
	BgWhite       Color = "\033[47m"
	Bold          Color = "\033[1m" // Text formatting
	Italic        Color = "\033[3m"
	Underline     Color = "\033[4m"
	StrikeThrough Color = "\033[9m"
	Concealed     Color = "\033[8m"
	Blink         Color = "\033[5m"
)

func SetColorsEnabled(state bool) {
	ansiEnabled = state
}

func printColors(ansi ...Color) {
	if !ansiEnabled {
		return
	}
	for _, a := range ansi {
		fmt.Print(a)
	}
}
