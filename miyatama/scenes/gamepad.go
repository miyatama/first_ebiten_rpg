package scenes

import (
	"bytes"
	MiyatamaImages "first_rpg/miyatama/assets/images"
	gamestatus "first_rpg/miyatama/game_status"
	"first_rpg/miyatama/util"
	"image"
	"image/color"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	GAME_PAD_BUTTON_WIDTH = 50

	GAME_PAD_IMAGE_A_BUTTON_OFF = 0
	GAME_PAD_IMAGE_A_BUTTON_ON  = 1
	GAME_PAD_IMAGE_B_BUTTON_OFF = 2
	GAME_PAD_IMAGE_B_BUTTON_ON  = 3
	GAME_PAD_IMAGE_PAD_UP       = 4
	GAME_PAD_IMAGE_PAD_DOWN     = 5
	GAME_PAD_IMAGE_PAD_LEFT     = 6
	GAME_PAD_IMAGE_PAD_RIGHT    = 7
	GAME_PAD_IMAGE_PAD_NONE     = 8
)

type GamePad struct {
	calcPanelRectHeight int
	calcPanelRectWidth  int
	gamepadBackground   color.Color
	gamepadImages       []*ebiten.Image
	basePanelRect       *util.Rect
	gamepadRect         *util.Rect
	aButtonRect         *util.Rect
	bButtonRect         *util.Rect
	gamepadScale        float64
	buttonScale         float64
}

func (g *GamePad) Init() error {
	slog.Info("GamePad.Init()")
	img, _, err := image.Decode(bytes.NewReader(MiyatamaImages.GamePadImage))
	if err != nil {
		slog.Error("gamepad_50x50.png decode error")
		slog.String("error", err.Error())
		return err
	}
	gamepadImage := ebiten.NewImageFromImage(img)
	g.gamepadImages = []*ebiten.Image{}
	// A off
	left := 0
	top := 0
	right := GAME_PAD_BUTTON_WIDTH - 1
	bottom := GAME_PAD_BUTTON_WIDTH - 1
	g.gamepadImages = append(g.gamepadImages, gamepadImage.SubImage(image.Rect(left, top, right, bottom)).(*ebiten.Image))

	// A on
	left = GAME_PAD_BUTTON_WIDTH
	right = GAME_PAD_BUTTON_WIDTH*2 - 1
	g.gamepadImages = append(g.gamepadImages, gamepadImage.SubImage(image.Rect(left, top, right, bottom)).(*ebiten.Image))

	// B off
	left = GAME_PAD_BUTTON_WIDTH * 2
	right = GAME_PAD_BUTTON_WIDTH*3 - 1
	g.gamepadImages = append(g.gamepadImages, gamepadImage.SubImage(image.Rect(left, top, right, bottom)).(*ebiten.Image))

	// B on
	left = GAME_PAD_BUTTON_WIDTH * 3
	right = GAME_PAD_BUTTON_WIDTH*4 - 1
	g.gamepadImages = append(g.gamepadImages, gamepadImage.SubImage(image.Rect(left, top, right, bottom)).(*ebiten.Image))

	// pad up
	left = 0
	top = GAME_PAD_BUTTON_WIDTH
	right = GAME_PAD_BUTTON_WIDTH - 1
	bottom = GAME_PAD_BUTTON_WIDTH*2 - 1
	g.gamepadImages = append(g.gamepadImages, gamepadImage.SubImage(image.Rect(left, top, right, bottom)).(*ebiten.Image))

	// pad down
	left = GAME_PAD_BUTTON_WIDTH
	right = GAME_PAD_BUTTON_WIDTH*2 - 1
	g.gamepadImages = append(g.gamepadImages, gamepadImage.SubImage(image.Rect(left, top, right, bottom)).(*ebiten.Image))

	// pad left
	left = GAME_PAD_BUTTON_WIDTH * 2
	right = GAME_PAD_BUTTON_WIDTH*3 - 1
	g.gamepadImages = append(g.gamepadImages, gamepadImage.SubImage(image.Rect(left, top, right, bottom)).(*ebiten.Image))

	// pad right
	left = GAME_PAD_BUTTON_WIDTH * 3
	right = GAME_PAD_BUTTON_WIDTH*4 - 1
	g.gamepadImages = append(g.gamepadImages, gamepadImage.SubImage(image.Rect(left, top, right, bottom)).(*ebiten.Image))

	// pad none
	left = 0
	top = GAME_PAD_BUTTON_WIDTH * 2
	right = GAME_PAD_BUTTON_WIDTH - 1
	bottom = GAME_PAD_BUTTON_WIDTH*3 - 1
	g.gamepadImages = append(g.gamepadImages, gamepadImage.SubImage(image.Rect(left, top, right, bottom)).(*ebiten.Image))

	g.gamepadBackground = color.RGBA{0x25, 0x27, 0x2a, 0xff}
	return nil
}

func (g *GamePad) Update(data *gamestatus.GameData) {
	// calculate game pad position
	if g.calcPanelRectHeight != data.LayoutHeight || g.calcPanelRectWidth != data.LayoutWidth {
		slog.Info("GamePad.Update()",
			slog.String("recalculate", "layout"),
		)

		// base
		panelHeight := int(float32(data.LayoutHeight) * 0.2)
		panelY := data.LayoutHeight - panelHeight - 10
		panelWidth := int(float32(data.LayoutWidth) * 0.9)
		panelX := (data.LayoutWidth - panelWidth) / 2
		g.basePanelRect =
			&util.Rect{
				Left:   panelX,
				Top:    panelY,
				Right:  panelX + panelWidth,
				Bottom: panelY + panelHeight,
			}

		paddingWidth := 5

		// Gamepad layout
		gamepadHeight := min(g.basePanelRect.Height()-paddingWidth*2, g.basePanelRect.Width()-paddingWidth*3)
		gamepadWidth := gamepadHeight
		gamepadY := g.basePanelRect.Top + paddingWidth
		gamepadX := g.basePanelRect.Left + paddingWidth
		g.gamepadRect =
			&util.Rect{
				Left:   gamepadX,
				Top:    gamepadY,
				Right:  gamepadX + gamepadWidth,
				Bottom: gamepadY + gamepadHeight,
			}
		g.gamepadScale = float64(gamepadWidth) / float64(GAME_PAD_BUTTON_WIDTH)

		// button layout
		buttonWidth := int(float32(g.gamepadRect.Width()) * 0.5)
		buttonY := int(float32(g.basePanelRect.Top) + float32(g.basePanelRect.Height())*0.5 - float32(buttonWidth)*0.5)
		g.buttonScale = float64(buttonWidth) / float64(GAME_PAD_BUTTON_WIDTH)

		// A Button layout
		buttonX := g.basePanelRect.Right - paddingWidth - buttonWidth
		g.aButtonRect =
			&util.Rect{
				Left:   buttonX,
				Top:    buttonY,
				Right:  buttonX + buttonWidth,
				Bottom: buttonY + buttonWidth,
			}

		// B Button layout
		buttonX = g.aButtonRect.Left - 5 - buttonWidth
		g.bButtonRect =
			&util.Rect{
				Left:   buttonX,
				Top:    buttonY,
				Right:  buttonX + buttonWidth,
				Bottom: buttonY + buttonWidth,
			}

		g.calcPanelRectHeight = data.LayoutHeight
		g.calcPanelRectWidth = data.LayoutWidth

		slog.Info("GamePad.Update()",
			slog.String("base panel rect", g.basePanelRect.ToString()),
			slog.String("gamepad rect", g.gamepadRect.ToString()),
			slog.String("a button rect", g.aButtonRect.ToString()),
			slog.String("b button rect", g.bButtonRect.ToString()),
			slog.Float64("gamepad scale", g.gamepadScale),
			slog.Float64("button scale", g.buttonScale),
		)
	}

	for _, touchId := range data.TouchIds {
		touchPosition := data.TouchPositions[touchId]
		slog.Info("GamePad.Update()",
			slog.Int("touchPosition.X", touchPosition.X),
			slog.Int("touchPosition.Y", touchPosition.Y),
		)
	}
	//todo process user device tap

}

func (g *GamePad) Draw(screen *ebiten.Image, data *gamestatus.GameData) {
	if g.calcPanelRectHeight != data.LayoutHeight || g.calcPanelRectWidth != data.LayoutWidth {
		return
	}
	vector.DrawFilledRect(
		screen,
		float32(g.basePanelRect.Left),
		float32(g.basePanelRect.Top),
		float32(g.basePanelRect.Width()),
		float32(g.basePanelRect.Height()),
		g.gamepadBackground,
		false,
	)

	// gamepad
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.gamepadScale, g.gamepadScale)
	op.GeoM.Translate(float64(g.gamepadRect.Left), float64(g.gamepadRect.Top))
	screen.DrawImage(g.gamepadImages[GAME_PAD_IMAGE_PAD_NONE], op)

	// A button
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.buttonScale, g.buttonScale)
	op.GeoM.Translate(float64(g.aButtonRect.Left), float64(g.aButtonRect.Top))
	screen.DrawImage(g.gamepadImages[GAME_PAD_IMAGE_A_BUTTON_OFF], op)

	// B button
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.buttonScale, g.buttonScale)
	op.GeoM.Translate(float64(g.bButtonRect.Left), float64(g.bButtonRect.Top))
	screen.DrawImage(g.gamepadImages[GAME_PAD_IMAGE_B_BUTTON_OFF], op)
}

func (g *GamePad) GetPadRect() *util.Rect {
	return g.gamepadRect
}
