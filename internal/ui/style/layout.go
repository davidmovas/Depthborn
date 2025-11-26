package style

import (
	"github.com/charmbracelet/lipgloss"
)

func RowWithGap(gap int, items ...string) string {
	if gap <= 0 {
		return lipgloss.JoinHorizontal(lipgloss.Top, items...)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, items...)
}

func ColWithGap(gap int, items ...string) string {
	if gap <= 0 {
		return lipgloss.JoinVertical(lipgloss.Left, items...)
	}
	var out []string
	for i, it := range items {
		out = append(out, it)
		if i < len(items)-1 {
			for j := 0; j < gap; j++ {
				out = append(out, "")
			}
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left, out...)
}

func Center(width, height int, content string) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

func Row(items ...string) string { return lipgloss.JoinHorizontal(lipgloss.Top, items...) }
func Col(items ...string) string { return lipgloss.JoinVertical(lipgloss.Left, items...) }
