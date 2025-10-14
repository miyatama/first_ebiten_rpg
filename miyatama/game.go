package miyatama

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	gamestatus "first_rpg/miyatama/game_status"
	"first_rpg/miyatama/scenes"
)

type Game struct {
	scene    scenes.Scene
	gameData gamestatus.GameData
}

func (g *Game) Update() error {
	g.gameData.UserAction = keyToUserAction(inpututil.AppendPressedKeys([]ebiten.Key{}))
	if g.scene == nil {
		g.scene = &scenes.TitleScene{}
		if err := g.scene.Init(); err != nil {
			return err
		}
	}
	g.scene.Update(&g.gameData)
	return nil
}

func keyToUserAction(keys []ebiten.Key) gamestatus.UserAction {
	if len(keys) <= 0 {
		return gamestatus.USER_ACTION_NONE
	}
	switch keys[0] {
	case ebiten.KeyW:
		return gamestatus.USER_ACTION_UP
	case ebiten.KeyS:
		return gamestatus.USER_ACTION_DOWN
	case ebiten.KeyA:
		return gamestatus.USER_ACTION_LEFT
	case ebiten.KeyD:
		return gamestatus.USER_ACTION_RIGHT
	case ebiten.KeySpace:
		return gamestatus.USER_ACTION_DECIDE
	}
	return gamestatus.USER_ACTION_NONE
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.scene.Draw(screen, &g.gameData)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.gameData.ScreenWidth = outsideWidth
	g.gameData.ScreenHeight = outsideHeight
	g.gameData.LayoutWidth = outsideWidth
	g.gameData.LayoutHeight = outsideHeight
	return g.gameData.LayoutWidth, g.gameData.LayoutHeight
}
