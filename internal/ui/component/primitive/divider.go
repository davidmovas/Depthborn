package primitive

import (
	"strings"

	"github.com/davidmovas/Depthborn/internal/ui/component"
)

// Divider creates horizontal divider line
func Divider(width int) component.Component {
	return &divider{
		width: width,
		char:  "─",
	}
}

type divider struct {
	width int
	char  string
}

func (d *divider) Render(ctx *component.Context) string {
	if d.width <= 0 {
		d.width = 40 // default width
	}
	return strings.Repeat(d.char, d.width)
}

// DividerDouble creates double-line divider
func DividerDouble(width int) component.Component {
	return &divider{
		width: width,
		char:  "═",
	}
}

// DividerThick creates thick divider
func DividerThick(width int) component.Component {
	return &divider{
		width: width,
		char:  "━",
	}
}

// DividerDotted creates dotted divider
func DividerDotted(width int) component.Component {
	return &divider{
		width: width,
		char:  "·",
	}
}

// DividerCustom creates divider with custom character
func DividerCustom(width int, char string) component.Component {
	return &divider{
		width: width,
		char:  char,
	}
}

// DividerWithText creates divider with text in the middle
// Example: ───── Title ─────
func DividerWithText(width int, text string) component.Component {
	return &dividerWithText{
		width: width,
		text:  text,
	}
}

type dividerWithText struct {
	width int
	text  string
}

func (d *dividerWithText) Render(ctx *component.Context) string {
	if d.width <= 0 {
		d.width = 40
	}

	textLen := len(d.text)
	if textLen >= d.width-4 {
		// Not enough space, just return text
		return d.text
	}

	// Calculate padding
	totalPad := d.width - textLen - 2 // -2 for spaces around text
	leftPad := totalPad / 2
	rightPad := totalPad - leftPad

	left := strings.Repeat("─", leftPad)
	right := strings.Repeat("─", rightPad)

	return left + " " + d.text + " " + right
}
