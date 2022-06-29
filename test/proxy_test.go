package test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wamphlett/ledblinky-proxy/pkg/core/model"
)

const (
	LEDBLINKY_GAME_SELECT_EVENT = "9 game_name platform_name"
	LEDBLINKY_GAME_START_EVENT  = "3"
)

func TestProxyHandlesIndependentEventsCorrectly(t *testing.T) {
	tt := map[string]struct {
		incomingEvent          string
		expectedPublishedEvent *model.Event
	}{
		"Game Selected": {LEDBLINKY_GAME_SELECT_EVENT, &model.Event{
			Type:     model.EVENT_TYPE_GAME_SELECT,
			Game:     "game_name",
			Platform: "platform_name",
		}},
		"Game Started": {LEDBLINKY_GAME_START_EVENT, &model.Event{
			Type:     model.EVENT_TYPE_GAME_START,
			Game:     "",
			Platform: "",
		}},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			tv := setup()
			tv.proxy.Handle(strings.Split(tc.incomingEvent, " "))
			assert.NotNil(t, tv.publisher.GetLastEvent())
			assert.Equal(t, tc.expectedPublishedEvent, tv.publisher.GetLastEvent())
		})
	}
}

// LaunchBox does not send the game name or platform name to LEDBlinky when starting a game,
// instead, it first sends a "game selected" event when the user navigates to a game. When the game
// is started, it then sends a "game start" event. This means we have to internally record which game
// was last selected so that when we send the "game start" event, we can send the intended
// game and platform.
//
// This test runs through the expected behavior when starting a game through LaunchBox/BigBox
func TestProxyHandlesGameStartScenario(t *testing.T) {
	tv := setup()
	// Game selected
	tv.proxy.Handle(convertEventStringToArgs(LEDBLINKY_GAME_SELECT_EVENT))
	assert.Equal(t, &model.Event{
		Type:     model.EVENT_TYPE_GAME_SELECT,
		Game:     "game_name",
		Platform: "platform_name",
	}, tv.publisher.GetLastEvent())

	// Game start should also send the game/platform from the last "game selected" event
	tv.proxy.Handle(convertEventStringToArgs(LEDBLINKY_GAME_START_EVENT))
	assert.Equal(t, &model.Event{
		Type:     model.EVENT_TYPE_GAME_START,
		Game:     "game_name",
		Platform: "platform_name",
	}, tv.publisher.GetLastEvent())
}

func convertEventStringToArgs(eventString string) []string {
	return strings.Split(eventString, " ")
}
