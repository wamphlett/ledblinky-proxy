package proxying

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/wamphlett/ledblinky-proxy/config"
	"github.com/wamphlett/ledblinky-proxy/pkg/core/model"
	"github.com/wamphlett/ledblinky-proxy/pkg/publishing"
)

type Publisher interface {
	Publish(*model.Event) error
}

type Interceptor interface {
	Intercept([]string) *model.Event
}

type Proxy struct {
	interceptor   Interceptor
	publishers    []Publisher
	port          int64
	keepAlive     chan bool
	handlerWg     sync.WaitGroup
	ledBlinkyPath string
}

func New(interceptor Interceptor, ledBlinkyPath string, port int64) *Proxy {
	return &Proxy{
		interceptor:   interceptor,
		port:          port,
		keepAlive:     make(chan bool, 1),
		handlerWg:     sync.WaitGroup{},
		ledBlinkyPath: validateLEDBlinkyPath(ledBlinkyPath),
	}
}

func validateLEDBlinkyPath(ledBlinkyPath string) string {
	if ledBlinkyPath == "" {
		log.Print("warn. LEDBlinky path is not set. it is highly recommended to pass events to LEDBlinky to maintain existing functionality.")
		return ""
	}

	if !filepath.IsAbs(ledBlinkyPath) {
		ledBlinkyPath = filepath.Join(filepath.Dir(os.Args[0]), ledBlinkyPath)
	}
	// check the path actually exists
	if f, err := os.Stat(ledBlinkyPath); errors.Is(err, os.ErrNotExist) || f == nil {
		log.Print("warn. cannot find LEDBlinky at the given path. check the path set in your config file to maintain LEDBlinky functionality.")
		return ""
	}

	return ledBlinkyPath
}

func (p *Proxy) ConfigurePublishers(receivers *config.ReceiversConfig) {
	for _, exePath := range receivers.Executables {
		publisher, err := publishing.NewEXEPublisher(exePath)
		if err != nil {
			log.Printf("error creating executable publisher for (%s): %s\n", exePath, err.Error())
			continue
		}
		p.AddPublisher(publisher)
	}

	for _, address := range receivers.Webhooks {
		publisher, err := publishing.NewHTTPPublisher(address)
		if err != nil {
			log.Printf("error creating HTTP publisher for (%s): %s\n", address, err.Error())
			continue
		}
		p.AddPublisher(publisher)
	}
}

func (p *Proxy) AddPublisher(publisher Publisher) {
	p.publishers = append(p.publishers, publisher)
}

func (p *Proxy) Handle(args []string) {
	p.handlerWg.Add(1)
	defer p.handlerWg.Done()
	// pass the arguments straight through to LEDBlinky before doing anything else
	if err := p.CallLEDBlinkey(args); err != nil {
		// print any errors and continue
		log.Printf("failed to call LEDBlinky: %s\n", err.Error())
	}

	// build an event from the incoming arguments
	event := p.interceptor.Intercept(args)

	// asynchronously publish the event using all the configured publishers
	var wg sync.WaitGroup
	wg.Add(len(p.publishers))
	for _, p := range p.publishers {
		go func(p Publisher) {
			defer wg.Done()
			if err := p.Publish(event); err != nil {
				log.Printf("error publishing: %s\n", err.Error())
			}
		}(p)
	}
	wg.Wait()

	// if the event was a "frontend quit" event, kill the proxy too
	if event.Type == model.EVENT_TYPE_FE_QUIT {
		p.End()
	}
}

func (p *Proxy) CallLEDBlinkey(args []string) error {
	if p.ledBlinkyPath == "" {
		return nil
	}
	cmd := exec.Command(p.ledBlinkyPath, args...)
	if err := cmd.Run(); err != nil {
		log.Printf("failed to exectube LEDBlinky: %s", err.Error())
	}
	return nil
}

func (p *Proxy) Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", p.port))
	if err != nil {
		return errors.New(fmt.Sprintf("listen error: %s", err))
	}

	NewInboundHandler(p)

	log.Printf("listening for RPC messages on port %d", p.port)
	go http.Serve(l, nil)
	// keep alive until we get an end request
	<-p.keepAlive
	// wait for any in progress handles to be done
	p.handlerWg.Wait()

	log.Print("received end\nexiting")
	return nil
}

func (p *Proxy) End() {
	p.keepAlive <- true
}
