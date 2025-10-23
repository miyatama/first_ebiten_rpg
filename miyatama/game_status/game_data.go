package gamestatus

import (
	"first_rpg/miyatama/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type UserAction int

const (
	USER_ACTION_NONE UserAction = iota
	USER_ACTION_LEFT
	USER_ACTION_RIGHT
	USER_ACTION_UP
	USER_ACTION_DOWN
	USER_ACTION_DECIDE
)

const (
	USER_INPUT_WAIT_FRAME_COUNT = 20
)

type GameData struct {
	UserAction     UserAction
	ScreenWidth    int
	ScreenHeight   int
	LayoutWidth    int
	LayoutHeight   int
	TouchIds       []ebiten.TouchID
	TouchPositions map[ebiten.TouchID]util.TouchPosition
	GOOS           string
	Environemnt    *Environment
	// Audio
	AudioContext *audio.Context

	// TextSize
	TextSizeLarge  float64
	TextSizeMiddle float64
	TextSizeSmall  float64
}

func (g *GameData) IsMobile() bool {
	if g.GOOS == "android" || g.GOOS == "ios" {
		return true
	}
	return false
}
