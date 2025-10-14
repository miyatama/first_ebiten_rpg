package scenes

import (
	gamestatus "first_rpg/miyatama/game_status"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameScene interface {
	Init() error
	Update(data *gamestatus.GameData)
	Draw(screen *ebiten.Image, data *gamestatus.GameData)
	Msg() gamestatus.GameStateMsg
}
