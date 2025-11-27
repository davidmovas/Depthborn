package primitive

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type LoadingProps struct {
	ContainerProps

	// Loading text
	Text string

	// Spinner props
	SpinnerVariant string
	SpinnerFrame   int

	// Layout
	Vertical bool
}

// Loading creates a loading indicator with spinner and text
func Loading(props LoadingProps) component.Component {
	spinner := Spinner(SpinnerProps{
		Variant: props.SpinnerVariant,
		Frame:   props.SpinnerFrame,
	})

	text := Text(TextProps{
		Content:    props.Text,
		StyleProps: StyleProps{Style: style.S(style.Bold)},
	})

	if props.Vertical {
		return VStack(ContainerProps{
			ChildrenProps: Children(spinner, text),
		}, 1)
	}

	return HStack(ContainerProps{
		ChildrenProps: Children(spinner, text),
	}, 2)
}
