package screens

import (
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/component/layout"
	"github.com/davidmovas/Depthborn/internal/ui/component/primitive"
	"github.com/davidmovas/Depthborn/internal/ui/navigation"
	"github.com/davidmovas/Depthborn/internal/ui/style"
)

type mainMenuScreen struct {
	*navigation.BaseScreen
}

func NewMainMenuScreen() navigation.Screen {
	return &mainMenuScreen{
		BaseScreen: navigation.NewBaseScreen("main_menu"),
	}
}

func (m *mainMenuScreen) Render(ctx *component.Context) component.Component {
	countState := component.UseState(ctx, 0, "count")
	healthState := component.UseState(ctx, 10, "health")

	component.UseEffect(ctx, func() {
		healthState.Set(healthState.Value() + 10)
	}, []any{countState.Value()})

	return layout.V(
		primitive.Button(primitive.ButtonProps{
			Label: style.GradientFg("Gradient button from orange to red",
				style.Orange300,
				style.Red500,
				style.Purple700,
			),
		}),

		primitive.Button(primitive.ButtonProps{
			Label: style.Sprintf("🎮 New Game (%v)",
				style.Val(countState.Value(), style.Bold)),
			OnClick:    func() { countState.Set(countState.Value() + 1) },
			FocusStyle: style.S(style.Padding1, style.BorderRounded),
		}),

		primitive.Text(
			style.Sprintf("Player: %v | Score: %v | Level: %v",
				style.Val("John", style.Fg(style.Blue500), style.Bold),
				style.Val(12500, style.Fg(style.Green500), style.Bold),
				style.Val(5, style.Fg(style.Orange500), style.Underline),
			),
		),

		primitive.Text(
			style.Sprintf("Status: %v (%v)",
				style.Val("Online", style.Fg(style.Success), style.Bold),
				style.Val("2 players connected"),
			),
		),

		primitive.Button(primitive.ButtonProps{
			Label: style.GradientFg(
				style.Sprintf("🎯 Mission %v: %v/%v completed",
					style.Val("A", style.Fg(style.Primary), style.Bold),
					style.Val(3, style.Fg(style.Success), style.Bold),
					style.Val(5, style.FgGray, style.TextMuted),
				),
				style.Orange300,
				style.Red500,
				style.Purple700,
			),
			Style: style.S(style.Padding1, style.BorderRounded),
		}),

		primitive.Text(
			style.Sprintf("Health: %v | Mana: %v | Gold: %v",
				style.Val(healthState.Value(), style.FgRed, style.Bold),
				style.Val(42.5, style.FgBlue, style.Bold),
				style.Val(1250, style.FgYellow, style.Bold),
			),
		),

		primitive.Text(
			style.NewGradient(
				style.Cyan500,
				style.Blue500,
				style.Purple500,
			).Style(
				style.Merge(
					style.Br(),
					style.P(2),
					style.Bold,
				),
				"border",
			).Render(
				style.Sprintf("⚔️  Battle Stats  ⚔️\nDamage: %v | Defense: %v",
					style.Val(150, style.FgRed, style.Bold),
					style.Val(85, style.FgGreen, style.Bold),
				),
			),
		),

		primitive.Text(
			style.GradientBorderBox(
				"Красивая рамка!",
				200,
				style.Cyan500, style.Orange500,
			),
		),

		// С направлением
		primitive.Text(
			style.NewGradient(style.Red500, style.Orange500, style.Yellow500, style.Green500).
				Direction(style.DirectionBorderClockwise).
				BorderGradientBox("Радужная рамка", 25),
		),

		// Многострочный контент
		primitive.Text(
			style.GradientBorderBox(
				"Строка 1\nСтрока 2\nСтрока 3",
				30,
				style.Orange300, style.Red500, style.Purple700,
			),
		),
	)
}
