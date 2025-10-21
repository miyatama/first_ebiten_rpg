package gamestatus

type EventType int

const (
	EVENT_TYPE_MOB_TALK EventType = iota
	EVENT_TYPE_STORE
	EVENT_TYPE_CHANGE_SCENE
)

type EventTag int

const (
	EVENT_TAG_CHANGE_TO_HOUSE_SCENE EventTag = iota
)

type Event struct {
	Id          int
	EventType   EventType
	TalkTexts   []string
	NextEventId int
	StoreId     int
	EventTag    EventTag
}
