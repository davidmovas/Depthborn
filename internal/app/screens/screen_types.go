package screens

type ScreenID string

const (
	MainMenuScreenID ScreenID = "main_menu"
)

func (id ScreenID) String() string {
	return string(id)
}
