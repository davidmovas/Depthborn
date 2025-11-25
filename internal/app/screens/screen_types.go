package screens

type ScreenID string

const (
	MainMenuScreenID ScreenID = "main_menu_screen"
)

func (id ScreenID) String() string {
	return string(id)
}
