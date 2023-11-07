package out

type Color string

var (
	ansiEnabled = true
)

//goland:noinspection ALL
const (
	ColorBlack         Color = "\033[30m" // Foreground colors
	ColorRed           Color = "\033[31m"
	ColorGreen         Color = "\033[32m"
	ColorYellow        Color = "\033[33m"
	ColorBlue          Color = "\033[34m"
	ColorMagenta       Color = "\033[35m"
	ColorCyan          Color = "\033[36m"
	ColorWhite         Color = "\033[37m"
	ColorReset         Color = "\033[0m"
	ColorBgBlack       Color = "\033[40m" // Background colors
	ColorBgRed         Color = "\033[41m"
	ColorBgGreen       Color = "\033[42m"
	ColorBgYellow      Color = "\033[43m"
	ColorBgBlue        Color = "\033[44m"
	ColorBgMagenta     Color = "\033[45m"
	ColorBgCyan        Color = "\033[46m"
	ColorBgWhite       Color = "\033[47m"
	ColorBold          Color = "\033[1m" // Text formatting
	ColorItalic        Color = "\033[3m"
	ColorUnderline     Color = "\033[4m"
	ColorStrikeThrough Color = "\033[9m"
	ColorConcealed     Color = "\033[8m"
	ColorBlink         Color = "\033[5m"
)

func (c Color) S() string {
	if ansiEnabled {
		return string(c)
	}
	return ""
}

func SetColorsEnabled(state bool) {
	ansiEnabled = state
}
