package gamestatus

type GameStateMsg int

const (
	GAME_STATE_MSG_NONE GameStateMsg = iota
	GAME_STATE_MSG_TITLE
	GAME_STATE_MSG_HOUSE
)
