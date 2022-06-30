package publishing

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/wamphlett/ledblinky-proxy/pkg/core/model"
)

// HTTPClient defines methods required by a HTTP Client
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// HTTPPublisher is used to POST event information to a URL
type HTTPPublisher struct {
	address string
	client  HTTPClient
}

// payload defines the payload sent in HTTP requests
type payload struct {
	EventType    string `json:"event_type"`
	GameName     string `json:"game_name"`
	PlatformName string `json:"platform_name"`
}

// NewHTTPPublisher creates a new HTTPPublisher with the required dependencies
func NewHTTPPublisher(address string) (*HTTPPublisher, error) {
	return NewHTTPPublisherWithClient(address, &http.Client{})
}

// NewHTTPPublisherWithClient creates a new HTTPPublisher with the required dependencies and
// allows the HTTPClient to be defined
func NewHTTPPublisherWithClient(address string, client HTTPClient) (*HTTPPublisher, error) {
	if address == "" {
		return nil, errors.New("cannot create HTTP publisher with empty address")
	}

	return &HTTPPublisher{
		address: address,
		client:  client,
	}, nil
}

// Publish POSTs a JSON payload to the configured address.
//  example:
//    {
//      "event_type":    "[EVENT TYPE]",
//      "game_name":     "[GAME NAME]",
//      "platform_name": "[PLATFORM NAME]"
//    }
func (p *HTTPPublisher) Publish(event *model.Event) error {
	// build a payload from the given event
	payload, err := p.buildJsonPayload(event)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to marshal payload: %s", err.Error()))
	}

	// build a HTTP request
	request, err := http.NewRequest("POST", p.address, bytes.NewBuffer(payload))
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create request: %s", err.Error()))
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// execute the request
	response, err := p.client.Do(request)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to do request: %s", err.Error()))
	}
	defer func() {
		if response.Body != nil {
			response.Body.Close()
		}
	}()

	return nil
}

func (p *HTTPPublisher) buildJsonPayload(event *model.Event) ([]byte, error) {
	payload := &payload{
		EventType:    string(event.Type),
		GameName:     event.Game,
		PlatformName: event.Platform,
	}
	return json.Marshal(payload)
}
