package screens

import (
	"github.com/davidmovas/Depthborn/internal/app/components"
	"github.com/davidmovas/Depthborn/internal/ui/component"
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
	return components.ExampleFormModal(ctx)
}
