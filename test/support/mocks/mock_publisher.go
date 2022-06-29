package mocks

import (
	"github.com/wamphlett/ledblinky-proxy/pkg/core/model"
)

type MockPublisher struct {
	lastEvent *model.Event
}

func (p *MockPublisher) Publish(event *model.Event) error {
	p.lastEvent = event
	return nil
}

func (p *MockPublisher) GetLastEvent() *model.Event {
	return p.lastEvent
}
