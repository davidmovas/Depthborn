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
	rainbow := []style.Color{
		style.Red500,
		style.Orange500,
		style.Yellow500,
		style.Green500,
		style.Blue500,
		style.Purple500,
	}

	fire := []style.Color{
		style.Red500,
		style.Orange500,
		style.Yellow500,
	}

	cyber := []style.Color{
		style.Cyan500,
		style.Purple500,
		style.Pink500,
	}

	return layout.V(
		primitive.Text(
			style.GradientFg(
				"🌈 Rainbow Header",
				style.AnimatedGradient(ctx, 40, rainbow...)...,
			),
		),

		primitive.Text(
			style.GradientBorderBox(
				style.GradientFg(
					"🔥 Fire inside\n💜 Cyber border",
					style.AnimatedGradientWave(ctx, 15, fire...)...,
				),
				style.AnimatedGradient(ctx, 80, cyber...),
				style.BoxBorderDouble,
			),
		),

		primitive.Text(
			style.GradientBorderBox(
				style.GradientFg(
					"← Reverse content\n→ Normal border",
					style.AnimatedGradientReverse(ctx, 20, rainbow...)...,
				),
				style.AnimatedGradient(ctx, 60, rainbow...),
				style.BoxBorderThick,
			),
		),

		primitive.Text(
			style.GradientFg(
				"✨ Everything is animated!",
				style.AnimatedGradientPingPong(ctx, 25, rainbow...)...,
			),
		),
	)
}
