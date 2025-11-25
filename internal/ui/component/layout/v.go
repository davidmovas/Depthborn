package layout

import (
	"strings"

	"github.com/davidmovas/Depthborn/internal/ui/component"
)

// V creates vertical stack component (children stacked vertically)
// Short for "Vertical" - similar to SwiftUI's VStack
//
// Example:
//
//	V(
//	    Text("Header"),
//	    Text("Content"),
//	    Text("Footer"),
//	)
func V(children ...component.Component) component.Component {
	return &vStack{children: children}
}

type vStack struct {
	children []component.Component
}

func (v *vStack) Render(ctx *component.Context) string {
	if len(v.children) == 0 {
		return ""
	}

	var parts []string
	for _, child := range v.children {
		if child == nil {
			continue
		}
		rendered := child.Render(ctx)
		if rendered != "" {
			parts = append(parts, rendered)
		}
	}

	// Join with newlines for vertical stacking
	return strings.Join(parts, "\n")
}

// VSpaced creates vertical stack with custom spacing
// spacing: number of newlines between children
func VSpaced(spacing int, children ...component.Component) component.Component {
	return &vStackSpaced{
		children: children,
		spacing:  spacing,
	}
}

type vStackSpaced struct {
	children []component.Component
	spacing  int
}

func (v *vStackSpaced) Render(ctx *component.Context) string {
	if len(v.children) == 0 {
		return ""
	}

	var parts []string
	for _, child := range v.children {
		if child == nil {
			continue
		}
		rendered := child.Render(ctx)
		if rendered != "" {
			parts = append(parts, rendered)
		}
	}

	// Create spacing string
	spacer := strings.Repeat("\n", v.spacing)

	return strings.Join(parts, spacer)
}
