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

type StoreInfo struct {
	showpName   string
	baseCommand []string
	items       []string
}

type Store struct {
	Id                  int
	calcPanelRectHeight int
	calcPanelRectWidth  int
	shopNameRect        *util.Rect
	baseCommandRect     *util.Rect
	itemsRect           *util.Rect
	itemDescriptionRect *util.Rect
	ownerMessageRect    *util.Rect
	answerRect          *util.Rect
	StoreInfo           *StoreInfo
	font                *text.GoTextFace
	background          color.Color
}

func (s *Store) Init() error {
	s.background = color.RGBA{0x0, 0x0, 0x0, 0x5f}
	faceSouorce, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		slog.Error("Store.Init",
			slog.String("TextFaceSource loading error", err.Error()),
		)
	}
	japaneseFaceSource = faceSouorce
	return nil
}

func (s *Store) Update(data *gamestatus.GameData) {
	// Panel表示領域再計算
	if s.calcPanelRectHeight != data.LayoutHeight || s.calcPanelRectWidth != data.LayoutWidth {
		slog.Info("TalkPanel.Update()",
			slog.String("recalculate", "layout"),
		)
		baseHeight := int(float32(data.LayoutHeight) * 0.05)
		panelHeight := baseHeight
		panelY := 5
		panelWidth := int(float32(data.LayoutWidth) * 0.3)
		panelX := 5
		s.shopNameRect = &util.Rect{
			Left:   panelX,
			Top:    panelY,
			Right:  panelX + panelWidth,
			Bottom: panelY + panelHeight,
		}

		// ベースコマンド
		panelHeight = baseHeight * 4
		panelY = 5 + s.shopNameRect.Bottom
		panelWidth = int(float32(data.LayoutWidth) * 0.25)
		panelX = s.shopNameRect.Left
		s.baseCommandRect = &util.Rect{
			Left:   panelX,
			Top:    panelY,
			Right:  panelX + panelWidth,
			Bottom: panelY + panelHeight,
		}

		// アイテム
		panelHeight = int(float32(data.LayoutHeight) * 0.4)
		panelY = s.baseCommandRect.Top
		panelWidth = int(float32(data.LayoutWidth) * 0.6)
		panelX = s.baseCommandRect.Right + 5
		s.itemsRect = &util.Rect{
			Left:   panelX,
			Top:    panelY,
			Right:  panelX + panelWidth,
			Bottom: panelY + panelHeight,
		}

		s.calcPanelRectHeight = data.LayoutHeight
		s.calcPanelRectWidth = data.LayoutWidth

		// TODO 高さに応じてフォントサイズを変更
		s.font = &text.GoTextFace{
			Source:    japaneseFaceSource,
			Direction: text.DirectionLeftToRight,
			Size:      12,
			Language:  language.Japanese,
		}
	}
}

func (s *Store) Draw(screen *ebiten.Image, data *gamestatus.GameData) {
	if s.calcPanelRectHeight != data.LayoutHeight || s.calcPanelRectWidth != data.LayoutWidth {
		return
	}
	const lineSpacing = 48

	// 店名(左上)
	rect := s.shopNameRect
	vector.DrawFilledRect(screen, float32(rect.Left), float32(rect.Top), float32(rect.Width()), float32(rect.Height()), s.background, false)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(rect.Left+5), float64(rect.Top+5))
	op.LineSpacing = lineSpacing
	text.Draw(screen, s.StoreInfo.showpName, s.font, op)

	// ベースコマンド(店名の下長め)
	vector.DrawFilledRect(screen, float32(s.baseCommandRect.Left), float32(s.baseCommandRect.Top), float32(s.baseCommandRect.Width()), float32(s.baseCommandRect.Height()), s.background, false)
	for i, baseCommandText := range s.StoreInfo.baseCommand {
		op = &text.DrawOptions{}
		op.GeoM.Translate(float64(s.baseCommandRect.Left+12), float64(s.baseCommandRect.Top+5+i*25))
		op.LineSpacing = lineSpacing
		text.Draw(screen, baseCommandText, s.font, op)
	}

	// 買うもの(ベースコマンドの右重ね)
	vector.DrawFilledRect(screen, float32(s.itemsRect.Left), float32(s.itemsRect.Top), float32(s.itemsRect.Width()), float32(s.itemsRect.Height()), s.background, false)
	for i, itemsText := range s.StoreInfo.items {
		op = &text.DrawOptions{}
		op.GeoM.Translate(float64(s.itemsRect.Left+12), float64(s.itemsRect.Top+5+i*25))
		op.LineSpacing = lineSpacing
		text.Draw(screen, itemsText, s.font, op)
	}

	// TODO 詳細(中央)
	// TODO 店長の会話
	// TODO 回答 y/n
}
