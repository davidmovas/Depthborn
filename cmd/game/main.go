package main

import (
	"github.com/davidmovas/Depthborn/internal/app/screens"
	"github.com/davidmovas/Depthborn/internal/ui/navigation"
	"github.com/davidmovas/Depthborn/internal/ui/renderer"
	"github.com/davidmovas/Depthborn/internal/ui/renderer/tea"
)

func main() {
	nav := navigation.NewNavigator()

	nav.Register(screens.MainMenuScreenID.String(), screens.NewMainMenuScreen)

	if err := nav.Open(screens.MainMenuScreenID.String(), map[string]any{}); err != nil {
		panic(err)
	}

	r := tea.New(
		renderer.Config{Title: "Depthborn"},
		nav,
	)

	if err := r.Init(); err != nil {
		panic(err)
	}

	defer func() {
		if err := r.Stop(); err != nil {
			panic(err)
		}
	}()

	if err := r.Run(); err != nil {
		panic(err)
	}
}
