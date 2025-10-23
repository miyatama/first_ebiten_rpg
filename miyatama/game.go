package miyatama

import (
	"bytes"
	"fmt"
	"log/slog"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/text/language"

	miyatamaAudio "first_rpg/miyatama/assets/audio"
	"first_rpg/miyatama/assets/fonts"
	gamestatus "first_rpg/miyatama/game_status"
	"first_rpg/miyatama/scenes"
	"first_rpg/miyatama/util"
)

type Game struct {
	scene           scenes.GameScene
	gameData        gamestatus.GameData
	mobileInterface gamestatus.MobileInterface
}

func (g *Game) Init() {
	g.gameData.TouchIds = []ebiten.TouchID{}
	g.gameData.GOOS = runtime.GOOS
	g.gameData.Environemnt = &gamestatus.Environment{}
	g.gameData.AudioContext = audio.NewContext(miyatamaAudio.DEFAULT_SAMPLE_RATE)
	slog.Info("Game.Init()",
		slog.String("GOOS", g.gameData.GOOS))
}

func (g *Game) Update() error {
	g.gameData.UserAction = keyToUserAction(inpututil.AppendPressedKeys([]ebiten.Key{}))
	if g.scene == nil {
		g.scene = &scenes.TitleScene{}
		if err := g.scene.Init(&g.gameData); err != nil {
			return err
		}
	} else {
		switch g.scene.Msg() {

		case gamestatus.GAME_STATE_MSG_TITLE:
			g.scene = &scenes.TitleScene{}
			if err := g.scene.Init(&g.gameData); err != nil {
				return err
			}
		case gamestatus.GAME_STATE_MSG_HOUSE:
			g.scene = &scenes.HouseScene{}
			if err := g.scene.Init(&g.gameData); err != nil {
				return err
			}
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
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()))

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if g.gameData.ScreenWidth == outsideWidth && g.gameData.ScreenHeight == outsideHeight {
		return g.gameData.LayoutWidth, g.gameData.LayoutHeight
	}

	slog.Info("Game.Layout()",
		slog.String("outside rect", fmt.Sprintf("{width: %d, height: %d}", outsideWidth, outsideHeight)),
	)

	g.gameData.ScreenWidth = outsideWidth
	g.gameData.ScreenHeight = outsideHeight
	g.gameData.LayoutWidth = outsideWidth
	g.gameData.LayoutHeight = outsideHeight

	// デバイスの縦型横型に合わせてテキストサイズを決定
	japaneseFaceSource, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		slog.Error("Store.Init",
			slog.String("TextFaceSource loading error", err.Error()),
		)
	}
	font := &text.GoTextFace{
		Source:    japaneseFaceSource,
		Direction: text.DirectionLeftToRight,
		Size:      12,
		Language:  language.Japanese,
	}
	messagePanelWidth := 0
	if outsideHeight > outsideWidth {
		// (padding(5px) + border(5px) + margin(5px)) * 2 = 30
		messagePanelWidth = outsideWidth - 30
	} else {
		// 横幅の8割 - (border(5px) + mergin(5px)) * 2
		messagePanelWidth = int(float32(outsideWidth)*0.8) - 20
	}
	largeSize, err := getTextSize("０１２３４０１２３４", float64(messagePanelWidth), font)
	if err != nil {
		slog.Error("Game.Layout() get large text font size",
			slog.String("error", err.Error()),
		)
	}
	g.gameData.TextSizeLarge = largeSize
	middleSize, err := getTextSize("０１２３４０１２３４０１２３４０１２３４", float64(messagePanelWidth), font)
	if err != nil {
		slog.Error("Game.Layout() get middle text font size",
			slog.String("error", err.Error()),
		)
	}
	g.gameData.TextSizeMiddle = middleSize
	smallSize, err := getTextSize("０１２３４０１２３４０１２３４０１２３４０１２３４０１２３４０１２３４０１２３４", float64(messagePanelWidth), font)
	if err != nil {
		slog.Error("Game.Layout() get small text font size",
			slog.String("error", err.Error()),
		)
	}
	g.gameData.TextSizeSmall = smallSize
	slog.Info("Game.Layout()",
		slog.Int("messagePanelWidth", messagePanelWidth),
		slog.String("text size", fmt.Sprintf("{large: %f, middle: %f, small: %f}", largeSize, middleSize, smallSize)),
	)

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

func (g *Game) RegisterWorkDir(path string) {
	slog.Info("Game.RegisterWorkDir()",
		slog.String("path", path),
	)
	g.gameData.Environemnt.WorkDirPath = path
}

func getTextSize(sampleText string, panelWidth float64, font *text.GoTextFace) (float64, error) {
	for i := 100; i > 0; i-- {
		size := float64(i)
		font.Size = size
		lineSpacing := size * 1.2
		w, _ := text.Measure(sampleText, font, lineSpacing)
		slog.Info("Game.Layout() get text size",
			slog.String("text", sampleText),
			slog.Float64("font size", size),
			slog.Float64("width", w),
		)
		if w <= panelWidth {
			return size, nil
		}
	}
	return 0, fmt.Errorf("text size not found")
}
