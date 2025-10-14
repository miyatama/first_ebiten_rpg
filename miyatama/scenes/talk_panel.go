package scenes

import (
	"bytes"
	"first_rpg/miyatama/assets/fonts"
	gamestatus "first_rpg/miyatama/game_status"
	"first_rpg/miyatama/util"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/text/language"
)

var japaneseFaceSource *text.GoTextFaceSource

type TalkPanel struct {
	background          color.Color
	calcPanelRectHeight int
	calcPanelRectWidth  int
	rect                *util.Rect
	font                *text.GoTextFace
}

func (t *TalkPanel) Init() error {
	t.background = color.RGBA{0x0, 0x0, 0x0, 0x5f}
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	japaneseFaceSource = s
	return nil
}

func (t *TalkPanel) Update(data *gamestatus.GameData) {
	// Panel表示領域再計算
	if t.calcPanelRectHeight != data.LayoutHeight || t.calcPanelRectWidth != data.LayoutWidth {
		panelHeight := int(float32(data.LayoutHeight) * 0.2)
		panelY := data.LayoutHeight - panelHeight - 10
		panelWidth := int(float32(data.LayoutWidth) * 0.9)
		panelX := (data.LayoutWidth - panelWidth) / 2
		t.rect = &util.Rect{
			Position: util.MapPosition{
				X: panelX,
				Y: panelY,
			},
			Width:  panelWidth,
			Height: panelHeight,
		}
		t.calcPanelRectHeight = data.LayoutHeight
		t.calcPanelRectWidth = data.LayoutWidth

		// 高さに応じてフォントサイズを変更
		t.font = &text.GoTextFace{
			Source:    japaneseFaceSource,
			Direction: text.DirectionLeftToRight,
			Size:      12,
			Language:  language.Japanese,
		}
	}
}

func (t *TalkPanel) Draw(screen *ebiten.Image, data *gamestatus.GameData) {
	if data.Event == nil {
		return
	}
	vector.DrawFilledRect(screen, float32(t.rect.Position.X), float32(t.rect.Position.Y), float32(t.rect.Width), float32(t.rect.Height), t.background, false)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(t.rect.Position.X+5), float64(t.rect.Position.Y+5))
	const lineSpacing = 48
	op.LineSpacing = lineSpacing
	text.Draw(screen, data.Event.TalkTexts[0], t.font, op)
}
