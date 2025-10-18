package gamestatus

import (
	"first_rpg/miyatama/assets/events"
	"first_rpg/miyatama/util"

	"github.com/hajimehoshi/ebiten/v2"
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

type GameData struct {
	UserAction      UserAction
	ScreenWidth     int
	ScreenHeight    int
	LayoutWidth     int
	LayoutHeight    int
	Event           *events.Event
	EventMessageSeq int
	TouchIds        []ebiten.TouchID
	TouchPositions  map[ebiten.TouchID]util.TouchPosition
}
