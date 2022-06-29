package intercepting

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wamphlett/ledblinky-proxy/pkg/core/model"
)

func TestInterceptorRecordsCurrentGameAndPlatform(t *testing.T) {
	interceptor := &Interceptor{}
	// send the arguments LaunchBox would send for a "game select" event and make
	// sure the interceptor records the game and platform correctly
	_ = interceptor.Intercept([]string{"9", "first_game_name", "first_platform_name"})
	assert.Equal(t, "first_game_name", interceptor.currentGame)
	assert.Equal(t, "first_platform_name", interceptor.currentPlatform)

	// send the arguments LaunchBox would send for a "game start" event and make sure
	// the current game and platform remain untouched
	_ = interceptor.Intercept([]string{"3"})
	assert.Equal(t, "first_game_name", interceptor.currentGame)
	assert.Equal(t, "first_platform_name", interceptor.currentPlatform)

	// send the arguments LaunchBox would send for a "game quit" event and make sure
	// the current game and platform get cleared out
	_ = interceptor.Intercept([]string{"4"})
	assert.Equal(t, "", interceptor.currentGame)
	assert.Equal(t, "", interceptor.currentPlatform)

	// send the arguments LaunchBox would send for a "game select" event and make
	// sure the interceptor records the new game and platform correctly
	_ = interceptor.Intercept([]string{"9", "second_game_name", "second_platform_name"})
	assert.Equal(t, "second_game_name", interceptor.currentGame)
	assert.Equal(t, "second_platform_name", interceptor.currentPlatform)
}

func TestInterceptorEnrichesTheStartGameEvent(t *testing.T) {
	interceptor := &Interceptor{}
	// send the arguments LaunchBox would send for a "game select" event so that the
	// interceptor knows which game is currently selected
	_ = interceptor.Intercept([]string{"9", "My Game Name", "The_Platform"})

	// send the arguments LaunchBox would send for a "game start" event and make sure
	// the current game and platform are present on the new "game start" event
	event := interceptor.Intercept([]string{"3"})
	assert.Equal(t, &model.Event{
		Type:     model.EVENT_TYPE_GAME_START,
		Game:     "My Game Name",
		Platform: "The_Platform",
	}, event)
}

func TestInterceptorConvertsTypesCorrectly(t *testing.T) {
	tt := map[string]struct {
		inputEvent        string
		expectedEventType model.EventType
	}{
		// invalid events
		"non integer":    {"x", model.EVENT_TYPE_UNKNOWN},
		"unknown number": {"1000", model.EVENT_TYPE_UNKNOWN},
		// game events
		"game start":   {"3", model.EVENT_TYPE_GAME_START},
		"game quit":    {"4", model.EVENT_TYPE_GAME_QUIT},
		"game select":  {"9", model.EVENT_TYPE_GAME_SELECT},
		"game pause":   {"16", model.EVENT_TYPE_GAME_PAUSE},
		"game unpause": {"17", model.EVENT_TYPE_GAME_UNPAUSE},

		// frontend events
		"frontend start":             {"1", model.EVENT_TYPE_FE_START},
		"frontend quit":              {"2", model.EVENT_TYPE_FE_QUIT},
		"frontend list change":       {"8", model.EVENT_TYPE_FE_LIST_CHANGE},
		"frontend screensaver start": {"5", model.EVENT_TYPE_FE_SS_START},
		"frontend screensaver stop":  {"6", model.EVENT_TYPE_FE_SS_STOP},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			interceptor := &Interceptor{}
			assert.Equal(t, tc.expectedEventType, interceptor.convertEventType(tc.inputEvent))
		})
	}
}
