package screens

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davidmovas/Depthborn/internal/ui/component"
	"github.com/davidmovas/Depthborn/internal/ui/component/layout"
	"github.com/davidmovas/Depthborn/internal/ui/component/primitive"
	"github.com/davidmovas/Depthborn/internal/ui/navigation"
)

var _ navigation.Screen = (*mainMenuScreen)(nil)

type mainMenuScreen struct {
	id ScreenID
}

func NewMainMenuScreen() navigation.Screen {
	return &mainMenuScreen{
		id: MainMenuScreenID,
	}
}

func (m *mainMenuScreen) ID() string {
	return string(m.id)
}

func (m *mainMenuScreen) OnInit() {}

func (m *mainMenuScreen) OnEnter(params map[string]any) {}

func (m *mainMenuScreen) OnUpdateStart() {}

func (m *mainMenuScreen) OnUpdate(dt time.Duration) {}

func (m *mainMenuScreen) OnUpdateEnd() {}

func (m *mainMenuScreen) OnPause() {}

func (m *mainMenuScreen) OnResume() {}

func (m *mainMenuScreen) OnExit() {}

func (m *mainMenuScreen) Render(ctx *component.Context) component.Component {
	var count = 10

	countState := component.UseState(ctx, count)

	time.AfterFunc(time.Second*5, func() {
		countState.Set(countState.Value() + 10)
	})

	return layout.V(
		primitive.Title("My Game"),
		primitive.Textf("Count: %d", countState.Value()),
		primitive.PrimaryButton("New Game", "N", func() {
			fmt.Println("New game clicked")
		}),
		primitive.Button(primitive.ButtonConfig{
			Label: "Continue",
			Key:   "C",
			OnClick: func() {
				fmt.Println("Continue clicked")
			},
		}),
		primitive.Button(primitive.ButtonConfig{
			Label: "Exit",
			Key:   "Q",
			OnClick: func() {
				fmt.Println("Exit clicked")
			},
		}),
	)
}

func (m *mainMenuScreen) HandleInput(msg any) bool {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "n":
			fmt.Println("New game")
			return true
		case "c":
			fmt.Println("Continue")
			return true
		case "q":
			fmt.Println("Exit")
			return true
		}
	}
	return false
}

func (m *mainMenuScreen) CanClose() bool {
	return true
}
