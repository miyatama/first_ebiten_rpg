package scenes

import (
	"bytes"
	MiyatamaImages "first_rpg/miyatama/assets/images"
	gamestatus "first_rpg/miyatama/game_status"
	"first_rpg/miyatama/util"
	"image"
	_ "image/png"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	userDirection util.Direction
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
			left := j * util.CHARACTER_IMG_WIDTH
			top := i * util.CHARACTER_IMG_HEIGTH
			right := left + util.CHARACTER_IMG_WIDTH - 1
			bottom := top + util.CHARACTER_IMG_HEIGTH - 1
			progressImages = append(progressImages, playerImage.SubImage(image.Rect(left, top, right, bottom)).(*ebiten.Image))
		}
		p.playerImages = append(p.playerImages, progressImages)
	}
	return nil
}

func (p *Player) Update(data *gamestatus.GameData) {
	// 状態の変更
	p.frame++
	if p.frame > util.CHARACTER_ANIMATION_SPAN*3 {
		p.frame = 0
	}
}

func (p *Player) Draw(screen *ebiten.Image, data *gamestatus.GameData) {
	progressIndex := p.frame / util.CHARACTER_ANIMATION_SPAN
	if progressIndex >= 3 {
		progressIndex = 2
	}

	directionIndex := 0
	switch p.userDirection {
	case util.DIRECTION_UP:
		directionIndex = 3
	case util.DIRECTION_LEFT:
		directionIndex = 1
	case util.DIRECTION_RIGHT:
		directionIndex = 2
	}
	// プレイヤーは常に画面中央
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(util.CHARACTER_IMG_WIDTH)/2, -float64(util.CHARACTER_IMG_HEIGTH)/2)
	op.GeoM.Translate(float64(data.LayoutWidth)/2, float64(data.LayoutHeight)/2)
	screen.DrawImage(p.playerImages[directionIndex][progressIndex], op)
}

func (p *Player) SetUserAction(userAction gamestatus.UserAction) {
	switch userAction {
	case gamestatus.USER_ACTION_LEFT:
		p.userDirection = util.DIRECTION_LEFT
	case gamestatus.USER_ACTION_UP:
		p.userDirection = util.DIRECTION_UP
	case gamestatus.USER_ACTION_RIGHT:
		p.userDirection = util.DIRECTION_RIGHT
	case gamestatus.USER_ACTION_DOWN:
		p.userDirection = util.DIRECTION_DOWN
	}
}
