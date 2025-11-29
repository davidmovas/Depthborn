package primitive

import "github.com/davidmovas/Depthborn/internal/ui/component"

// PortalProps configures the Portal component.
type PortalProps struct {
	ID       string
	Layer    component.PortalLayer
	ZIndex   int
	Open     bool
	Children []component.Component
}

// Portal renders children in a portal layer.
func Portal(props PortalProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		if !props.Open && len(props.Children) == 0 {
			return ""
		}

		// Create wrapper component for portal content
		portalContent := Box(ContainerProps{
			ChildrenProps: ChildrenProps{Children: props.Children},
		})

		// Register with portal manager
		id := props.ID
		if id == "" {
			id = "portal_default"
		}

		layer := props.Layer
		if layer == 0 {
			layer = component.LayerOverlay
		}

		ctx.Portals().OpenWithZIndex(id, layer, props.ZIndex, portalContent)

		// Return empty - actual rendering happens via PortalManager.RenderPortals()
		return ""
	})
}
