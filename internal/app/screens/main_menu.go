package screens

import (
	"fmt"

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
	styledLabel := style.Merge(
		style.Fg(style.Orange600),
		style.FgWhite,
		style.Bold,
	).Render(" David (lvl 50) ")

	return primitive.Box(primitive.ContainerProps{
		Children: []component.Component{
			primitive.Card(primitive.CardProps{
				Label:         ptr("Top Left"),
				LabelPosition: ptr(primitive.LabelTopLeft),
				Children: []component.Component{
					primitive.Text("Лейбл в левом верхнем углу"),
				},
			}),

			primitive.ActionCard(primitive.ActionCardProps{
				Label:         &styledLabel,
				LabelPosition: ptr(primitive.LabelTopLeft),
				Content: `My name is David. I'm 50 leveled soldier of the Empire.
Do you want to know more about me?`,
				Actions: []primitive.ButtonProps{
					{
						Label:   "More",
						Hotkeys: []string{"M", "enter"},
						OnClick: func() { fmt.Println("Player clicked More button") },
						OnFocus: func(s string) string { return style.FgBlue.Render(s) },
					},
					{
						Label:   "No",
						Hotkeys: []string{"N"},
						OnClick: func() { fmt.Println("Player clicked No button") },
						OnFocus: func(s string) string { return style.FgBlue.Render(s) },
					},
				},
				Padding: ptr(1),
			}),
		},
	})
}

func ptr[T any](v T) *T {
	return &v
}
