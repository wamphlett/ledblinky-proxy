package publishing

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wamphlett/ledblinky-proxy/pkg/core/model"
)

// MockExecutor defines an executor which records the last execution
type MockExecutor struct {
	lastPath string
	lastArgs []string
}

// Run records the given path and arguments
func (e *MockExecutor) Run(path string, args []string) error {
	e.lastPath = path
	e.lastArgs = args
	return nil
}

// LastCommandString returns the last execution as a string
func (e *MockExecutor) LastCommandString() string {
	return fmt.Sprintf("%s \"%s\"", e.lastPath, strings.Join(e.lastArgs, "\" \""))
}

func TestEXEPublisherConstructionWithBadConfig(t *testing.T) {
	// should throw an error when trying to create a publisher with bad config
	publisher, err := NewEXEPublisher("")
	assert.Nil(t, publisher)
	assert.NotNil(t, err)
}

func TestEXEPublisherIssuesCorrectCommands(t *testing.T) {
	tt := map[string]struct {
		expectedCommandString string
		inputEvent            *model.Event
	}{
		"event with game and platform": {
			expectedCommandString: "./app.exe \"GAME_START\" \"Some Game\" \"PLATFORM_NAME\"",
			inputEvent: &model.Event{
				Type:     model.EVENT_TYPE_GAME_START,
				Game:     "Some Game",
				Platform: "PLATFORM_NAME",
			},
		},
		"event without game or platform": {
			expectedCommandString: "./app.exe \"GAME_QUIT\"",
			inputEvent: &model.Event{
				Type: model.EVENT_TYPE_GAME_QUIT,
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			mockExecutor := &MockExecutor{}
			publisher, err := NewEXEPublisherWithExecutor("./app.exe", mockExecutor)
			assert.Nil(t, err, "expected no error but got: %s", err)
			publisher.Publish(tc.inputEvent)
			assert.Equal(t, tc.expectedCommandString, mockExecutor.LastCommandString())
		})
	}
}
