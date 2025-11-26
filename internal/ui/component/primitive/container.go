package primitive

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

func Container(props ContainerProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		childrenOutput := ""
		for _, child := range props.Children {
			childrenOutput += child.Render(ctx)
		}

		result := childrenOutput
		if props.Style != nil {
			result = props.Style.Render(result)
		}

		return result
	})
}

func Box(props ContainerProps) component.Component {
	baseStyle := lipgloss.NewStyle()

	// Применяем размеры если указаны
	if props.Width != nil {
		baseStyle = baseStyle.Width(*props.Width)
	}
	if props.Height != nil {
		baseStyle = baseStyle.Height(*props.Height)
	}
	if props.Padding != nil {
		baseStyle = baseStyle.Padding(*props.Padding)
	}
	if props.Margin != nil {
		baseStyle = baseStyle.Margin(*props.Margin)
	}
	if props.Style != nil {
		baseStyle = baseStyle.Inherit(*props.Style)
	}

	props.Style = &baseStyle
	return Container(props)
}

func Panel(props ContainerProps) component.Component {
	panelStyle := lipgloss.NewStyle().
		Background(style.Grey100).
		Padding(1).
		Border(lipgloss.RoundedBorder())

	if props.Style != nil {
		panelStyle = panelStyle.Inherit(*props.Style)
	}

	props.Style = &panelStyle
	return Box(props)
}
