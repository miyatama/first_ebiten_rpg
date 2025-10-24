package scenes

import (
	"bytes"
	"first_rpg/miyatama/ui/assets/fonts"
	MiyatamaImages "first_rpg/miyatama/ui/assets/images"
	gamestatus "first_rpg/miyatama/ui/game_status"
	"first_rpg/miyatama/ui/util"
	"image"
	"image/color"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/text/language"
)

type STORE_STATUS int

const (
	STORE_STATUS_IDLE STORE_STATUS = iota
	// ベースコマンドを選択できる状態
	STORE_STATUS_OPEN
	// アイテムを選択できる状態
	STORE_STATUS_SELECT_ITEM
	// 詳細を表示している状態(yes/noが表示されている)
	STORE_STATUS_ITEM_DESCRIPTION
	STORE_STATUS_CLOSE
)

type StoreInfo struct {
	showpName              string
	baseCommand            []string
	items                  []string
	baseCommandIndex       int
	itemsIndex             int
	baseCommandDecideIndex int
	itemsDecideIndex       int
}

type Store struct {
	status              STORE_STATUS
	Id                  int
	calcPanelRectHeight int
	calcPanelRectWidth  int
	shopNameRect        *util.Rect[int]
	baseCommandRect     *util.Rect[int]
	itemsRect           *util.Rect[int]
	itemDescriptionRect *util.Rect[int]
	ownerMessageRect    *util.Rect[int]
	answerRect          *util.Rect[int]
	StoreInfo           *StoreInfo

	font            *text.GoTextFace
	decideTextColor color.RGBA

	background       color.Color
	choiceImage      *ebiten.Image
	choiceImageScale float64
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

	img, _, err := image.Decode(bytes.NewReader(MiyatamaImages.AllowRight))
	if err != nil {
		slog.Error("player.png decode error")
		slog.String("error: {}", err.Error())
		return err
	}
	s.choiceImage = ebiten.NewImageFromImage(img)
	s.decideTextColor.R = 0xea
	s.decideTextColor.G = 0xe4
	s.decideTextColor.B = 0x4e
	s.decideTextColor.A = 0xff
	return nil
}

func (s *Store) Update(data *gamestatus.GameData) {
	// Panel表示領域再計算
	if s.calcPanelRectHeight != data.LayoutHeight || s.calcPanelRectWidth != data.LayoutWidth {
		slog.Info("TalkPanel.Update()",
			slog.String("recalculate", "layout"),
		)
		// 店名
		baseHeight := int(data.TextSizeSmallRect.Height() + 4)
		panelHeight := baseHeight
		panelY := 5
		panelWidth := int(float32(data.LayoutWidth) * 0.3)
		panelX := 5
		s.shopNameRect = &util.Rect[int]{
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
		s.baseCommandRect = &util.Rect[int]{
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
		s.itemsRect = &util.Rect[int]{
			Left:   panelX,
			Top:    panelY,
			Right:  panelX + panelWidth,
			Bottom: panelY + panelHeight,
		}

		s.calcPanelRectHeight = data.LayoutHeight
		s.calcPanelRectWidth = data.LayoutWidth

		s.font = &text.GoTextFace{
			Source:    japaneseFaceSource,
			Direction: text.DirectionLeftToRight,
			Size:      data.TextSizeSmall,
			Language:  language.Japanese,
		}

		s.choiceImageScale = data.TextSizeSmallRect.Height() / float64(s.choiceImage.Bounds().Dx())
	}

	// 入力の処理
	switch s.status {
	case STORE_STATUS_OPEN:
		// 開店時: ベースコマンドの選択が可能
		switch data.UserAction {
		case gamestatus.USER_ACTION_DECIDE:
			// 何かしら選択されている
			if s.StoreInfo.baseCommandIndex != -1 {
				// ベースコマンド選択
				s.StoreInfo.baseCommandDecideIndex = s.StoreInfo.baseCommandIndex
				if s.StoreInfo.baseCommandDecideIndex+1 >= len(s.StoreInfo.baseCommand) {
					// 決定:店を出る -> 閉店
					s.StoreInfo.baseCommandDecideIndex = -1
					s.status = STORE_STATUS_CLOSE
				} else {
					// 決定:ベースコマンド選択 -> STORE_STATUS_SELECT_ITEM
					s.status = STORE_STATUS_SELECT_ITEM
					s.StoreInfo.itemsIndex = 0
				}
			}
		case gamestatus.USER_ACTION_UP:
			if s.StoreInfo.baseCommandIndex <= 0 {
				s.StoreInfo.baseCommandIndex = len(s.StoreInfo.baseCommand) - 1
			} else if s.StoreInfo.baseCommandIndex > 0 {
				s.StoreInfo.baseCommandIndex--
			}
		case gamestatus.USER_ACTION_DOWN:
			if s.StoreInfo.baseCommandIndex+1 >= len(s.StoreInfo.baseCommand) {
				s.StoreInfo.baseCommandIndex = 0
			} else if s.StoreInfo.baseCommandIndex+1 <= len(s.StoreInfo.baseCommand) {
				s.StoreInfo.baseCommandIndex++
			}
		}
	case STORE_STATUS_SELECT_ITEM:
		// アイテムの選択が可能
		// TODO 上下 -> アイテム選択
		// TODO 決定:アイテム選択 -> STORE_STATUS_ITEM_DESCRIPTION
		switch data.UserAction {
		case gamestatus.USER_ACTION_DECIDE:
			// ベースコマンド選択
			s.StoreInfo.itemsDecideIndex = s.StoreInfo.itemsIndex
			// 決定:ベースコマンド選択 -> STORE_STATUS_SELECT_ITEM
			s.status = STORE_STATUS_ITEM_DESCRIPTION
		case gamestatus.USER_ACTION_UP:
			if s.StoreInfo.itemsIndex <= 0 {
				s.StoreInfo.itemsIndex = len(s.StoreInfo.items) - 1
			} else if s.StoreInfo.itemsIndex > 0 {
				s.StoreInfo.itemsIndex--
			}
		case gamestatus.USER_ACTION_DOWN:
			if s.StoreInfo.itemsIndex+1 >= len(s.StoreInfo.items) {
				s.StoreInfo.itemsIndex = 0
			} else if s.StoreInfo.itemsIndex+1 <= len(s.StoreInfo.items) {
				s.StoreInfo.itemsIndex++
			}
		case gamestatus.USER_ACTION_CANCEL:
			s.StoreInfo.itemsDecideIndex = -1
			s.StoreInfo.itemsIndex = -1
			s.status = STORE_STATUS_OPEN
		}
	case STORE_STATUS_ITEM_DESCRIPTION:
		// 詳細の選択が可能
	case STORE_STATUS_CLOSE:
		slog.Info("Store.Update()",
			slog.String("store status", "STORE_STATUS_CLOSE"),
		)

	}
}

func (s *Store) Draw(screen *ebiten.Image, data *gamestatus.GameData) {
	if s.calcPanelRectHeight != data.LayoutHeight || s.calcPanelRectWidth != data.LayoutWidth {
		return
	}
	padding := 5
	textHeight := data.TextSizeSmallRect.Height()
	lineSpacing := data.TextSizeSmall * 1.2

	// 店名(左上)
	rect := s.shopNameRect
	vector.DrawFilledRect(screen, float32(rect.Left), float32(rect.Top), float32(rect.Width()), float32(rect.Height()), s.background, false)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(rect.Left+5), float64(rect.Top+5))
	op.LineSpacing = lineSpacing
	text.Draw(screen, s.StoreInfo.showpName, s.font, op)

	// ベースコマンド(店名の下長め)
	vector.DrawFilledRect(screen, float32(s.baseCommandRect.Left), float32(s.baseCommandRect.Top), float32(s.baseCommandRect.Width()), float32(s.baseCommandRect.Height()), s.background, false)
	baseCommandLeft := s.baseCommandRect.Left + padding
	baseCommandTop := s.baseCommandRect.Top + padding
	cursorImageWidth := float64(s.choiceImage.Bounds().Dx()) * s.choiceImageScale
	for i, baseCommandText := range s.StoreInfo.baseCommand {
		// カーソルの描画
		if i == s.StoreInfo.baseCommandIndex {
			choiceImageOp := &ebiten.DrawImageOptions{}
			choiceImageOp.GeoM.Translate(-float64(s.choiceImage.Bounds().Dx())/2, -float64(s.choiceImage.Bounds().Dy())/2)
			choiceImageOp.GeoM.Scale(s.choiceImageScale, s.choiceImageScale)
			choiceImageOp.GeoM.Translate(
				float64(baseCommandLeft)+cursorImageWidth/2,
				float64(baseCommandTop)+cursorImageWidth/2+float64(i)*textHeight)
			screen.DrawImage(s.choiceImage, choiceImageOp)
		}

		op = &text.DrawOptions{}
		op.GeoM.Translate(
			float64(baseCommandLeft+int(cursorImageWidth)+5),
			float64(float64(baseCommandTop)+float64(i)*textHeight),
		)
		if i == s.StoreInfo.baseCommandDecideIndex {
			op.ColorScale.ScaleWithColor(s.decideTextColor)
		}
		op.LineSpacing = lineSpacing
		text.Draw(screen, baseCommandText, s.font, op)
	}

	// 買うもの(ベースコマンドの右重ね)
	if s.status == STORE_STATUS_SELECT_ITEM || s.status == STORE_STATUS_ITEM_DESCRIPTION {
		vector.DrawFilledRect(screen, float32(s.itemsRect.Left), float32(s.itemsRect.Top), float32(s.itemsRect.Width()), float32(s.itemsRect.Height()), s.background, false)
		for i, itemsText := range s.StoreInfo.items {
			// カーソルの描画
			if i == s.StoreInfo.itemsIndex {
				choiceImageOp := &ebiten.DrawImageOptions{}
				choiceImageOp.GeoM.Translate(-float64(s.choiceImage.Bounds().Dx())/2, -float64(s.choiceImage.Bounds().Dy())/2)
				choiceImageOp.GeoM.Scale(s.choiceImageScale, s.choiceImageScale)
				choiceImageOp.GeoM.Translate(
					float64(s.itemsRect.Left+padding)+cursorImageWidth/2,
					float64(s.itemsRect.Top+padding)+cursorImageWidth/2+float64(i)*textHeight)
				screen.DrawImage(s.choiceImage, choiceImageOp)
			}

			op = &text.DrawOptions{}
			op.GeoM.Translate(
				float64(s.itemsRect.Left+padding+int(cursorImageWidth)+5),
				float64(float64(s.itemsRect.Top+padding)+float64(i)*textHeight),
			)
			if i == s.StoreInfo.itemsDecideIndex {
				op.ColorScale.ScaleWithColor(s.decideTextColor)
			}
			op.LineSpacing = lineSpacing
			text.Draw(screen, itemsText, s.font, op)
		}
	}

	if s.status == STORE_STATUS_ITEM_DESCRIPTION {
		// 詳細(中央)
		// TODO 店長の会話
		// TODO 回答 y/n
	}
}

func (s *Store) Open() {
	s.status = STORE_STATUS_OPEN
}

func (s *Store) IsClosed() bool {
	return s.status == STORE_STATUS_CLOSE
}
