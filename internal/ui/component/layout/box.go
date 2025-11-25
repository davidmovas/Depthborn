package layout

import (
	"strings"

	"github.com/davidmovas/Depthborn/internal/ui/component"
)

// Box creates a container with padding/margin
// Useful for adding spacing around components
//
// Example:
//
//	Box(BoxProps{Padding: 2}, Text("Content"))
func Box(props BoxProps, child component.Component) component.Component {
	return &box{
		props: props,
		child: child,
	}
}

// BoxProps configures box padding and borders
type BoxProps struct {
	Padding int    // Inner padding (spaces)
	Border  bool   // Draw border
	Title   string // Optional title (shown in border)
}

type box struct {
	props BoxProps
	child component.Component
}

func (b *box) Render(ctx *component.Context) string {
	if b.child == nil {
		return ""
	}

	content := b.child.Render(ctx)
	if content == "" {
		return ""
	}

	// Apply padding
	if b.props.Padding > 0 {
		content = b.applyPadding(content, b.props.Padding)
	}

	// Apply border
	if b.props.Border {
		content = b.applyBorder(content, b.props.Title)
	}

	return content
}

func (b *box) applyPadding(content string, padding int) string {
	lines := strings.Split(content, "\n")
	paddingStr := strings.Repeat(" ", padding)

	var padded []string

	// Top padding
	for i := 0; i < padding; i++ {
		padded = append(padded, "")
	}

	// Content with side padding
	for _, line := range lines {
		padded = append(padded, paddingStr+line+paddingStr)
	}

	// Bottom padding
	for i := 0; i < padding; i++ {
		padded = append(padded, "")
	}

	return strings.Join(padded, "\n")
}

func (b *box) applyBorder(content string, title string) string {
	lines := strings.Split(content, "\n")

	// Calculate width (longest line)
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	var result []string

	// Top border
	if title != "" {
		// ╔══ Title ══╗
		titleLen := len(title)
		leftPad := (maxWidth - titleLen - 2) / 2
		rightPad := maxWidth - titleLen - 2 - leftPad

		top := "╔" + strings.Repeat("═", leftPad) + " " + title + " " + strings.Repeat("═", rightPad) + "╗"
		result = append(result, top)
	} else {
		// ╔════════╗
		top := "╔" + strings.Repeat("═", maxWidth) + "╗"
		result = append(result, top)
	}

	// Content with side borders
	for _, line := range lines {
		// Pad line to max width
		padded := line + strings.Repeat(" ", maxWidth-len(line))
		result = append(result, "║"+padded+"║")
	}

	// Bottom border
	bottom := "╚" + strings.Repeat("═", maxWidth) + "╝"
	result = append(result, bottom)

	return strings.Join(result, "\n")
}

// Spacer creates empty space (for vertical spacing)
func Spacer(lines int) component.Component {
	return component.Func(func(ctx *component.Context) string {
		return strings.Repeat("\n", lines)
	})
}

// HSpacer creates horizontal space
func HSpacer(spaces int) component.Component {
	return component.Func(func(ctx *component.Context) string {
		return strings.Repeat(" ", spaces)
	})
}
