package scenes

import (
	"bytes"
	MiyatamaImages "first_rpg/miyatama/assets/images"
	gamestatus "first_rpg/miyatama/game_status"
	"image"
	_ "image/png"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

type UserDirection int

const (
	DIRECTION_DOWN UserDirection = iota
	DIRECTION_UP
	DIRECTION_LEFT
	DIRECTION_RIGHT
)

const (
	PLAYER_IMG_WIDTH  = 32
	PLAYER_IMG_HEIGTH = 32
	ANIMATION_SPAN    = 15
)

type Player struct {
	userDirection UserDirection
	frame         int
	playerImages  [][]*ebiten.Image
}

func (p *Player) Init() error {
	slog.Info("Player.Init()")
	img, _, err := image.Decode(bytes.NewReader(MiyatamaImages.Player_png))
	if err != nil {
		slog.Error("player.png decode error")
		slog.String("error: {}", err.Error())
		return err
	}
	playerImage := ebiten.NewImageFromImage(img)
	p.playerImages = [][]*ebiten.Image{}
	for i := 0; i < 4; i++ {
		progressImages := []*ebiten.Image{}
		for j := 0; j < 3; j++ {
			left := j * PLAYER_IMG_WIDTH
			top := i * PLAYER_IMG_HEIGTH
			right := left + PLAYER_IMG_WIDTH - 1
			bottom := top + PLAYER_IMG_HEIGTH - 1
			progressImages = append(progressImages, playerImage.SubImage(image.Rect(left, top, right, bottom)).(*ebiten.Image))
		}
		p.playerImages = append(p.playerImages, progressImages)
	}
	return nil
}

func (p *Player) Update(data *gamestatus.GameData) {
	switch data.UserAction {
	case gamestatus.USER_ACTION_LEFT:
		p.userDirection = DIRECTION_LEFT
	case gamestatus.USER_ACTION_UP:
		p.userDirection = DIRECTION_UP
	case gamestatus.USER_ACTION_RIGHT:
		p.userDirection = DIRECTION_RIGHT
	case gamestatus.USER_ACTION_DOWN:
		p.userDirection = DIRECTION_DOWN
	}
	// 状態の変更
	p.frame++
	if p.frame > ANIMATION_SPAN*3 {
		p.frame = 0
	}
}

func (p *Player) Draw(screen *ebiten.Image, data *gamestatus.GameData) {
	progressIndex := p.frame / ANIMATION_SPAN
	if progressIndex >= 3 {
		progressIndex = 2
	}

	directionIndex := 0
	switch p.userDirection {
	case DIRECTION_UP:
		directionIndex = 3
	case DIRECTION_LEFT:
		directionIndex = 1
	case DIRECTION_RIGHT:
		directionIndex = 2
	}
	// プレイヤーは常に画面中央
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(PLAYER_IMG_WIDTH)/2, -float64(PLAYER_IMG_HEIGTH)/2)
	op.GeoM.Translate(float64(data.LayoutWidth)/2, float64(data.LayoutHeight)/2)
	screen.DrawImage(p.playerImages[directionIndex][progressIndex], op)
}
