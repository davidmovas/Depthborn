package primitive

import (
	"strings"

	"github.com/davidmovas/Depthborn/internal/ui/component"
)

// List renders all children as a single string separated by newlines
func List(props ContainerProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		items := make([]string, len(props.Children))
		for i, child := range props.Children {
			items[i] = child.Render(ctx)
		}

		result := strings.Join(items, "\n")
		if props.Style != nil {
			result = props.Style.Render(result)
		}

		return result
	})
}

// Inline renders all children as a single string separated by spaces
func Inline(props ContainerProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		items := make([]string, len(props.Children))
		for i, child := range props.Children {
			items[i] = child.Render(ctx)
		}

		result := strings.Join(items, " ")
		if props.Style != nil {
			result = props.Style.Render(result)
		}

		return result
	})
}

// Grid renders all children as a grid of columns
func Grid(props ContainerProps, columns int) component.Component {
	return component.Func(func(ctx *component.Context) string {
		items := make([]string, len(props.Children))
		for i, child := range props.Children {
			items[i] = child.Render(ctx)
		}

		var rows []string
		for i := 0; i < len(items); i += columns {
			end := i + columns
			if end > len(items) {
				end = len(items)
			}
			row := strings.Join(items[i:end], "  ")
			rows = append(rows, row)
		}

		result := strings.Join(rows, "\n")
		if props.Style != nil {
			result = props.Style.Render(result)
		}

		return result
	})
}
