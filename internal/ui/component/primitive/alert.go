package primitive

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type AlertProps struct {
	BaseProps
	Type    BadgeVariant
	Title   *string
	Message string
	OnClose *func()
}

// Alert renders styled alert message
func Alert(props AlertProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		bgColor, fgColor, borderColor := getAlertColors(props.Type)

		content := []component.Component{
			Text(TextProps{
				Content:    props.Message,
				StyleProps: WithStyle(style.Fg(fgColor)),
			}),
		}

		// Add title if provided
		if props.Title != nil {
			content = append([]component.Component{
				Text(TextProps{
					Content: *props.Title,
					StyleProps: WithStyle(
						style.Merge(style.Fg(fgColor), style.Bold),
					),
				}),
				VSpacer(1),
			}, content...)
		}

		// Add close button if handler provided
		if props.OnClose != nil {
			closeBtn := Button(
				InteractiveProps{
					StyleProps: WithStyle(
						style.Merge(style.Fg(fgColor), style.Bold),
					),
					FocusProps: FocusProps{
						OnClick: *props.OnClose,
					},
				},
				"✕",
			)

			content = []component.Component{
				HStack(ContainerProps{
					ChildrenProps: Children(
						Box(ContainerProps{
							LayoutProps:   LayoutProps{Width: 60},
							ChildrenProps: ChildrenProps{Children: content},
						}),
						closeBtn,
					),
				}, 2),
			}
		}

		return Box(ContainerProps{
			StyleProps: StyleProps{
				Style: style.S(
					style.Bg(bgColor),
					style.P(2),
					style.Br(),
					style.BrColor(style.New(), borderColor),
					ApplyLayoutProps(style.New(), props.LayoutProps),
				),
			},
			ChildrenProps: ChildrenProps{Children: content},
		}).Render(ctx)
	})
}
