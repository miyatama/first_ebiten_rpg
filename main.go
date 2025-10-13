package main

import (
	"first_rpg/miyatama"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game, err := miyatama.NewGame()
	if err != nil {
		panic(err)
	}

	ebiten.SetWindowSize(640, 480)
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
