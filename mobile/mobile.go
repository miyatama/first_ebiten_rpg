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

type MobileIface interface {
	OutputDebugLog(text string)
	OutputInfoLog(text string)
	OutputErrorLog(text string)
	SelectImportPhotos()
}

func RegisterMobileInterface(callback MobileIface) {
	game.RegisterMobileInterface(
		func(text string) { callback.OutputDebugLog(text) },
		func(text string) { callback.OutputInfoLog(text) },
		func(text string) { callback.OutputErrorLog(text) },
		func() { callback.SelectImportPhotos() },
	)
}

func IsInitializedGame() bool {
	return game != nil
}

func RegisterWorkDir(path string) {
	game.RegisterWorkDir(path)
}

func ImportPhotos(paths *miyatama.StringArray) {
	game.ImportPhotos(paths)
}

func ImportPhoto(path string) {
	game.ImportPhoto(path)
}

func RegisterWorkDir2(path string) {
	game.RegisterWorkDir(path)
}
