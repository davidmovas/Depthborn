package primitive

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

// If conditionally renders a component.
func If(condition bool, comp component.Component) component.Component {
	if condition {
		return comp
	}
	return component.Empty()
}

// IfElse conditionally renders one of two components.
func IfElse(condition bool, ifTrue, ifFalse component.Component) component.Component {
	if condition {
		return ifTrue
	}
	return ifFalse
}

// --- Modal Size Helpers ---

// ModalSize represents modal dialog size.
type ModalSize string

const (
	ModalSizeSmall  ModalSize = "small"  // 40 cols
	ModalSizeMedium ModalSize = "medium" // 60 cols
	ModalSizeLarge  ModalSize = "large"  // 80 cols
	ModalSizeFull   ModalSize = "full"   // Full width
)

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

// --- Badge Variant Helpers ---

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

// --- Utility Functions ---

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
