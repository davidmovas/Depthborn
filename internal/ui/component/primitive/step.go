package primitive

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type Step struct {
	Label    string
	Complete bool
	Active   bool
}

type StepperProps struct {
	ContainerProps

	// Steps
	Steps []Step

	// Layout
	Vertical bool

	// Visual
	ShowLabels bool
}

// Stepper creates a step indicator component
func Stepper(props StepperProps) component.Component {
	return component.Func(func(ctx *component.Context) string {
		if len(props.Steps) == 0 {
			return ""
		}

		stepComponents := make([]component.Component, 0, len(props.Steps))

		for i, step := range props.Steps {
			var icon string
			var iconColor style.Color

			if step.Complete {
				icon = "✓"
				iconColor = style.Success
			} else if step.Active {
				icon = fmt.Sprintf("%d", i+1)
				iconColor = style.Primary
			} else {
				icon = fmt.Sprintf("%d", i+1)
				iconColor = style.Grey400
			}

			// Icon circle
			iconStyle := lipgloss.NewStyle().
				Foreground(iconColor).
				Bold(true).
				Width(3).
				Align(lipgloss.Center)

			iconComp := component.Raw(iconStyle.Render(icon))

			// Label
			var stepComp component.Component
			if props.ShowLabels && step.Label != "" {
				labelStyle := lipgloss.NewStyle().
					Foreground(iconColor)

				labelComp := component.Raw(labelStyle.Render(step.Label))

				if props.Vertical {
					stepComp = VStack(ContainerProps{
						ChildrenProps: Children(iconComp, labelComp),
					}, 0)
				} else {
					stepComp = VStack(ContainerProps{
						ChildrenProps: Children(iconComp, labelComp),
					}, 0)
				}
			} else {
				stepComp = iconComp
			}

			stepComponents = append(stepComponents, stepComp)

			// Add connector line between steps (except last)
			if i < len(props.Steps)-1 {
				lineColor := style.Grey300
				if step.Complete {
					lineColor = style.Success
				}

				var connector component.Component
				if props.Vertical {
					line := lipgloss.NewStyle().
						Foreground(lineColor).
						Render(" │\n │\n │")
					connector = component.Raw(line)
				} else {
					line := lipgloss.NewStyle().
						Foreground(lineColor).
						Render(" ——— ")
					connector = component.Raw(line)
				}

				stepComponents = append(stepComponents, connector)
			}
		}

		if props.Vertical {
			return VStack(ContainerProps{
				ChildrenProps: Children(stepComponents...),
			}, 0).Render(ctx)
		}

		return HStack(ContainerProps{
			ChildrenProps: Children(stepComponents...),
		}, 0).Render(ctx)
	})
}
