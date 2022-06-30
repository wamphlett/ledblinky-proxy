package mocks

import (
	"github.com/wamphlett/ledblinky-proxy/pkg/core/model"
)

// MockPublisher records events which have been published
type MockPublisher struct {
	lastEvent *model.Event
}

// Publish records the published event
func (p *MockPublisher) Publish(event *model.Event) error {
	p.lastEvent = event
	return nil
}

// GetLastEvent returns the last recorded event
func (p *MockPublisher) GetLastEvent() *model.Event {
	return p.lastEvent
}
