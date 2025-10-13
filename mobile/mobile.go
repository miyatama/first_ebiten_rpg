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

/*
type AppLoggerCallback interface {
	OutputLog()
}

func RegisterAppLoggerCallback(callback AppLoggerCallback) {
	game.RegisterAppLoggerCallback(func() { callback.OutputLog() })
}

type AdsCallback interface {
	ShowAds()
}

func RegisterAdsCallback(callback AdsCallback) {
	game.RegisterAdsCallback(func() { callback.ShowAds() })
}

func IsInitializedGame() bool {
	return game != nil
}
*/
