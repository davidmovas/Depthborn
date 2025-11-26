package primitive

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

var (
	// T is an alias for Text
	T = Text
	// P is an alias for Text
	P = Text
)

// Text creates a styled text component
func Text(content string, styles ...style.Style) component.Component {
	return component.Func(func(ctx *component.Context) string {
		return style.Render(content, styles...)
	})
}
