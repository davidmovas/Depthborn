package primitive

import "github.com/davidmovas/Depthborn/internal/ui/component"

type PortalProps struct {
	Children []component.Component
	Layer    component.PortalLayer
	ID       string
	ZIndex   *int
}

// Portal renders children in specified layer (modal, toast, tooltip)
func Portal(props PortalProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		if len(props.Children) == 0 {
			return ""
		}

		zIndex := 1000
		if props.ZIndex != nil {
			zIndex = *props.ZIndex
		}

		// Create wrapper component for portal content
		portalContent := Box(ContainerProps{
			ChildrenProps: ChildrenProps{Children: props.Children},
		})

		// Register in portal manager
		ctx.PortalManager().Register(component.PortalEntry{
			ID:        props.ID,
			Layer:     props.Layer,
			Component: portalContent,
			ZIndex:    zIndex,
		})

		// Return empty - actual rendering happens in portal layer
		return ""
	})
}
