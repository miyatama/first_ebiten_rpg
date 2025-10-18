package mobile

import (
	"first_rpg/miyatama"

	"github.com/hajimehoshi/ebiten/v2/mobile"
)

var game *miyatama.Game

func init() {
	g, err := miyatama.NewGame()
	if err != nil {
		panic(err)
	}
	game = g
	if game == nil {
		print("[ebitengine] game is nil")
	} else {
		print("[ebitengine] game is not nil")

	}
	mobile.SetGame(g)
}

type AppLoggerCallback interface {
	OutputDebugLog(text string)
	OutputInfoLog(text string)
	OutputErrorLog(text string)
}

func RegisterMobileInterface(callback AppLoggerCallback) {
	game.RegisterMobileInterface(
		func(text string) { callback.OutputDebugLog(text) },
		func(text string) { callback.OutputInfoLog(text) },
		func(text string) { callback.OutputErrorLog(text) },
	)
}

func IsInitializedGame() bool {
	return game != nil
}
