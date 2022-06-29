package publishing

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/wamphlett/ledblinky-proxy/pkg/core/model"
)

type Executor interface {
	Run(path string, args []string) error
}

type EXEPublisher struct {
	executablePath string
	executor       Executor
}

func NewEXEPublisher(executablePath string) (*EXEPublisher, error) {
	return NewEXEPublisherWithExecutor(executablePath, &RealExecutor{})
}

func NewEXEPublisherWithExecutor(executablePath string, executor Executor) (*EXEPublisher, error) {
	if executablePath == "" {
		return nil, errors.New("cannot create executor with empty path")
	}
	return &EXEPublisher{
		executablePath: executablePath,
		executor:       executor,
	}, nil
}

func (p *EXEPublisher) Publish(event *model.Event) error {
	args := []string{string(event.Type)}
	if event.Game != "" {
		args = append(args, event.Game)
	}
	if event.Platform != "" {
		args = append(args, event.Platform)
	}
	err := p.executor.Run(p.executablePath, args)
	if err != nil {
		return errors.New(fmt.Sprintf("executor failed to run: %s", err.Error()))
	}
	return nil
}

type RealExecutor struct {
}

func (e *RealExecutor) Run(path string, args []string) error {
	cmd := exec.Command(path, args...)
	return cmd.Run()
}
