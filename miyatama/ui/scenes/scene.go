package scenes

import (
	gamestatus "first_rpg/miyatama/ui/game_status"

	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Init() error
	Update(data *gamestatus.GameData)
	Draw(screen *ebiten.Image, data *gamestatus.GameData)
}
