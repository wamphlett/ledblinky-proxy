package proxying

import (
	"errors"
	"fmt"
	"net/rpc"
)

type InboundHandler struct {
	proxy *Proxy
}

func NewInboundHandler(proxy *Proxy) (*InboundHandler, error) {
	h := &InboundHandler{proxy}
	if err := rpc.Register(h); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to register proxy handler: %s", err.Error()))
	}
	rpc.HandleHTTP()
	return h, nil
}

func (h *InboundHandler) Handle(args []string, reply *string) error {
	h.proxy.Handle(args)
	status := "OK"
	reply = &status
	return nil
}
