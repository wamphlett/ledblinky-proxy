package publishing

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wamphlett/ledblinky-proxy/pkg/core/model"
)

// MockHTTPClient defines a HTTP Client which records the last request
type MockHTTPClient struct {
	lastRequest *http.Request
}

// Do records the given request
func (c *MockHTTPClient) Do(request *http.Request) (*http.Response, error) {
	c.lastRequest = request
	return &http.Response{}, nil
}

// LastRequest returns the last recorded request
func (c *MockHTTPClient) LastRequest() *http.Request {
	return c.lastRequest
}

func TestHTTPPublisherConstructionWithBadConfig(t *testing.T) {
	// should throw an error when trying to create a publisher with bad config
	publisher, err := NewHTTPPublisher("")
	assert.Nil(t, publisher)
	assert.NotNil(t, err)
}

func TestHTTPPublisherSendsCorrectPayload(t *testing.T) {
	tt := map[string]struct {
		expectedAddress string
		expectedPayload string
		inputEvent      *model.Event
	}{
		"event with game and platform": {
			expectedPayload: `{"event_type":"GAME_START","game_name":"Some Game","platform_name":"PLATFORM_NAME"}`,
			inputEvent: &model.Event{
				Type:     model.EVENT_TYPE_GAME_START,
				Game:     "Some Game",
				Platform: "PLATFORM_NAME",
			},
		},
		"event without game or platform": {
			expectedPayload: `{"event_type":"GAME_QUIT","game_name":"","platform_name":""}`,
			inputEvent: &model.Event{
				Type: model.EVENT_TYPE_GAME_QUIT,
			},
		},
		"event with game and platform with quotes": {
			expectedPayload: `{"event_type":"GAME_START","game_name":"Some \"Game","platform_name":"PLATFORM_NAME"}`,
			inputEvent: &model.Event{
				Type:     model.EVENT_TYPE_GAME_START,
				Game:     "Some \"Game",
				Platform: "PLATFORM_NAME",
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			mockClient := &MockHTTPClient{}
			publisher, err := NewHTTPPublisherWithClient("localhost:8000", mockClient)
			assert.Nil(t, err, "expected no error but got: %s", err)
			publisher.Publish(tc.inputEvent)

			lastRequest := mockClient.LastRequest()
			assert.NotNil(t, lastRequest, "expected request to be not nil")
			assert.Equal(t, "POST", lastRequest.Method)
			assert.Equal(t, "localhost:8000", lastRequest.URL.String())

			b, err := io.ReadAll(lastRequest.Body)
			assert.NoError(t, err, "failed to read request body")
			assert.Equal(t, tc.expectedPayload, string(b))
		})
	}
}
