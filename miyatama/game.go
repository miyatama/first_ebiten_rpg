package miyatama

import (
	"fmt"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	gamestatus "first_rpg/miyatama/game_status"
	"first_rpg/miyatama/scenes"
	"first_rpg/miyatama/util"
)

type Game struct {
	scene           scenes.Scene
	gameData        gamestatus.GameData
	mobileInterface gamestatus.MobileInterface
}

func (g *Game) Init() {
	g.gameData.TouchIds = []ebiten.TouchID{}
}

func (g *Game) Update() error {
	g.gameData.UserAction = keyToUserAction(inpututil.AppendPressedKeys([]ebiten.Key{}))
	if g.scene == nil {
		g.scene = &scenes.TitleScene{}
		if err := g.scene.Init(); err != nil {
			return err
		}
	}

	g.gameData.TouchIds = ebiten.AppendTouchIDs(g.gameData.TouchIds[:0])
	g.gameData.TouchPositions = map[ebiten.TouchID]util.TouchPosition{}
	for _, id := range g.gameData.TouchIds {
		x, y := ebiten.TouchPosition(id)

		slog.Info("Game.Update()",
			slog.Int("id", int(id)),
			slog.Int("x", x),
			slog.Int("y", y),
		)
		g.deviceOutputDebugLog(fmt.Sprintf("Game.Update() touch: {id: %d, x: %d, y: %d}", int(id), x, y))
		g.gameData.TouchPositions[id] = util.TouchPosition{X: x, Y: y}
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
	if g.gameData.ScreenWidth != outsideWidth || g.gameData.ScreenHeight != outsideHeight {
		slog.Info("Game.Layout()",
			slog.String("outside rect", fmt.Sprintf("{width: %d, height: %d}", outsideWidth, outsideHeight)),
		)
	}

	g.gameData.ScreenWidth = outsideWidth
	g.gameData.ScreenHeight = outsideHeight
	g.gameData.LayoutWidth = outsideWidth
	g.gameData.LayoutHeight = outsideHeight
	return g.gameData.LayoutWidth, g.gameData.LayoutHeight
}

func (g *Game) RegisterMobileInterface(
	ouptutDebug func(string),
	ouptutInfo func(string),
	ouptutError func(string),
) {
	g.mobileInterface = gamestatus.MobileInterface{
		OutputDebugLog: ouptutDebug,
		OutputInfoLog:  ouptutInfo,
		OutputErrorLog: ouptutError,
	}
}

func (g *Game) deviceOutputDebugLog(text string) {
	if g.mobileInterface.OutputDebugLog != nil {
		g.mobileInterface.OutputDebugLog(text)
	}
}
