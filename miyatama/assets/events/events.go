package events

type EventType int

const (
	EVENT_TYPE_MOB_TALK EventType = iota
	EVENT_TYPE_STORE
)

type Event struct {
	Id          int
	EventType   EventType
	TalkTexts   []string
	NextEventId int
	StoreId     int
}
