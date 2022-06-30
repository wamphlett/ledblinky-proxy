package publishing

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/wamphlett/ledblinky-proxy/pkg/core/model"
)

// Executor defines the methods required to execute a shell script
type Executor interface {
	Run(path string, args []string) error
}

// EXEPublisher is used to pass event information to another executable
type EXEPublisher struct {
	executablePath string
	executor       Executor
}

// NewEXEPublisher creates a new EXEPublisher with the required dependencies
func NewEXEPublisher(executablePath string) (*EXEPublisher, error) {
	return NewEXEPublisherWithExecutor(executablePath, &RealExecutor{})
}

// NewEXEPublisherWithExecutor creates a new EXEPublisher with the required
// dependencies and allows for the Executor to be specified
func NewEXEPublisherWithExecutor(executablePath string, executor Executor) (*EXEPublisher, error) {
	if executablePath == "" {
		return nil, errors.New("cannot create executor with empty path")
	}
	return &EXEPublisher{
		executablePath: executablePath,
		executor:       executor,
	}, nil
}

// Publish publishes an event the configured executable path.
//  example:
//    app.exe [EVENT TYPE] [GAME NAME] [PLATFORM NAME]
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

// RealExecutor defines an executor which uses the exe.Cmd library
type RealExecutor struct{}

// Run executes the given path and arguments
func (e *RealExecutor) Run(path string, args []string) error {
	cmd := exec.Command(path, args...)
	return cmd.Run()
}
