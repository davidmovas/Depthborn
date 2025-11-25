package layout

import (
	"strings"

	"github.com/davidmovas/Depthborn/internal/ui/component"
)

// Stack creates overlay stack (children rendered on top of each other)
// Last child is rendered on top
// Useful for modals, tooltips, overlays
//
// Example:
//
//	Stack(
//	    Background(),
//	    Modal(),
//	    Tooltip(),
//	)
func Stack(children ...component.Component) component.Component {
	return &stack{children: children}
}

type stack struct {
	children []component.Component
}

func (s *stack) Render(ctx *component.Context) string {
	if len(s.children) == 0 {
		return ""
	}

	// For terminal UI, we can't really overlay
	// So we'll just render them vertically with separation
	// In a GUI renderer, this would do actual overlaying

	var layers []string
	for _, child := range s.children {
		if child == nil {
			continue
		}
		rendered := child.Render(ctx)
		if rendered != "" {
			layers = append(layers, rendered)
		}
	}

	// For terminal, just separate with blank line
	return strings.Join(layers, "\n\n")
}

// ZStackLayer creates z-indexed stack (named children for explicit ordering)
// Useful when you need precise control over layer order
type ZStackLayer struct {
	ZIndex int
	Child  component.Component
}

// ZStack creates z-indexed stack
func ZStack(layers ...ZStackLayer) component.Component {
	return &zStack{layers: layers}
}

type zStack struct {
	layers []ZStackLayer
}

func (z *zStack) Render(ctx *component.Context) string {
	if len(z.layers) == 0 {
		return ""
	}

	// Sort layers by z-index (lower first)
	// In terminal, we just render in order

	var rendered []string
	for _, layer := range z.layers {
		if layer.Child == nil {
			continue
		}
		content := layer.Child.Render(ctx)
		if content != "" {
			rendered = append(rendered, content)
		}
	}

	return strings.Join(rendered, "\n\n")
}
