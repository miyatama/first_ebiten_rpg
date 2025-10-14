package gamestatus

import "first_rpg/miyatama/assets/events"

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
	UserAction   UserAction
	ScreenWidth  int
	ScreenHeight int
	LayoutWidth  int
	LayoutHeight int
	Event        *events.Event
}
