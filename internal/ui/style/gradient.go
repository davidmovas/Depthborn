package style

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

// GradientBuilder builds gradients from colors
type GradientBuilder struct {
	colors    []Color
	steps     int
	direction GradientDirection
}

// GradientDirection defines the direction of a gradient
type GradientDirection int

const (
	DirectionLeftToRight GradientDirection = iota
	DirectionRightToLeft
	DirectionTopToBottom
	DirectionBottomToTop
	DirectionDiagonalTLBR
	DirectionDiagonalTRBL
	DirectionDiagonalBLTR
	DirectionDiagonalBRTL
	DirectionBorderClockwise
	DirectionBorderCounterClockwise
	DirectionBorderVertical
	DirectionBorderHorizontal
)

// BoxBorderStyle defines the characters used for drawing a box border
type BoxBorderStyle struct {
	TopLeft     string
	Top         string
	TopRight    string
	Right       string
	BottomRight string
	Bottom      string
	BottomLeft  string
	Left        string
}

var (
	// BoxBorderRounded uses rounded corners: ╭─╮│╰─╯│
	BoxBorderRounded = BoxBorderStyle{
		TopLeft:     "╭",
		Top:         "─",
		TopRight:    "╮",
		Right:       "│",
		BottomRight: "╯",
		Bottom:      "─",
		BottomLeft:  "╰",
		Left:        "│",
	}

	// BoxBorderNormal uses normal corners: ┌─┐│└─┘│
	BoxBorderNormal = BoxBorderStyle{
		TopLeft:     "┌",
		Top:         "─",
		TopRight:    "┐",
		Right:       "│",
		BottomRight: "┘",
		Bottom:      "─",
		BottomLeft:  "└",
		Left:        "│",
	}

	// BoxBorderThick uses thick lines: ┏━┓┃┗━┛┃
	BoxBorderThick = BoxBorderStyle{
		TopLeft:     "┏",
		Top:         "━",
		TopRight:    "┓",
		Right:       "┃",
		BottomRight: "┛",
		Bottom:      "━",
		BottomLeft:  "┗",
		Left:        "┃",
	}

	// BoxBorderDouble uses double lines: ╔═╗║╚═╝║
	BoxBorderDouble = BoxBorderStyle{
		TopLeft:     "╔",
		Top:         "═",
		TopRight:    "╗",
		Right:       "║",
		BottomRight: "╝",
		Bottom:      "═",
		BottomLeft:  "╚",
		Left:        "║",
	}

	// BoxBorderHidden uses spaces (invisible border)
	BoxBorderHidden = BoxBorderStyle{
		TopLeft:     " ",
		Top:         " ",
		TopRight:    " ",
		Right:       " ",
		BottomRight: " ",
		Bottom:      " ",
		BottomLeft:  " ",
		Left:        " ",
	}
)

// NewGradient creates a new gradient builder from colors
func NewGradient(colors ...Color) *GradientBuilder {
	return &GradientBuilder{
		colors:    colors,
		direction: DirectionLeftToRight,
	}
}

// Direction sets the gradient direction
func (g *GradientBuilder) Direction(dir GradientDirection) *GradientBuilder {
	g.direction = dir
	return g
}

// Steps sets the number of gradient steps
func (g *GradientBuilder) Steps(steps int) *GradientBuilder {
	g.steps = steps
	return g
}

// Foreground applies gradient to text (foreground color)
// Skips already styled parts
func (g *GradientBuilder) Foreground(text string) string {
	result := g.applySmartGradient(text, "foreground")
	return g.applyTextDirection(result, text)
}

// ForegroundRaw applies gradient to all text without considering existing styles
func (g *GradientBuilder) ForegroundRaw(text string) string {
	runes := []rune(text)
	if len(runes) == 0 {
		return ""
	}

	gradient := g.generate(len(runes))
	var result string

	for i, char := range runes {
		style := lipgloss.NewStyle().Foreground(gradient[i])
		result += style.Render(string(char))
	}

	return result
}

// Background applies gradient as background
// Skips already styled parts
func (g *GradientBuilder) Background(text string) string {
	return g.applySmartGradient(text, "background")
}

// BackgroundRaw applies gradient as background to all text
func (g *GradientBuilder) BackgroundRaw(text string) string {
	runes := []rune(text)
	if len(runes) == 0 {
		return ""
	}

	gradient := g.generate(len(runes))
	var result string

	for i, char := range runes {
		style := lipgloss.NewStyle().Background(gradient[i])
		result += style.Render(string(char))
	}

	return result
}

// Style applies gradient to a base style
// baseStyle - your pre-configured style (with border, padding, etc.)
// applyTo - where to apply: "foreground"/"fg", "background"/"bg", "border"
// Note: For borders, only applies first color (lipgloss limitation)
func (g *GradientBuilder) Style(baseStyle Style, applyTo string) Style {
	gradient := g.generate(1)
	if len(gradient) == 0 {
		return baseStyle
	}

	result := baseStyle

	switch applyTo {
	case "foreground", "fg", "color":
		result = result.Foreground(gradient[0])
	case "background", "bg":
		result = result.Background(gradient[0])
	case "border":
		result = result.BorderForeground(gradient[0])
	}

	return result
}

// Styles returns an array of styles with gradient applied
// baseStyle - your pre-configured style
// steps - number of variations
// applyTo - where to apply gradient
func (g *GradientBuilder) Styles(baseStyle Style, steps int, applyTo string) []Style {
	gradient := g.generate(steps)
	styles := make([]Style, steps)

	for i, color := range gradient {
		result := baseStyle

		switch applyTo {
		case "foreground", "fg", "color":
			result = result.Foreground(color)
		case "background", "bg":
			result = result.Background(color)
		case "border":
			result = result.BorderForeground(color)
		}

		styles[i] = result
	}

	return styles
}

// BorderGradient applies gradient to each side of border
// Note: Each side will be a single color (lipgloss limitation)
// For smooth gradients, use BorderGradientBox
func (g *GradientBuilder) BorderGradient(baseStyle Style) Style {
	var topColor, rightColor, bottomColor, leftColor Color

	switch g.direction {
	case DirectionBorderClockwise:
		gradient := g.generate(4)
		if len(gradient) >= 4 {
			topColor = gradient[0]
			rightColor = gradient[1]
			bottomColor = gradient[2]
			leftColor = gradient[3]
		}

	case DirectionBorderCounterClockwise:
		gradient := g.generate(4)
		if len(gradient) >= 4 {
			topColor = gradient[0]
			leftColor = gradient[1]
			bottomColor = gradient[2]
			rightColor = gradient[3]
		}

	case DirectionBorderVertical:
		gradient := g.generate(4)
		if len(gradient) >= 4 {
			topColor = gradient[0]
			bottomColor = gradient[3]
			leftColor = gradient[1]
			rightColor = gradient[2]
		}

	case DirectionBorderHorizontal:
		gradient := g.generate(4)
		if len(gradient) >= 4 {
			leftColor = gradient[0]
			rightColor = gradient[3]
			topColor = gradient[1]
			bottomColor = gradient[2]
		}

	case DirectionDiagonalTLBR:
		gradient := g.generate(4)
		if len(gradient) >= 4 {
			topColor = gradient[0]
			leftColor = gradient[1]
			rightColor = gradient[2]
			bottomColor = gradient[3]
		}

	case DirectionDiagonalBLTR:
		gradient := g.generate(4)
		if len(gradient) >= 4 {
			bottomColor = gradient[0]
			leftColor = gradient[1]
			topColor = gradient[2]
			rightColor = gradient[3]
		}

	case DirectionDiagonalTRBL:
		gradient := g.generate(4)
		if len(gradient) >= 4 {
			topColor = gradient[0]
			rightColor = gradient[1]
			leftColor = gradient[2]
			bottomColor = gradient[3]
		}

	case DirectionDiagonalBRTL:
		gradient := g.generate(4)
		if len(gradient) >= 4 {
			bottomColor = gradient[0]
			rightColor = gradient[1]
			topColor = gradient[2]
			leftColor = gradient[3]
		}

	default:
		gradient := g.generate(4)
		if len(gradient) >= 4 {
			topColor = gradient[0]
			rightColor = gradient[1]
			bottomColor = gradient[2]
			leftColor = gradient[3]
		}
	}

	result := baseStyle

	result = result.
		BorderTopForeground(topColor).
		BorderRightForeground(rightColor).
		BorderBottomForeground(bottomColor).
		BorderLeftForeground(leftColor)

	return result
}

// BorderGradientBox creates a box with smooth gradient border
// Content width is calculated automatically based on the longest line
func (g *GradientBuilder) BorderGradientBox(content string, borderStyle ...BoxBorderStyle) string {
	border := BoxBorderRounded
	if len(borderStyle) > 0 {
		border = borderStyle[0]
	}

	lines := strings.Split(content, "\n")
	height := len(lines)

	// Calculate width using clean text (without ANSI) with proper rune width calculation
	width := 0
	for _, line := range lines {
		cleanLine := stripAnsi(line)
		lineWidth := runewidth.StringWidth(cleanLine)
		if lineWidth > width {
			width = lineWidth
		}
	}

	// Calculate gradient for full perimeter
	// Top: 1 (corner) + width + 1 (corner)
	// Bottom: 1 (corner) + width + 1 (corner)
	// Left/Right: height * 2 (both sides)
	perimeterLength := (width+2)*2 + height*2
	gradient := g.generate(perimeterLength)

	if len(gradient) == 0 {
		gradient = []Color{Grey400}
	}

	var result strings.Builder
	gradIndex := 0

	getColor := func() Color {
		if gradIndex >= len(gradient) {
			return gradient[len(gradient)-1]
		}
		color := gradient[gradIndex]
		gradIndex++
		return color
	}

	// Top border: corner + horizontal + corner
	result.WriteString(lipgloss.NewStyle().Foreground(getColor()).Render(border.TopLeft))
	for i := 0; i < width; i++ {
		result.WriteString(lipgloss.NewStyle().Foreground(getColor()).Render(border.Top))
	}
	result.WriteString(lipgloss.NewStyle().Foreground(getColor()).Render(border.TopRight))
	result.WriteString("\n")

	// Content lines: left + content + right
	for i := 0; i < height; i++ {
		result.WriteString(lipgloss.NewStyle().Foreground(getColor()).Render(border.Left))

		if i < len(lines) {
			line := lines[i]
			cleanLine := stripAnsi(line)
			lineWidth := runewidth.StringWidth(cleanLine)

			// Write the original line (with styles)
			result.WriteString(line)

			// Pad with spaces to reach target width
			if lineWidth < width {
				padding := width - lineWidth
				result.WriteString(strings.Repeat(" ", padding))
			}
		} else {
			result.WriteString(strings.Repeat(" ", width))
		}

		result.WriteString(lipgloss.NewStyle().Foreground(getColor()).Render(border.Right))
		result.WriteString("\n")
	}

	// Bottom border: corner + horizontal + corner
	result.WriteString(lipgloss.NewStyle().Foreground(getColor()).Render(border.BottomLeft))
	for i := 0; i < width; i++ {
		result.WriteString(lipgloss.NewStyle().Foreground(getColor()).Render(border.Bottom))
	}
	result.WriteString(lipgloss.NewStyle().Foreground(getColor()).Render(border.BottomRight))

	return result.String()
}

// Colors return an array of interpolated colors
func (g *GradientBuilder) Colors(steps int) []Color {
	return g.generate(steps)
}

// ForegroundStyle returns an array of styles with gradient for foreground
func (g *GradientBuilder) ForegroundStyle(steps int) []Style {
	gradient := g.generate(steps)
	styles := make([]Style, steps)

	for i, color := range gradient {
		styles[i] = lipgloss.NewStyle().Foreground(color)
	}

	return styles
}

// BackgroundStyle returns an array of styles with gradient for background
func (g *GradientBuilder) BackgroundStyle(steps int) []Style {
	gradient := g.generate(steps)
	styles := make([]Style, steps)

	for i, color := range gradient {
		styles[i] = lipgloss.NewStyle().Background(color)
	}

	return styles
}

// GradientFg applies gradient to text (foreground), skipping styled parts
func GradientFg(text string, colors ...Color) string {
	return NewGradient(colors...).Foreground(text)
}

// GradientBg applies gradient as background, skipping styled parts
func GradientBg(text string, colors ...Color) string {
	return NewGradient(colors...).Background(text)
}

// GradientFgRaw applies gradient to all text
func GradientFgRaw(text string, colors ...Color) string {
	return NewGradient(colors...).ForegroundRaw(text)
}

// GradientBgRaw applies gradient as background to all text
func GradientBgRaw(text string, colors ...Color) string {
	return NewGradient(colors...).BackgroundRaw(text)
}

// GradientBorder creates a style with gradient border from base style
// Note: Applies first gradient color to entire border (lipgloss limitation)
func GradientBorder(baseStyle Style, colors ...Color) Style {
	return NewGradient(colors...).Style(baseStyle, "border")
}

// GradientBorderFull creates a style with full gradient on border
// Note: Each side will be single color (lipgloss limitation)
// For smooth gradients, use GradientBorderBox
func GradientBorderFull(baseStyle Style, colors ...Color) Style {
	return NewGradient(colors...).BorderGradient(baseStyle)
}

// GradientBorderBox creates a box with smooth gradient border
// Border style is optional, defaults to BoxBorderRounded
func GradientBorderBox(content string, colors []Color, borderStyle ...BoxBorderStyle) string {
	return NewGradient(colors...).BorderGradientBox(content, borderStyle...)
}

// GradientStyles returns an array of styles with gradient
func GradientStyles(baseStyle Style, steps int, applyTo string, colors ...Color) []Style {
	return NewGradient(colors...).Styles(baseStyle, steps, applyTo)
}

// stripAnsi removes all ANSI escape codes from string
func stripAnsi(s string) string {
	ansiPattern := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiPattern.ReplaceAllString(s, "")
}

func (g *GradientBuilder) applyTextDirection(renderedText, originalText string) string {
	switch g.direction {
	case DirectionRightToLeft:
		return g.reverseGradient(originalText, "foreground")
	default:
		return renderedText
	}
}

func (g *GradientBuilder) reverseGradient(text string, applyTo string) string {
	runes := []rune(text)
	if len(runes) == 0 {
		return ""
	}

	gradient := g.generate(len(runes))
	for i, j := 0, len(gradient)-1; i < j; i, j = i+1, j-1 {
		gradient[i], gradient[j] = gradient[j], gradient[i]
	}

	var result string
	for i, char := range runes {
		var style lipgloss.Style
		switch applyTo {
		case "foreground":
			style = lipgloss.NewStyle().Foreground(gradient[i])
		case "background":
			style = lipgloss.NewStyle().Background(gradient[i])
		}
		result += style.Render(string(char))
	}

	return result
}

func (g *GradientBuilder) applySmartGradient(text string, applyTo string) string {
	if len(text) == 0 {
		return ""
	}

	ansiPattern := regexp.MustCompile(`\x1b\[[0-9;]*m`)

	var parts []struct {
		text   string
		isAnsi bool
	}

	lastIndex := 0
	matches := ansiPattern.FindAllStringIndex(text, -1)

	for _, match := range matches {
		if match[0] > lastIndex {
			parts = append(parts, struct {
				text   string
				isAnsi bool
			}{text[lastIndex:match[0]], false})
		}
		parts = append(parts, struct {
			text   string
			isAnsi bool
		}{text[match[0]:match[1]], true})
		lastIndex = match[1]
	}

	if lastIndex < len(text) {
		parts = append(parts, struct {
			text   string
			isAnsi bool
		}{text[lastIndex:], false})
	}

	if len(parts) == 0 {
		return text
	}

	visibleChars := 0
	for _, part := range parts {
		if !part.isAnsi {
			visibleChars += len([]rune(part.text))
		}
	}

	if visibleChars == 0 {
		return text
	}

	gradient := g.generate(visibleChars)

	var result strings.Builder
	gradientIndex := 0

	for _, part := range parts {
		if part.isAnsi {
			result.WriteString(part.text)
		} else {
			for _, char := range part.text {
				if gradientIndex >= len(gradient) {
					break
				}

				var style lipgloss.Style
				switch applyTo {
				case "foreground":
					style = lipgloss.NewStyle().Foreground(gradient[gradientIndex])
				case "background":
					style = lipgloss.NewStyle().Background(gradient[gradientIndex])
				}

				result.WriteString(style.Render(string(char)))
				gradientIndex++
			}
		}
	}

	return result.String()
}

func (g *GradientBuilder) generate(steps int) []Color {
	if len(g.colors) == 0 {
		return []Color{}
	}
	if len(g.colors) == 1 {
		result := make([]Color, steps)
		for i := range result {
			result[i] = g.colors[0]
		}
		return result
	}
	if steps <= 0 {
		return []Color{}
	}
	if steps == 1 {
		return []Color{g.colors[0]}
	}

	gradient := make([]Color, steps)

	for i := 0; i < steps; i++ {
		position := float64(i) / float64(steps-1)

		segmentSize := 1.0 / float64(len(g.colors)-1)

		segmentIndex := int(position / segmentSize)

		if segmentIndex < 0 {
			segmentIndex = 0
		}
		if segmentIndex >= len(g.colors)-1 {
			segmentIndex = len(g.colors) - 2
		}

		localPosition := (position - float64(segmentIndex)*segmentSize) / segmentSize

		if localPosition < 0 {
			localPosition = 0
		}
		if localPosition > 1 {
			localPosition = 1
		}

		color1 := g.colors[segmentIndex]
		color2 := g.colors[segmentIndex+1]

		gradient[i] = interpolateColor(color1, color2, localPosition)
	}

	return gradient
}

func interpolateColor(color1, color2 Color, t float64) Color {
	r1, g1, b1, _ := color1.RGBA()
	r2, g2, b2, _ := color2.RGBA()

	r18 := uint8(r1 >> 8)
	g18 := uint8(g1 >> 8)
	b18 := uint8(b1 >> 8)

	r28 := uint8(r2 >> 8)
	g28 := uint8(g2 >> 8)
	b28 := uint8(b2 >> 8)

	r := uint8(float64(r18) + (float64(r28)-float64(r18))*t)
	g := uint8(float64(g18) + (float64(g28)-float64(g18))*t)
	b := uint8(float64(b18) + (float64(b28)-float64(b18))*t)

	return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, b))
}
