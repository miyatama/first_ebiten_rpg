package scenes

import (
	gamestatus "first_rpg/miyatama/game_status"

	"github.com/hajimehoshi/ebiten/v2"
)

type Store struct {
	Id int
}

func (s *Store) Init() error {
	return nil
}

func (s *Store) Update(data *gamestatus.GameData) {

}

func (s *Store) Draw(screen *ebiten.Image, data *gamestatus.GameData) {

}
