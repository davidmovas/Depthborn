package screens

import (
	"fmt"

	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/component/layout"
	"github.com/davidmovas/Depthborn/internal/ui/component/primitive"
	"github.com/davidmovas/Depthborn/internal/ui/navigation"
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
	// Example with counter state (key MUST be stable)
	countState := component.UseState(ctx, 0, "count")

	component.UseEffect(ctx, func() {
		fmt.Printf("Count effect triggered, count=%d\n", countState.Value())
	}, []any{countState.Value()})

	return layout.V(
		primitive.Title("My Game"),
		primitive.Textf("Count: %d", countState.Value()),
		component.Raw("\n"),

		// Buttons with STABLE IDs
		primitive.Button(primitive.ButtonConfig{
			ID:    "btn_new_game", // <- Stable ID!
			Label: fmt.Sprintf("New game (%d)", countState.Value()),
			Key:   "N",
			OnClick: func() {
				countState.Set(countState.Value() + 1)
			},
		}),

		primitive.Button(primitive.ButtonConfig{
			ID:    "btn_continue",
			Label: "Continue",
			OnClick: func() {
				fmt.Println("Continue clicked")
			},
		}),

		primitive.Button(primitive.ButtonConfig{
			ID:    "btn_exit",
			Label: "Exit",
			Key:   "Q",
			OnClick: func() {
				fmt.Println("Exit clicked")
			},
		}),
	)
}
