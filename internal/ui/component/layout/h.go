package layout

import (
	"strings"

	"github.com/davidmovas/Depthborn/internal/ui/component"
)

// H creates horizontal stack component (children stacked horizontally)
// Short for "Horizontal" - similar to SwiftUI's HStack
//
// Example:
//
//	H(
//	    Text("Left"),
//	    Text("Center"),
//	    Text("Right"),
//	)
func H(children ...component.Component) component.Component {
	return &hStack{children: children}
}

type hStack struct {
	children []component.Component
}

func (h *hStack) Render(ctx *component.Context) string {
	if len(h.children) == 0 {
		return ""
	}

	var parts []string
	for _, child := range h.children {
		if child == nil {
			continue
		}
		rendered := child.Render(ctx)
		if rendered != "" {
			parts = append(parts, rendered)
		}
	}

	// Join with spaces for horizontal alignment
	return strings.Join(parts, " ")
}

// HSpaced creates horizontal stack with custom spacing
// spacing: number of spaces between children
func HSpaced(spacing int, children ...component.Component) component.Component {
	return &hStackSpaced{
		children: children,
		spacing:  spacing,
	}
}

type hStackSpaced struct {
	children []component.Component
	spacing  int
}

func (h *hStackSpaced) Render(ctx *component.Context) string {
	if len(h.children) == 0 {
		return ""
	}

	var parts []string
	for _, child := range h.children {
		if child == nil {
			continue
		}
		rendered := child.Render(ctx)
		if rendered != "" {
			parts = append(parts, rendered)
		}
	}

	// Create spacing string
	spacer := strings.Repeat(" ", h.spacing)

	return strings.Join(parts, spacer)
}

// HJoin creates horizontal stack with custom separator
// Useful for creating lists like "A | B | C"
func HJoin(separator string, children ...component.Component) component.Component {
	return &hStackJoin{
		children:  children,
		separator: separator,
	}
}

type hStackJoin struct {
	children  []component.Component
	separator string
}

func (h *hStackJoin) Render(ctx *component.Context) string {
	if len(h.children) == 0 {
		return ""
	}

	var parts []string
	for _, child := range h.children {
		if child == nil {
			continue
		}
		rendered := child.Render(ctx)
		if rendered != "" {
			parts = append(parts, rendered)
		}
	}

	return strings.Join(parts, h.separator)
}
