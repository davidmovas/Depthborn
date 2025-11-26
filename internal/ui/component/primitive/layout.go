package primitive

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type StackDirection string

const (
	StackVertical   StackDirection = "vertical"
	StackHorizontal StackDirection = "horizontal"
)

type StackProps struct {
	ContainerProps
	Direction StackDirection
	Gap       *int
	Reverse   *bool
}

// Stack renders children in vertical or horizontal layout
func Stack(props StackProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		if len(props.Children) == 0 {
			return ""
		}

		gap := 0
		if props.Gap != nil {
			gap = *props.Gap
		}

		reverse := false
		if props.Reverse != nil {
			reverse = *props.Reverse
		}

		children := props.Children
		if reverse {
			// Reverse children order
			reversed := make([]component.Component, len(children))
			for i := range children {
				reversed[i] = children[len(children)-1-i]
			}
			children = reversed
		}

		var result strings.Builder

		if props.Direction == StackHorizontal {
			// Horizontal layout
			for i, child := range children {
				if i > 0 && gap > 0 {
					result.WriteString(strings.Repeat(" ", gap))
				}
				result.WriteString(child.Render(ctx))
			}
		} else {
			// Vertical layout (default)
			for i, child := range children {
				if i > 0 && gap > 0 {
					result.WriteString(strings.Repeat("\n", gap))
				}
				result.WriteString(child.Render(ctx))
				if i < len(children)-1 {
					result.WriteString("\n")
				}
			}
		}

		content := result.String()

		s := style.New()
		s = ApplyAllProps(s, props.LayoutProps, props.StyleProps, props.ContentProps)

		return s.Render(content)
	})
}

// VStack renders children vertically
func VStack(props ContainerProps, gap int) component.Component {
	return Stack(StackProps{
		ContainerProps: props,
		Direction:      StackVertical,
		Gap:            &gap,
	})
}

// HStack renders children horizontally
func HStack(props ContainerProps, gap int) component.Component {
	return Stack(StackProps{
		ContainerProps: props,
		Direction:      StackHorizontal,
		Gap:            &gap,
	})
}

type FlexProps struct {
	ContainerProps
	Direction FlexDirection
	Justify   *FlexJustify
	Align     *FlexAlign
	Wrap      *bool
	Gap       *int
}

type FlexDirection string

const (
	FlexRow    FlexDirection = "row"
	FlexColumn FlexDirection = "column"
)

type FlexJustify string

const (
	JustifyStart   FlexJustify = "start"
	JustifyCenter  FlexJustify = "center"
	JustifyEnd     FlexJustify = "end"
	JustifyBetween FlexJustify = "between"
	JustifyAround  FlexJustify = "around"
	JustifyEvenly  FlexJustify = "evenly"
)

type FlexAlign string

const (
	FlexStart   FlexAlign = "start"
	FlexCenter  FlexAlign = "center"
	FlexEnd     FlexAlign = "end"
	FlexStretch FlexAlign = "stretch"
)

// Flex renders flexible layout container
func Flex(props FlexProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		if len(props.Children) == 0 {
			return ""
		}

		gap := 0
		if props.Gap != nil {
			gap = *props.Gap
		}

		isColumn := props.Direction == FlexColumn

		// Render children
		renderedChildren := make([]string, 0, len(props.Children))
		for _, child := range props.Children {
			if child != nil {
				renderedChildren = append(renderedChildren, child.Render(ctx))
			}
		}

		// Join with gap
		var result string
		if isColumn {
			separator := "\n"
			if gap > 0 {
				separator = strings.Repeat("\n", gap+1)
			}
			result = strings.Join(renderedChildren, separator)
		} else {
			separator := ""
			if gap > 0 {
				separator = strings.Repeat(" ", gap)
			}
			result = lipgloss.JoinHorizontal(lipgloss.Top, renderedChildren...)
			if gap > 0 && len(renderedChildren) > 1 {
				// Manual gap insertion for horizontal
				parts := make([]string, 0, len(renderedChildren)*2-1)
				for i, part := range renderedChildren {
					if i > 0 {
						parts = append(parts, separator)
					}
					parts = append(parts, part)
				}
				result = lipgloss.JoinHorizontal(lipgloss.Top, parts...)
			}
		}

		s := style.New()
		s = ApplyAllProps(s, props.LayoutProps, props.StyleProps, props.ContentProps)

		// Apply alignment
		if props.Justify != nil || props.Align != nil {
			// This is simplified - full implementation would need layout calculations
			if props.Align != nil {
				switch *props.Align {
				case FlexCenter:
					s = s.AlignHorizontal(lipgloss.Center)
				case FlexEnd:
					s = s.AlignHorizontal(lipgloss.Right)
				}
			}
		}

		return s.Render(result)
	})
}

type GridProps struct {
	ContainerProps
	Columns  int
	Gap      *int
	ColGap   *int
	RowGap   *int
	AutoFlow *GridAutoFlow
}

type GridAutoFlow string

const (
	GridFlowRow    GridAutoFlow = "row"
	GridFlowColumn GridAutoFlow = "column"
)

// Grid renders children in grid layout
func Grid(props GridProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		if len(props.Children) == 0 {
			return ""
		}

		columns := props.Columns
		if columns <= 0 {
			columns = 1
		}

		gap := 1
		if props.Gap != nil {
			gap = *props.Gap
		}

		colGap := gap
		if props.ColGap != nil {
			colGap = *props.ColGap
		}

		rowGap := gap
		if props.RowGap != nil {
			rowGap = *props.RowGap
		}

		// Render all children
		renderedChildren := make([]string, 0, len(props.Children))
		for _, child := range props.Children {
			if child != nil {
				renderedChildren = append(renderedChildren, child.Render(ctx))
			}
		}

		rows := make([][]string, 0)
		for i := 0; i < len(renderedChildren); i += columns {
			end := i + columns
			if end > len(renderedChildren) {
				end = len(renderedChildren)
			}
			rows = append(rows, renderedChildren[i:end])
		}

		var result strings.Builder
		for rowIdx, row := range rows {
			if rowIdx > 0 && rowGap > 0 {
				result.WriteString(strings.Repeat("\n", rowGap))
			}

			rowContent := lipgloss.JoinHorizontal(
				lipgloss.Top,
				row...,
			)

			if colGap > 0 && len(row) > 1 {
				parts := make([]string, 0, len(row)*2-1)
				for i, part := range row {
					if i > 0 {
						parts = append(parts, strings.Repeat(" ", colGap))
					}
					parts = append(parts, part)
				}
				rowContent = lipgloss.JoinHorizontal(lipgloss.Top, parts...)
			}

			result.WriteString(rowContent)
			if rowIdx < len(rows)-1 {
				result.WriteString("\n")
			}
		}

		content := result.String()

		// Apply container styles
		s := style.New()
		s = ApplyAllProps(s, props.LayoutProps, props.StyleProps, props.ContentProps)

		return s.Render(content)
	})
}

type CenterProps struct {
	ContainerProps
	Horizontal *bool
	Vertical   *bool
}

// Center centers content horizontally and/or vertically
func Center(props CenterProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		content := RenderChildren(ctx, props.Children)

		horizontal := true
		if props.Horizontal != nil {
			horizontal = *props.Horizontal
		}

		vertical := false
		if props.Vertical != nil {
			vertical = *props.Vertical
		}

		s := style.New()

		if horizontal {
			s = s.AlignHorizontal(lipgloss.Center)
		}

		if vertical {
			s = s.AlignVertical(lipgloss.Center)
		}

		s = ApplyAllProps(s, props.LayoutProps, props.StyleProps, props.ContentProps)

		return s.Render(content)
	})
}
