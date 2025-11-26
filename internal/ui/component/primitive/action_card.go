package primitive

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type ActionCardProps struct {
	CommonProps
	Label         *string
	LabelPosition *LabelPosition
	Content       string
	Actions       []ButtonProps
	Padding       *int
	Margin        *int
}

func ActionCard(props ActionCardProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		actionButtons := make([]component.Component, len(props.Actions))
		for i, action := range props.Actions {
			actionButtons[i] = Button(action)
		}

		children := []component.Component{
			Text(props.Content),

			Container(ContainerProps{
				CommonProps: CommonProps{
					Style: style.S(style.AlignRight),
				},
				Children: []component.Component{
					Inline(ContainerProps{
						Children: actionButtons,
					}),
				},
			}),
		}

		cardProps := CardProps{
			CommonProps:   props.CommonProps,
			Label:         props.Label,
			LabelPosition: props.LabelPosition,
			Padding:       props.Padding,
			Margin:        props.Margin,
			Children: []component.Component{
				Container(ContainerProps{
					Children: children,
				}),
			},
		}

		return Card(cardProps).Render(ctx)
	})
}
