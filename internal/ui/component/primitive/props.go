package primitive

import (
	"bytes"

	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

// StyleProps provides styling capabilities
type StyleProps struct {
	Style      *style.Style
	ClassName  *string           // For CSS-like class system (future)
	Sx         map[string]string // Inline style overrides (future)
	FocusStyle *style.Style
	HoverStyle *style.Style // Future: hover effects
}

// FocusProps provides focus/interaction capabilities
type FocusProps struct {
	Position  *component.FocusPosition
	Hotkeys   []string
	AutoFocus *bool
	OnFocus   func()
	OnBlur    func()
	OnClick   func()
}

// LayoutProps provides sizing and spacing
type LayoutProps struct {
	Width         *int
	Height        *int
	MaxW          *int
	MaxH          *int
	Padding       *int
	PaddingX      *int
	PaddingY      *int
	PaddingTop    *int
	PaddingBottom *int
	PaddingLeft   *int
	PaddingRight  *int
	Margin        *int
	MarginX       *int
	MarginY       *int
	MarginTop     *int
	MarginBottom  *int
	MarginLeft    *int
	MarginRight   *int
}

// ContentProps provides content alignment and wrapping
type ContentProps struct {
	Align    *Align  // horizontal alignment
	VAlign   *VAlign // vertical alignment
	Wrap     *bool   // text wrapping
	Truncate *bool   // truncate overflow
	Ellipsis *bool   // show ... on truncate
}

// ChildrenProps provides children rendering
type ChildrenProps struct {
	Children []component.Component
}

func Children(children ...component.Component) ChildrenProps {
	return ChildrenProps{Children: children}
}

type Align string

const (
	AlignLeft    Align = "left"
	AlignCenter  Align = "center"
	AlignRight   Align = "right"
	AlignJustify Align = "justify"
)

type VAlign string

const (
	VAlignTop    VAlign = "top"
	VAlignMiddle VAlign = "middle"
	VAlignBottom VAlign = "bottom"
)

// BaseProps - most minimal component props
type BaseProps struct {
	StyleProps
	LayoutProps
}

// InteractiveProps - for clickable/focusable components
type InteractiveProps struct {
	StyleProps
	FocusProps
	LayoutProps
}

// ContainerProps - for layout containers
type ContainerProps struct {
	StyleProps
	LayoutProps
	ContentProps
	ChildrenProps
}

// TextProps - for text components
type TextProps struct {
	StyleProps
	LayoutProps
	ContentProps
	Content string
}

// ApplyLayoutProps applies layout props to a style
func ApplyLayoutProps(s style.Style, props LayoutProps) style.Style {
	if props.Width != nil {
		s = s.Width(*props.Width)
	}
	if props.Height != nil {
		s = s.Height(*props.Height)
	}
	if props.MaxW != nil {
		s = s.MaxWidth(*props.MaxW)
	}
	if props.MaxH != nil {
		s = s.MaxHeight(*props.MaxH)
	}

	// Padding - check specific sides first, then generic
	if props.PaddingTop != nil {
		s = s.PaddingTop(*props.PaddingTop)
	}
	if props.PaddingBottom != nil {
		s = s.PaddingBottom(*props.PaddingBottom)
	}
	if props.PaddingLeft != nil {
		s = s.PaddingLeft(*props.PaddingLeft)
	}
	if props.PaddingRight != nil {
		s = s.PaddingRight(*props.PaddingRight)
	}
	if props.PaddingX != nil {
		s = s.PaddingLeft(*props.PaddingX).PaddingRight(*props.PaddingX)
	}
	if props.PaddingY != nil {
		s = s.PaddingTop(*props.PaddingY).PaddingBottom(*props.PaddingY)
	}
	if props.Padding != nil {
		s = s.Padding(*props.Padding)
	}

	// Margin - check specific sides first, then generic
	if props.MarginTop != nil {
		s = s.MarginTop(*props.MarginTop)
	}
	if props.MarginBottom != nil {
		s = s.MarginBottom(*props.MarginBottom)
	}
	if props.MarginLeft != nil {
		s = s.MarginLeft(*props.MarginLeft)
	}
	if props.MarginRight != nil {
		s = s.MarginRight(*props.MarginRight)
	}
	if props.MarginX != nil {
		s = s.MarginLeft(*props.MarginX).MarginRight(*props.MarginX)
	}
	if props.MarginY != nil {
		s = s.MarginTop(*props.MarginY).MarginBottom(*props.MarginY)
	}
	if props.Margin != nil {
		s = s.Margin(*props.Margin)
	}

	return s
}

// ApplyStyleProps merges style props
func ApplyStyleProps(s style.Style, props StyleProps) style.Style {
	if props.Style != nil {
		s = s.Inherit(*props.Style)
	}
	return s
}

// ApplyContentProps applies alignment props
func ApplyContentProps(s style.Style, props ContentProps) style.Style {
	if props.Align != nil {
		switch *props.Align {
		case AlignLeft:
			s = s.AlignHorizontal(0)
		case AlignCenter:
			s = s.AlignHorizontal(0.5)
		case AlignRight:
			s = s.AlignHorizontal(1)
		}
	}

	if props.VAlign != nil {
		switch *props.VAlign {
		case VAlignTop:
			s = s.AlignVertical(0)
		case VAlignMiddle:
			s = s.AlignVertical(0.5)
		case VAlignBottom:
			s = s.AlignVertical(1)
		}
	}

	return s
}

// ApplyAllProps applies all prop groups to style
func ApplyAllProps(s style.Style, layout LayoutProps, styleProps StyleProps, content ContentProps) style.Style {
	s = ApplyStyleProps(s, styleProps)
	s = ApplyLayoutProps(s, layout)
	s = ApplyContentProps(s, content)
	return s
}

// RenderChildren efficiently renders array of children
func RenderChildren(ctx *component.Context, children []component.Component) string {
	if len(children) == 0 {
		return ""
	}

	// Pre-allocate buffer
	totalSize := 0
	for _, child := range children {
		if child != nil {
			totalSize += 256 // estimate
		}
	}

	b := bytes.NewBuffer(make([]byte, 0, totalSize))
	for _, child := range children {
		b.WriteString(child.Render(ctx))
	}

	return b.String()
}

// RenderChildrenWithSeparator renders children with separator between
func RenderChildrenWithSeparator(ctx *component.Context, children []component.Component, separator string) string {
	if len(children) == 0 {
		return ""
	}

	b := bytes.NewBuffer(make([]byte, 0, len(children)*256))
	for i, child := range children {
		b.WriteString(child.Render(ctx))
		if i < len(children)-1 {
			b.WriteString(separator)
		}
	}

	return b.String()
}

// Ptr returns pointer to value (helper for optional props)
func Ptr[T any](v T) *T {
	return &v
}

// WithStyle creates new StyleProps with given style
func WithStyle(s style.Style) StyleProps {
	return StyleProps{Style: &s}
}

// WithLayout creates new LayoutProps with common settings
func WithLayout(width, height, padding, margin int) LayoutProps {
	return LayoutProps{
		Width:   &width,
		Height:  &height,
		Padding: &padding,
		Margin:  &margin,
	}
}
