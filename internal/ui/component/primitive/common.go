package primitive

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type CommonProps struct {
	ID    *string
	Class *string

	Style      *style.Style
	FocusStyle *style.Style

	AutoFocus *bool
	Position  *component.FocusPosition

	OnClick func()
	OnFocus func(string) string
	OnHover func(string) string
}

type ContainerProps struct {
	CommonProps
	Children []component.Component
	Width    *int
	Height   *int
	Padding  *int
	Margin   *int
}

type TextProps struct {
	CommonProps
	Content string
	Align   *string // "left", "center", "right"
	Wrap    *bool
}

func ptr[T any](v T) *T { return &v }
