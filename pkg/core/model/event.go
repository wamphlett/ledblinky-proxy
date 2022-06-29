package model

type EventType string

const (
	EVENT_TYPE_GAME_SELECT    EventType = "GAME_SELECT"
	EVENT_TYPE_GAME_START               = "GAME_START"
	EVENT_TYPE_GAME_QUIT                = "GAME_QUIT"
	EVENT_TYPE_GAME_PAUSE               = "GAME_PAUSE"
	EVENT_TYPE_GAME_UNPAUSE             = "GAME_UNPAUSE"
	EVENT_TYPE_FE_START                 = "FE_START"
	EVENT_TYPE_FE_QUIT                  = "FE_QUIT"
	EVENT_TYPE_FE_LIST_CHANGE           = "FE_LIST_CHANGE"
	EVENT_TYPE_FE_SS_START              = "FE_SS_START"
	EVENT_TYPE_FE_SS_STOP               = "FE_SS_STOP"
	EVENT_TYPE_UNKNOWN                  = "UNKNOWN"
)

type Event struct {
	Type     EventType
	Game     string
	Platform string
}
