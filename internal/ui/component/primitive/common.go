package primitive

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

func If(condition bool, comp component.Component) component.Component {
	if condition {
		return comp
	}
	return component.Empty()
}

func ptr[T any](v T) *T { return &v }

func calculateModalWidth(size ModalSize, screenWidth int) int {
	switch size {
	case ModalSizeSmall:
		return minInt(40, screenWidth-4)
	case ModalSizeMedium:
		return minInt(60, screenWidth-4)
	case ModalSizeLarge:
		return minInt(80, screenWidth-4)
	case ModalSizeFull:
		return screenWidth - 4
	default:
		return minInt(60, screenWidth-4)
	}
}

func buildModalHeader(title *string, showClose *bool, onClose func()) component.Component {
	titleComp := component.Empty()
	if title != nil {
		titleComp = Heading(TextProps{Content: *title})
	}

	closeBtn := component.Empty()
	if showClose != nil && *showClose && onClose != nil {
		closeBtn = Button(
			InteractiveProps{
				StyleProps: WithStyle(
					style.Merge(style.Fg(style.Grey600), style.Bold),
				),
				FocusProps: FocusProps{
					OnClick: onClose,
				},
			},
			"✕",
		)
	}

	return HStack(ContainerProps{
		ChildrenProps: Children(
			titleComp, closeBtn,
		),
	}, 2)
}

func buildModalOverlay(size component.ScreenSize, onClose func(), closeOnClick *bool) component.Component {
	clickable := closeOnClick != nil && *closeOnClick && onClose != nil

	overlayStyle := style.Merge(
		style.Bg(style.Black),
		// TODO: Add opacity/dim effect
	)

	overlay := Box(ContainerProps{
		LayoutProps: LayoutProps{
			Width:  Ptr(size.Width),
			Height: Ptr(size.Height),
		},
		StyleProps: StyleProps{Style: &overlayStyle},
	})

	if clickable {
		return component.MakeFocusable(overlay, component.FocusableConfig{
			CanFocus: true,
			OnActivateCallback: func() bool {
				onClose()
				return true
			},
		})
	}

	return overlay
}

func getVariantColor(variant BadgeVariant) style.Color {
	switch variant {
	case BadgePrimary:
		return style.Primary
	case BadgeSuccess:
		return style.Success
	case BadgeWarning:
		return style.Warning
	case BadgeError:
		return style.Error
	case BadgeInfo:
		return style.Info
	default:
		return style.Grey500
	}
}

func getAlertColors(variant BadgeVariant) (bg, fg, border style.Color) {
	switch variant {
	case BadgeSuccess:
		return style.Green100, style.Green800, style.Green500
	case BadgeWarning:
		return style.Yellow100, style.Yellow800, style.Yellow500
	case BadgeError:
		return style.Red100, style.Red800, style.Red500
	case BadgeInfo:
		return style.Blue100, style.Blue800, style.Blue500
	default:
		return style.Grey100, style.Grey800, style.Grey500
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
