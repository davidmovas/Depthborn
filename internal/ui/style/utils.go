package style

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// S is an alias for returning styles variadically — convenience
func S(styles ...Style) *Style {
	s := Merge(styles...)
	return &s
}

// Merge combines multiple lipgloss styles into one via Inherit
func Merge(styles ...Style) Style {
	out := lipgloss.NewStyle()
	for _, s := range styles {
		out = out.Inherit(s)
	}
	return out
}

// If returns style if cond true, otherwise empty style
func If(cond bool, s Style) Style {
	if cond {
		return s
	}
	return lipgloss.NewStyle()
}

// Render applies merged style to content
func Render(content string, styles ...Style) string {
	return Merge(styles...).Render(content)
}

// StyledArg for Sprintf-like helper
type StyledArg struct {
	Val    any
	Styles []Style
}

func Val(val any, styles ...Style) StyledArg {
	return StyledArg{Val: val, Styles: styles}
}

// Sprintf renders styled args and formats into a string
func Sprintf(format string, args ...StyledArg) string {
	vals := make([]any, 0, len(args))
	for _, a := range args {
		vals = append(vals, Merge(a.Styles...).Render(fmt.Sprintf("%v", a.Val)))
	}
	return fmt.Sprintf(format, vals...)
}
