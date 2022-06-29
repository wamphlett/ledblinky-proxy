package test

import (
	"github.com/wamphlett/ledblinky-proxy/pkg/intercepting"
	"github.com/wamphlett/ledblinky-proxy/pkg/proxying"
	"github.com/wamphlett/ledblinky-proxy/test/support/mocks"
)

type testVars struct {
	proxy     *proxying.Proxy
	publisher *mocks.MockPublisher
}

func setup() *testVars {
	proxy := proxying.New(intercepting.New(), "", 3000)
	publisher := &mocks.MockPublisher{}
	proxy.AddPublisher(publisher)
	return &testVars{
		proxy:     proxy,
		publisher: publisher,
	}
}
