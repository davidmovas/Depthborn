package primitive

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

// ModalProps configures a modal dialog.
type ModalProps struct {
	ContainerProps

	// State
	Open    bool
	OnClose func()

	// Content
	Title       string
	Description string
	Footer      []component.Component

	// Appearance
	Size           ModalSize
	CloseOnEscape  bool
	CloseOnOverlay bool
	ShowCloseBtn   bool
	Overlay        bool

	// Advanced
	ID string
}

// Modal renders a modal dialog with overlay.
func Modal(props ModalProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		portalID := props.ID
		if portalID == "" {
			portalID = "modal_default"
		}

		if !props.Open {
			// Close portal if it was open
			ctx.Portals().Close(portalID)
			return ""
		}

		// Create the modal component that will be rendered in portal context
		modalComponent := component.Func(func(portalCtx *component.Context) string {
			width, height := portalCtx.ScreenSize()
			modalWidth := calculateModalWidth(props.Size, width)

			// Build modal content
			var content []component.Component

			// Header (title + close button)
			if props.Title != "" || props.ShowCloseBtn {
				header := buildHeader(props.Title, props.ShowCloseBtn, props.OnClose)
				content = append(content, header)
				content = append(content, VSpacer(1))
			}

			// Description
			if props.Description != "" {
				content = append(content,
					Text(TextProps{
						Content: props.Description,
						StyleProps: StyleProps{
							Style: Ptr(style.Fg(style.Grey400)),
						},
					}),
					VSpacer(1),
				)
			}

			// Children content
			content = append(content, props.Children...)

			// Footer
			if len(props.Footer) > 0 {
				content = append(content, VSpacer(1))
				content = append(content, Divider(DividerProps{
					BaseProps: BaseProps{
						LayoutProps: LayoutProps{Width: modalWidth - 4},
					},
				}))
				content = append(content, VSpacer(1))
				content = append(content, props.Footer...)
			}

			// Wrap content in card with dark background for visibility
			modalCard := Card(ContainerProps{
				ChildrenProps: Children(content...),
				LayoutProps: LayoutProps{
					Width:    modalWidth,
					MaxWidth: width - 4,
					Padding:  2,
				},
				StyleProps: StyleProps{
					Style: Ptr(style.Merge(
						style.Bg(style.Grey800),
						style.Fg(style.Grey100),
						style.BrColor(style.New(), style.Grey600),
					)),
				},
			})

			// Center modal on screen
			centeredStyle := lipgloss.NewStyle().
				Width(width).
				Height(height).
				Align(lipgloss.Center, lipgloss.Center)

			// Render with portal context so focusable elements register correctly
			return centeredStyle.Render(modalCard.Render(portalCtx))
		})

		// Register with portal system
		ctx.Portals().Open(portalID, component.LayerModal, modalComponent)

		return ""
	})
}

func buildHeader(title string, showClose bool, onClose func()) component.Component {
	return component.Func(func(ctx *component.Context) string {
		var left, right string

		if title != "" {
			left = Heading(TextProps{Content: title}).Render(ctx)
		}

		if showClose && onClose != nil {
			closeBtn := Button(InteractiveProps{
				StyleProps: WithStyle(
					style.Merge(style.Fg(style.Grey600), style.Bold),
				),
				FocusProps: FocusProps{
					OnClick: onClose,
				},
			}, "x")
			right = closeBtn.Render(ctx)
		}

		return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	})
}
