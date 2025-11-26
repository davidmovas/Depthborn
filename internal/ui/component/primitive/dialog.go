package primitive

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type DialogProps struct {
	Open        bool
	Title       string
	Message     string
	OnConfirm   func()
	OnCancel    func()
	ConfirmText *string
	CancelText  *string
	Variant     *BadgeVariant
}

// Dialog renders simple confirmation dialog
func Dialog(props DialogProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		confirmText := "Confirm"
		if props.ConfirmText != nil {
			confirmText = *props.ConfirmText
		}

		cancelText := "Cancel"
		if props.CancelText != nil {
			cancelText = *props.CancelText
		}

		variant := BadgeDefault
		if props.Variant != nil {
			variant = *props.Variant
		}

		footer := Children(
			HStack(ContainerProps{
				ChildrenProps: Children(
					Button(
						InteractiveProps{
							StyleProps: WithStyle(
								style.Merge(
									style.Bg(style.Grey200),
									style.Fg(style.Grey700),
									style.P(1),
									style.Br(),
								),
							),
							FocusProps: FocusProps{
								OnClick: props.OnCancel,
							},
						},
						cancelText,
					),

					Button(
						InteractiveProps{
							StyleProps: WithStyle(
								style.Merge(
									style.Bg(getVariantColor(variant)),
									style.Fg(style.White),
									style.P(1),
									style.Br(),
								),
							),
							FocusProps: FocusProps{
								OnClick:   props.OnConfirm,
								AutoFocus: Ptr(true),
							},
						},
						confirmText,
					),
				),
			}, 2),
		)

		return Modal(ModalProps{
			Open:    props.Open,
			OnClose: props.OnCancel,
			Title:   &props.Title,
			Size:    ModalSizeSmall,
			ContainerProps: ContainerProps{
				ChildrenProps: Children(
					Text(TextProps{Content: props.Message}),
				),
			},
			Footer:         footer.Children,
			CloseOnEscape:  Ptr(true),
			CloseOnOverlay: Ptr(true),
			ShowCloseBtn:   Ptr(false),
		}).Render(ctx)
	})
}
