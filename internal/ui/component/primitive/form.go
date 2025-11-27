package primitive

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type FormField struct {
	Name      string
	Label     string
	Component component.Component
	Required  bool
	Validator func(string) error
}

type FormProps struct {
	ContainerProps

	// Form fields
	Fields []FormField

	// Form state
	Values map[string]string
	Errors map[string]string

	// Callbacks
	OnSubmit func(values map[string]string) error
	OnCancel func()

	// Visual
	Title       string
	SubmitLabel string
	CancelLabel string

	// Layout
	Vertical bool // Stack fields vertically (default: true)
}

// Form creates a form with multiple input fields
func Form(props FormProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		if props.Values == nil {
			props.Values = make(map[string]string)
		}
		if props.Errors == nil {
			props.Errors = make(map[string]string)
		}

		fieldComponents := make([]component.Component, 0, len(props.Fields))

		for _, field := range props.Fields {
			labelText := field.Label
			if field.Required {
				labelText += " *"
			}

			labelComp := Text(TextProps{
				Content: labelText,
				StyleProps: StyleProps{
					Style: style.S(style.Bold, style.Fg(style.Grey700)),
				},
			})

			// Error for this field
			_ = props.Errors[field.Name]

			// Field component with error
			fieldWithError := VStack(ContainerProps{
				ChildrenProps: Children(
					labelComp,
					field.Component),
			}, 0)

			fieldComponents = append(fieldComponents, fieldWithError)
		}

		submitLabel := props.SubmitLabel
		if submitLabel == "" {
			submitLabel = "Submit"
		}

		cancelLabel := props.CancelLabel
		if cancelLabel == "" {
			cancelLabel = "Cancel"
		}

		buttonRow := HStack(ContainerProps{
			ChildrenProps: Children(
				Button(InteractiveProps{
					StyleProps: WithStyle(style.Merge(style.Bg(style.Primary), style.Fg(style.White))),
					FocusProps: FocusProps{
						OnClick: func() {
							if props.OnSubmit != nil {
								_ = props.OnSubmit(props.Values)
							}
						},
					},
				}, submitLabel),
				Button(InteractiveProps{
					FocusProps: FocusProps{
						OnClick: func() {
							if props.OnCancel != nil {
								props.OnCancel()
							}
						},
					},
				}, cancelLabel),
			),
		}, 2)

		allComponents := append(fieldComponents, buttonRow)

		formContent := VStack(ContainerProps{
			ChildrenProps: Children(allComponents...),
		}, 1)

		if props.Title != "" {
			return Panel(
				ContainerProps{
					ChildrenProps: Children(formContent),
				},
			).Render(ctx)
		}

		return formContent.Render(ctx)
	})
}
