package primitive

import (
	"bytes"

	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

// --- Alignment Types ---

// Align represents horizontal alignment.
type Align string

const (
	AlignLeft    Align = "left"
	AlignCenter  Align = "center"
	AlignRight   Align = "right"
	AlignJustify Align = "justify"
)

// VAlign represents vertical alignment.
type VAlign string

const (
	VAlignTop    VAlign = "top"
	VAlignMiddle VAlign = "middle"
	VAlignBottom VAlign = "bottom"
)

// --- Core Props ---

// StyleProps provides styling capabilities.
type StyleProps struct {
	Style      *style.Style // Base style
	FocusStyle *style.Style // Style when focused
	HoverStyle *style.Style // Style when hovered (future)
}

// FocusProps provides focus/interaction capabilities.
type FocusProps struct {
	ID        string                   // Explicit ID (auto-generated if empty)
	Position  *component.FocusPosition // 2D grid position for navigation
	Hotkeys   []string                 // Keys that activate this component
	AutoFocus bool                     // Focus when registered
	Disabled  bool                     // Cannot receive focus
	IsInput   bool                     // Is an input field (blocks hotkeys)
	OnFocus   func()                   // Called when focused
	OnBlur    func()                   // Called when unfocused
	OnClick   func()                   // Called on activation
}

// LayoutProps provides sizing and spacing.
type LayoutProps struct {
	Width         int
	Height        int
	MaxWidth      int
	MaxHeight     int
	Padding       int
	PaddingX      int
	PaddingY      int
	PaddingTop    int
	PaddingBottom int
	PaddingLeft   int
	PaddingRight  int
	Margin        int
	MarginX       int
	MarginY       int
	MarginTop     int
	MarginBottom  int
	MarginLeft    int
	MarginRight   int
}

// ContentProps provides content alignment and wrapping.
type ContentProps struct {
	Align    Align  // Horizontal alignment
	VAlign   VAlign // Vertical alignment
	Wrap     bool   // Text wrapping
	Truncate bool   // Truncate overflow
	Ellipsis bool   // Show ... on truncate
}

// ChildrenProps provides children rendering.
type ChildrenProps struct {
	Children []component.Component
}

// Children creates ChildrenProps from variadic components.
func Children(children ...component.Component) ChildrenProps {
	return ChildrenProps{Children: children}
}

// --- Composite Props ---

// BaseProps combines style and layout.
type BaseProps struct {
	StyleProps
	LayoutProps
}

// InteractiveProps for clickable/focusable components.
type InteractiveProps struct {
	StyleProps
	FocusProps
	LayoutProps
}

// ContainerProps for layout containers.
type ContainerProps struct {
	StyleProps
	LayoutProps
	ContentProps
	ChildrenProps
}

// TextProps for text components.
type TextProps struct {
	StyleProps
	LayoutProps
	ContentProps
	Content string
}

// --- Style Application Functions ---

// ApplyLayoutProps applies layout props to a style.
func ApplyLayoutProps(s style.Style, props LayoutProps) style.Style {
	if props.Width > 0 {
		s = s.Width(props.Width)
	}
	if props.Height > 0 {
		s = s.Height(props.Height)
	}
	if props.MaxWidth > 0 {
		s = s.MaxWidth(props.MaxWidth)
	}
	if props.MaxHeight > 0 {
		s = s.MaxHeight(props.MaxHeight)
	}

	// Padding - specific sides first, then generic
	if props.PaddingTop > 0 {
		s = s.PaddingTop(props.PaddingTop)
	}
	if props.PaddingBottom > 0 {
		s = s.PaddingBottom(props.PaddingBottom)
	}
	if props.PaddingLeft > 0 {
		s = s.PaddingLeft(props.PaddingLeft)
	}
	if props.PaddingRight > 0 {
		s = s.PaddingRight(props.PaddingRight)
	}
	if props.PaddingX > 0 {
		s = s.PaddingLeft(props.PaddingX).PaddingRight(props.PaddingX)
	}
	if props.PaddingY > 0 {
		s = s.PaddingTop(props.PaddingY).PaddingBottom(props.PaddingY)
	}
	if props.Padding > 0 {
		s = s.Padding(props.Padding)
	}

	// Margin - specific sides first, then generic
	if props.MarginTop > 0 {
		s = s.MarginTop(props.MarginTop)
	}
	if props.MarginBottom > 0 {
		s = s.MarginBottom(props.MarginBottom)
	}
	if props.MarginLeft > 0 {
		s = s.MarginLeft(props.MarginLeft)
	}
	if props.MarginRight > 0 {
		s = s.MarginRight(props.MarginRight)
	}
	if props.MarginX > 0 {
		s = s.MarginLeft(props.MarginX).MarginRight(props.MarginX)
	}
	if props.MarginY > 0 {
		s = s.MarginTop(props.MarginY).MarginBottom(props.MarginY)
	}
	if props.Margin > 0 {
		s = s.Margin(props.Margin)
	}

	return s
}

// ApplyStyleProps merges style props.
func ApplyStyleProps(s style.Style, props StyleProps) style.Style {
	if props.Style != nil {
		s = s.Inherit(*props.Style)
	}
	return s
}

// ApplyContentProps applies alignment props.
func ApplyContentProps(s style.Style, props ContentProps) style.Style {
	switch props.Align {
	case AlignLeft:
		s = s.AlignHorizontal(0)
	case AlignCenter:
		s = s.AlignHorizontal(0.5)
	case AlignRight:
		s = s.AlignHorizontal(1)
	}

	switch props.VAlign {
	case VAlignTop:
		s = s.AlignVertical(0)
	case VAlignMiddle:
		s = s.AlignVertical(0.5)
	case VAlignBottom:
		s = s.AlignVertical(1)
	}

	return s
}

// ApplyAllProps applies all prop groups to style.
func ApplyAllProps(s style.Style, layout LayoutProps, styleProps StyleProps, content ContentProps) style.Style {
	s = ApplyStyleProps(s, styleProps)
	s = ApplyLayoutProps(s, layout)
	s = ApplyContentProps(s, content)
	return s
}

// --- Children Rendering ---

// RenderChildren renders array of children.
func RenderChildren(ctx *component.Context, children []component.Component) string {
	if len(children) == 0 {
		return ""
	}

	var b bytes.Buffer
	for _, child := range children {
		if child != nil {
			b.WriteString(child.Render(ctx))
		}
	}

	return b.String()
}

// RenderChildrenWithSeparator renders children with separator between.
func RenderChildrenWithSeparator(ctx *component.Context, children []component.Component, separator string) string {
	if len(children) == 0 {
		return ""
	}

	var b bytes.Buffer
	for i, child := range children {
		if child != nil {
			b.WriteString(child.Render(ctx))
			if i < len(children)-1 {
				b.WriteString(separator)
			}
		}
	}

	return b.String()
}

// --- Helper Functions ---

// Ptr returns pointer to value (for optional props).
func Ptr[T any](v T) *T {
	return &v
}

// WithStyle creates StyleProps with given style.
func WithStyle(s style.Style) StyleProps {
	return StyleProps{Style: &s}
}

// WithLayout creates LayoutProps with common settings.
func WithLayout(width, height, padding, margin int) LayoutProps {
	return LayoutProps{
		Width:   width,
		Height:  height,
		Padding: padding,
		Margin:  margin,
	}
}
