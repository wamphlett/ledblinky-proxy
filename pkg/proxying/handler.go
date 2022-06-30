package proxying

import (
	"errors"
	"fmt"
	"net/rpc"
)

// InboundHandler defines the RPC handler
type InboundHandler struct {
	proxy *Proxy
}

// NewInboundHandler creates a new InboundHandler with the required dependencies
func NewInboundHandler(proxy *Proxy) (*InboundHandler, error) {
	h := &InboundHandler{proxy}
	if err := rpc.Register(h); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to register proxy handler: %s", err.Error()))
	}
	rpc.HandleHTTP()
	return h, nil
}

// Handle takes arguments and passes them off to the already running proxy handler
func (h *InboundHandler) Handle(args []string, reply *string) error {
	h.proxy.Handle(args)
	status := "OK"
	reply = &status
	return nil
}
