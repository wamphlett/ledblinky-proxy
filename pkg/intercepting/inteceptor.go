package intercepting

import (
	"strconv"

	"github.com/wamphlett/ledblinky-proxy/pkg/core/model"
)

// Interceptor is responsible for taking events from LaunchBox which were intended for
// LEDBlinkey. It transforms these events into proxy events and holds information about
// the current game/platform so that it can enrich events which did not originally
// have this information.
type Interceptor struct {
	currentGame     string
	currentPlatform string
}

// New creates a new Interceptor
func New() *Interceptor {
	return &Interceptor{}
}

// Intercept takes the raw arguments which were intended for LEDBlinky and returns a new
// proxy event.
func (i *Interceptor) Intercept(rawArgs []string) *model.Event {
	eventType := i.convertEventType(rawArgs[0])
	// the game selected event contains information about about the game and
	// platform. Use this to enrich other event types which do not have this information
	// such as "game_start" which does not contain this information.
	if eventType == model.EVENT_TYPE_GAME_SELECT {
		i.currentGame = rawArgs[1]
		i.currentPlatform = rawArgs[2]
	}

	// whenever "game_quit" is called, we should remove the current game / platform
	if eventType == model.EVENT_TYPE_GAME_QUIT {
		i.currentGame = ""
		i.currentPlatform = ""
	}

	event := &model.Event{
		Type: eventType,
	}

	// Enrich specific events with the game and platform type
	if eventType == model.EVENT_TYPE_GAME_SELECT || eventType == model.EVENT_TYPE_GAME_START {
		event.Game = i.currentGame
		event.Platform = i.currentPlatform
	}

	return event
}

// convertEventType is a 1:1 mapping of the LEDBlinky event types to proxy types
func (i *Interceptor) convertEventType(eventTypeString string) model.EventType {
	// we expect all LEDBlinky event types to be integers
	eventInt, err := strconv.ParseInt(eventTypeString, 10, 64)
	if err != nil {
		return model.EVENT_TYPE_UNKNOWN
	}

	LEDBlinkyEventMap := map[int64]model.EventType{
		1:  model.EVENT_TYPE_FE_START,
		2:  model.EVENT_TYPE_FE_QUIT,
		3:  model.EVENT_TYPE_GAME_START,
		4:  model.EVENT_TYPE_GAME_QUIT,
		5:  model.EVENT_TYPE_FE_SS_START,
		6:  model.EVENT_TYPE_FE_SS_STOP,
		8:  model.EVENT_TYPE_FE_LIST_CHANGE,
		9:  model.EVENT_TYPE_GAME_SELECT,
		16: model.EVENT_TYPE_GAME_PAUSE,
		17: model.EVENT_TYPE_GAME_UNPAUSE,
	}

	if eventType, ok := LEDBlinkyEventMap[eventInt]; ok {
		return eventType
	}

	return model.EVENT_TYPE_UNKNOWN
}
