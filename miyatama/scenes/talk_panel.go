package scenes

import (
	"bytes"
	"first_rpg/miyatama/assets/fonts"
	gamestatus "first_rpg/miyatama/game_status"
	"first_rpg/miyatama/util"
	"image/color"
	"log/slog"

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
		slog.Error("Store.Init",
			slog.String("TextFaceSource loading error", err.Error()),
		)
	}
	japaneseFaceSource = s
	return nil
}

func (t *TalkPanel) Update(data *gamestatus.GameData, gamepadRect *util.Rect) {
	// Panel表示領域再計算
	if t.calcPanelRectHeight != data.LayoutHeight || t.calcPanelRectWidth != data.LayoutWidth {
		slog.Info("TalkPanel.Update()",
			slog.String("recalculate", "layout"),
		)
		panelHeight := int(float32(data.LayoutHeight) * 0.2)
		panelY := gamepadRect.Top - panelHeight - 10
		panelWidth := int(float32(data.LayoutWidth) * 0.9)
		panelX := (data.LayoutWidth - panelWidth) / 2
		t.rect = &util.Rect{
			Left:   panelX,
			Top:    panelY,
			Right:  panelX + panelWidth,
			Bottom: panelY + panelHeight,
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
	vector.DrawFilledRect(screen, float32(t.rect.Left), float32(t.rect.Top), float32(t.rect.Width()), float32(t.rect.Height()), t.background, false)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(t.rect.Left+5), float64(t.rect.Top+5))
	const lineSpacing = 48
	op.LineSpacing = lineSpacing
	text.Draw(screen, data.Event.TalkTexts[data.EventMessageSeq], t.font, op)
}
