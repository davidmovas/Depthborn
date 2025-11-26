package style

import "github.com/charmbracelet/lipgloss"

// Style represents a lipgloss style
type Style = lipgloss.Style

// Color represents a terminal color
type Color = lipgloss.TerminalColor

// Border represents a lipgloss border style
type Border = lipgloss.Border

var baseStyle = lipgloss.NewStyle()

// Text formatting
var (
	Bold          = baseStyle.Bold(true)
	Italic        = baseStyle.Italic(true)
	Underline     = baseStyle.Underline(true)
	Strikethrough = baseStyle.Strikethrough(true)
	Reverse       = baseStyle.Reverse(true)
	Blink         = baseStyle.Blink(true)
	Faint         = baseStyle.Faint(true)
	Dim           = baseStyle.Faint(true) // Alias for Faint
)

// Layout - Margins
var (
	Margin0 = baseStyle.Margin(0)
	Margin1 = baseStyle.Margin(1)
	Margin2 = baseStyle.Margin(2)
	Margin3 = baseStyle.Margin(3)
	Margin4 = baseStyle.Margin(4)

	// Horizontal margins

	MarginX0 = baseStyle.MarginLeft(0).MarginRight(0)
	MarginX1 = baseStyle.MarginLeft(1).MarginRight(1)
	MarginX2 = baseStyle.MarginLeft(2).MarginRight(2)
	MarginX3 = baseStyle.MarginLeft(3).MarginRight(3)
	MarginX4 = baseStyle.MarginLeft(4).MarginRight(4)

	// Vertical margins

	MarginY0 = baseStyle.MarginTop(0).MarginBottom(0)
	MarginY1 = baseStyle.MarginTop(1).MarginBottom(1)
	MarginY2 = baseStyle.MarginTop(2).MarginBottom(2)
	MarginY3 = baseStyle.MarginTop(3).MarginBottom(3)
	MarginY4 = baseStyle.MarginTop(4).MarginBottom(4)

	// Individual margins

	MarginTop0    = baseStyle.MarginTop(0)
	MarginTop1    = baseStyle.MarginTop(1)
	MarginTop2    = baseStyle.MarginTop(2)
	MarginBottom0 = baseStyle.MarginBottom(0)
	MarginBottom1 = baseStyle.MarginBottom(1)
	MarginBottom2 = baseStyle.MarginBottom(2)
	MarginLeft0   = baseStyle.MarginLeft(0)
	MarginLeft1   = baseStyle.MarginLeft(1)
	MarginLeft2   = baseStyle.MarginLeft(2)
	MarginRight0  = baseStyle.MarginRight(0)
	MarginRight1  = baseStyle.MarginRight(1)
	MarginRight2  = baseStyle.MarginRight(2)
)

// Layout - Padding
var (
	Padding0 = baseStyle.Padding(0)
	Padding1 = baseStyle.Padding(1)
	Padding2 = baseStyle.Padding(2)
	Padding3 = baseStyle.Padding(3)
	Padding4 = baseStyle.Padding(4)

	// Horizontal padding

	PaddingX0 = baseStyle.PaddingLeft(0).PaddingRight(0)
	PaddingX1 = baseStyle.PaddingLeft(1).PaddingRight(1)
	PaddingX2 = baseStyle.PaddingLeft(2).PaddingRight(2)
	PaddingX3 = baseStyle.PaddingLeft(3).PaddingRight(3)
	PaddingX4 = baseStyle.PaddingLeft(4).PaddingRight(4)

	// Vertical padding

	PaddingY0 = baseStyle.PaddingTop(0).PaddingBottom(0)
	PaddingY1 = baseStyle.PaddingTop(1).PaddingBottom(1)
	PaddingY2 = baseStyle.PaddingTop(2).PaddingBottom(2)
	PaddingY3 = baseStyle.PaddingTop(3).PaddingBottom(3)
	PaddingY4 = baseStyle.PaddingTop(4).PaddingBottom(4)

	// Individual padding

	PaddingTop0    = baseStyle.PaddingTop(0)
	PaddingTop1    = baseStyle.PaddingTop(1)
	PaddingTop2    = baseStyle.PaddingTop(2)
	PaddingBottom0 = baseStyle.PaddingBottom(0)
	PaddingBottom1 = baseStyle.PaddingBottom(1)
	PaddingBottom2 = baseStyle.PaddingBottom(2)
	PaddingLeft0   = baseStyle.PaddingLeft(0)
	PaddingLeft1   = baseStyle.PaddingLeft(1)
	PaddingLeft2   = baseStyle.PaddingLeft(2)
	PaddingRight0  = baseStyle.PaddingRight(0)
	PaddingRight1  = baseStyle.PaddingRight(1)
	PaddingRight2  = baseStyle.PaddingRight(2)
)

// Borders
var (
	BorderNone    = baseStyle.Border(lipgloss.HiddenBorder())
	BorderNormal  = baseStyle.Border(lipgloss.NormalBorder())
	BorderRounded = baseStyle.Border(lipgloss.RoundedBorder())
	BorderThick   = baseStyle.Border(lipgloss.ThickBorder())
	BorderDouble  = baseStyle.Border(lipgloss.DoubleBorder())

	// Border sides

	BorderTop    = baseStyle.BorderTop(true)
	BorderRight  = baseStyle.BorderRight(true)
	BorderBottom = baseStyle.BorderBottom(true)
	BorderLeft   = baseStyle.BorderLeft(true)
)

// Alignment
var (
	AlignCenter = baseStyle.Align(lipgloss.Center)
	AlignTop    = baseStyle.Align(lipgloss.Top)
	AlignBottom = baseStyle.Align(lipgloss.Bottom)
	AlignLeft   = baseStyle.Align(lipgloss.Left)
	AlignRight  = baseStyle.Align(lipgloss.Right)
)

// Sizing
var (
	WidthAuto = baseStyle.Width(0)
	Width10   = baseStyle.Width(10)
	Width20   = baseStyle.Width(20)
	Width30   = baseStyle.Width(30)
	Width40   = baseStyle.Width(40)
	Width50   = baseStyle.Width(50)

	HeightAuto = baseStyle.Height(0)
	Height5    = baseStyle.Height(5)
	Height10   = baseStyle.Height(10)
	Height15   = baseStyle.Height(15)
	Height20   = baseStyle.Height(20)
)
