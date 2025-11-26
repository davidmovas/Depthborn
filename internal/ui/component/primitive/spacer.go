package primitive

import (
	"strings"

	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type SpacerProps struct {
	CommonProps
	Width  *int
	Height *int
	Size   *int
}

func Spacer(props SpacerProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		width := 1
		height := 1

		if props.Size != nil {
			width = *props.Size
			height = *props.Size
		}
		if props.Width != nil {
			width = *props.Width
		}
		if props.Height != nil {
			height = *props.Height
		}

		content := strings.Repeat(" ", width)
		if height > 1 {
			lines := make([]string, height)
			for i := range lines {
				lines[i] = content
			}
			content = strings.Join(lines, "\n")
		}

		if props.Style != nil {
			content = props.Style.Render(content)
		}

		return content
	})
}

func HSpacer(props CommonProps) component.Component {
	return Spacer(SpacerProps{
		CommonProps: props,
		Width:       ptr(style.Space4),
		Height:      ptr(1),
	})
}

func VSpacer(props CommonProps) component.Component {
	return Spacer(SpacerProps{
		CommonProps: props,
		Width:       ptr(1),
		Height:      ptr(style.Space2),
	})
}

func FlexSpacer(props CommonProps) component.Component {
	return Spacer(SpacerProps{
		CommonProps: props,
		Width:       ptr(0),
	})
}
