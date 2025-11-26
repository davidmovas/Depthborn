package style

import "github.com/charmbracelet/lipgloss"

func P(n int) Style  { return baseStyle.Padding(n) }
func Px(n int) Style { return baseStyle.PaddingLeft(n).PaddingRight(n) }
func Py(n int) Style { return baseStyle.PaddingTop(n).PaddingBottom(n) }
func Pt(n int) Style { return baseStyle.PaddingTop(n) }
func Pb(n int) Style { return baseStyle.PaddingBottom(n) }
func Pl(n int) Style { return baseStyle.PaddingLeft(n) }
func Pr(n int) Style { return baseStyle.PaddingRight(n) }

func M(n int) Style  { return baseStyle.Margin(n) }
func Mx(n int) Style { return baseStyle.MarginLeft(n).MarginRight(n) }
func My(n int) Style { return baseStyle.MarginTop(n).MarginBottom(n) }
func Mt(n int) Style { return baseStyle.MarginTop(n) }
func Mb(n int) Style { return baseStyle.MarginBottom(n) }

func Bg(c ...Color) Style {
	out := lipgloss.NewStyle()
	for _, color := range c {
		out = out.Background(color)
	}
	return out
}
func Fg(c ...Color) Style {
	out := lipgloss.NewStyle()
	for _, color := range c {
		out = out.Foreground(color)
	}
	return out
}

func Br() Style { return baseStyle.Border(lipgloss.RoundedBorder()) }

func BrColor(s Style, c Color) Style {
	return s.BorderForeground(c)
}

func Rounded() lipgloss.Style { return baseStyle.Border(lipgloss.RoundedBorder()) }

func W(n int) Style    { return baseStyle.Width(n) }
func H(n int) Style    { return baseStyle.Height(n) }
func MaxW(n int) Style { return baseStyle.MaxWidth(n) }
func MaxH(n int) Style { return baseStyle.MaxHeight(n) }
