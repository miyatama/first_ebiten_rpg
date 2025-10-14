package scenes

import (
	"bytes"
	"first_rpg/miyatama/assets/images"
	gamestatus "first_rpg/miyatama/game_status"
	"first_rpg/miyatama/util"
	"fmt"
	"image"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

type MobType int

const (
	MOB_TYPE_NONE MobType = iota
	MOB_TYPE_BLACK_CAT
	MOB_TYPE_VILLAGE_BOY
)

type MobCharacter struct {
	MobType         MobType
	Position        util.MapPosition
	Direction       util.Direction
	EventId         int
	frame           int
	characterImages [][]*ebiten.Image
	savedSx         int
	savedSy         int
	mapCorrectionX  int
	mapCorrectionY  int
}

func (m *MobCharacter) Init() error {
	imageBytes := m.getImageBytes()
	if len(imageBytes) > 0 {
		img, _, err := image.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			slog.Error("player.png decode error")
			slog.String("error: {}", err.Error())
			return err
		}
		characterImages := ebiten.NewImageFromImage(img)
		m.characterImages = [][]*ebiten.Image{}
		for i := 0; i < 4; i++ {
			progressImages := []*ebiten.Image{}
			for j := 0; j < 3; j++ {
				left := j * util.CHARACTER_IMG_WIDTH
				top := i * util.CHARACTER_IMG_HEIGTH
				right := left + util.CHARACTER_IMG_WIDTH - 1
				bottom := top + util.CHARACTER_IMG_HEIGTH - 1
				progressImages = append(progressImages, characterImages.SubImage(image.Rect(left, top, right, bottom)).(*ebiten.Image))
			}
			m.characterImages = append(m.characterImages, progressImages)
		}
	}
	return nil
}

func (m *MobCharacter) Update(data *gamestatus.GameData) {
	m.frame++
	if m.frame > util.CHARACTER_ANIMATION_SPAN*3 {
		m.frame = 0
	}
}

func (m *MobCharacter) SetDrawCorrection(sx, sy int) {
	m.mapCorrectionX = sx
	m.mapCorrectionY = sy
}

func (m *MobCharacter) Draw(screen *ebiten.Image, data *gamestatus.GameData) {
	if m.MobType == MOB_TYPE_NONE {
		return
	}
	progressIndex := m.frame / util.CHARACTER_ANIMATION_SPAN
	if progressIndex >= 3 {
		progressIndex = 2
	}

	directionIndex := 0
	switch m.Direction {
	case util.DIRECTION_UP:
		directionIndex = 3
	case util.DIRECTION_LEFT:
		directionIndex = 1
	case util.DIRECTION_RIGHT:
		directionIndex = 2
	}

	op := &ebiten.DrawImageOptions{}
	sx := m.Position.X * images.MAP_TILE_WIDTH
	sy := m.Position.Y * images.MAP_TILE_WIDTH
	if sx != m.savedSx || sy != m.savedSy {
		m.savedSx = sx
		m.savedSy = sy
		slog.Info("MobCharacter.Draw()",
			slog.String("position", fmt.Sprintf("{x: %d, y: %d}", m.Position.X, m.Position.Y)),
			slog.String("map position", fmt.Sprintf("{x: %d, y: %d}", sx, sy)),
		)
	}

	op.GeoM.Translate(float64(sx), float64(sy))
	op.GeoM.Translate(-float64(m.mapCorrectionX), -float64(m.mapCorrectionY))
	op.GeoM.Translate(float64(data.LayoutWidth)/2, float64(data.LayoutHeight)/2)
	screen.DrawImage(m.characterImages[directionIndex][progressIndex], op)
}

func (m *MobCharacter) getImageBytes() []byte {
	switch m.MobType {
	case MOB_TYPE_BLACK_CAT:
		return images.BlackCat_png
	case MOB_TYPE_VILLAGE_BOY:
		return images.VillageBoy_png
	}
	return []byte{}
}
