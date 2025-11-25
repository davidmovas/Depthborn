package primitive

import (
	"fmt"
	"strings"

	"github.com/davidmovas/Depthborn/internal/ui/component"
)

// Text creates simple text component
func Text(text string) component.Component {
	return &textComponent{text: text}
}

type textComponent struct {
	text string
}

func (t *textComponent) Render(ctx *component.Context) string {
	return t.text
}

// Textf creates formatted text component (like fmt.Sprintf)
func Textf(format string, args ...any) component.Component {
	return Text(fmt.Sprintf(format, args...))
}

type TextStyle struct {
	Color      string // "red", "green", "blue", etc.
	Bold       bool
	Underline  bool
	Background string // background color
}

// TextStyled creates styled text with color/formatting
// Uses ANSI escape codes for terminal
func TextStyled(text string, style TextStyle) component.Component {
	return &styledText{
		text:  text,
		style: style,
	}
}

type styledText struct {
	text  string
	style TextStyle
}

func (st *styledText) Render(ctx *component.Context) string {
	// Apply ANSI styles
	result := st.text

	// Colors (basic ANSI)
	colorCodes := map[string]string{
		"black":   "\033[30m",
		"red":     "\033[31m",
		"green":   "\033[32m",
		"yellow":  "\033[33m",
		"blue":    "\033[34m",
		"magenta": "\033[35m",
		"cyan":    "\033[36m",
		"white":   "\033[37m",
		"gray":    "\033[90m",
	}

	bgColorCodes := map[string]string{
		"black":   "\033[40m",
		"red":     "\033[41m",
		"green":   "\033[42m",
		"yellow":  "\033[43m",
		"blue":    "\033[44m",
		"magenta": "\033[45m",
		"cyan":    "\033[46m",
		"white":   "\033[47m",
	}

	var codes []string

	if st.style.Bold {
		codes = append(codes, "\033[1m")
	}

	if st.style.Underline {
		codes = append(codes, "\033[4m")
	}

	if colorCode, ok := colorCodes[st.style.Color]; ok {
		codes = append(codes, colorCode)
	}

	if bgCode, ok := bgColorCodes[st.style.Background]; ok {
		codes = append(codes, bgCode)
	}

	if len(codes) > 0 {
		result = strings.Join(codes, "") + result + "\033[0m" // reset
	}

	return result
}

// Label creates label with optional styling
func Label(label string, value string) component.Component {
	return Textf("%s: %s", label, value)
}

// Title creates title text (bold, colored)
func Title(text string) component.Component {
	return TextStyled(text, TextStyle{
		Bold:  true,
		Color: "cyan",
	})
}

// Error creates error text (red)
func Error(text string) component.Component {
	return TextStyled(text, TextStyle{
		Color: "red",
		Bold:  true,
	})
}

// Success creates success text (green)
func Success(text string) component.Component {
	return TextStyled(text, TextStyle{
		Color: "green",
	})
}

// Warning creates warning text (yellow)
func Warning(text string) component.Component {
	return TextStyled(text, TextStyle{
		Color: "yellow",
	})
}

// Muted creates muted text (gray)
func Muted(text string) component.Component {
	return TextStyled(text, TextStyle{
		Color: "gray",
	})
}
