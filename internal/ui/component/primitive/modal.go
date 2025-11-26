package primitive

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type ModalSize string

const (
	ModalSizeSmall  ModalSize = "small"  // 40 cols
	ModalSizeMedium ModalSize = "medium" // 60 cols
	ModalSizeLarge  ModalSize = "large"  // 80 cols
	ModalSizeFull   ModalSize = "full"   // 100% width
)

type ModalProps struct {
	ContainerProps

	// State
	Open    bool
	OnClose func()

	// Content
	Title       *string
	Description *string
	Footer      []component.Component

	// Appearance
	Size           ModalSize
	CloseOnEscape  *bool
	CloseOnOverlay *bool
	ShowCloseBtn   *bool

	// Advanced
	ID        string
	Overlay   *bool // Show dimmed overlay
	FocusTrap *bool // Trap focus inside modal
}

// Modal renders modal dialog with overlay
func Modal(props ModalProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		if !props.Open {
			return ""
		}

		screenSize := ctx.ScreenSize()

		modalWidth := calculateModalWidth(props.Size, screenSize.Width)

		var modalContent []component.Component

		if props.Title != nil || (props.ShowCloseBtn != nil && *props.ShowCloseBtn) {
			header := buildModalHeader(props.Title, props.ShowCloseBtn, props.OnClose)
			modalContent = append(modalContent, header)
			modalContent = append(modalContent, VSpacer(1))
		}

		if props.Description != nil {
			modalContent = append(modalContent,
				Label(TextProps{
					Content: *props.Description,
				}),
				VSpacer(1),
			)
		}

		modalContent = append(modalContent, props.Children...)

		if len(props.Footer) > 0 {
			modalContent = append(modalContent, VSpacer(1))
			modalContent = append(modalContent, Divider(DividerProps{
				BaseProps: BaseProps{
					LayoutProps: LayoutProps{Width: Ptr(modalWidth - 4)},
				},
			}))
			modalContent = append(modalContent, VSpacer(1))
			modalContent = append(modalContent, props.Footer...)
		}

		// Wrap in card
		modalCard := Card(ContainerProps{
			ChildrenProps: ChildrenProps{Children: modalContent},
			LayoutProps: LayoutProps{
				Width:   Ptr(modalWidth),
				MaxW:    Ptr(screenSize.Width - 4),
				Padding: Ptr(2),
			},
			StyleProps: StyleProps{
				Style: style.S(
					style.Bg(style.White),
					style.BrColor(style.New(), style.Grey400),
				),
			},
		})

		// Center modal
		centeredModal := Center(CenterProps{
			Horizontal: Ptr(true),
			Vertical:   Ptr(true),
			ContainerProps: ContainerProps{
				LayoutProps: LayoutProps{
					Width:  Ptr(screenSize.Width),
					Height: Ptr(screenSize.Height),
				},
				ChildrenProps: ChildrenProps{
					Children: []component.Component{modalCard},
				},
			},
		})

		// Add to portal layer
		return Portal(PortalProps{
			ID:    props.ID,
			Layer: component.LayerModal,
			Children: []component.Component{
				// Overlay
				If(props.Overlay == nil || *props.Overlay,
					buildModalOverlay(screenSize, props.OnClose, props.CloseOnOverlay),
				),
				// Modal content
				centeredModal,
			},
		}).Render(ctx)
	})
}
