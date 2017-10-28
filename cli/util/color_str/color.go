package color_str

const (
	colorStart  = "\033["
	colorNormal = colorStart + "0m"
	colorRed    = colorStart + "31m"
	colorGreen  = colorStart + "32m"
	colorYellow = colorStart + "33m"
)

func Red(s string) string {
	return colorRed + s + colorNormal
}

func Green(s string) string {
	return colorGreen + s + colorNormal
}

func Yellow(s string) string {
	return colorYellow + s + colorNormal
}
