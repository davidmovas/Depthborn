package screens

import (
	"fmt"
	"time"

	"github.com/davidmovas/Depthborn/internal/ui/component"
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
	frameState := component.UseState(ctx, 0)
	lastTime := component.UseState(ctx, time.Now())
	smoothedFPS := component.UseState(ctx, 0.0)

	component.UseEffect(ctx, func() {
		frameState.Set(frameState.Value() + 1)

		now := time.Now()
		elapsed := now.Sub(lastTime.Value()).Seconds()

		if elapsed > 0 {
			currentFPS := 1.0 / elapsed
			// Сглаживание: 70% предыдущее значение + 30% текущее
			newFPS := smoothedFPS.Value()*0.7 + currentFPS*0.3
			smoothedFPS.Set(newFPS)
		}

		lastTime.Set(now)
	}, []any{frameState.Value()})

	return primitive.Box(
		primitive.ContainerProps{
			Children: []component.Component{
				primitive.ComplexCard(primitive.CardProps{
					Children: []component.Component{
						primitive.Label(primitive.TextProps{Content: "Main Menu"}),

						primitive.VSpacer(primitive.CommonProps{}),

						primitive.Text(
							style.Sprintf("FPS: %v", style.Val(fmt.Sprintf("%.1f", smoothedFPS.Value()), style.FgCyan)),
						),
					},
				}),
			},
		},
	)
}
