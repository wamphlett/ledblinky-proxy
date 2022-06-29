package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"

	"github.com/wamphlett/ledblinky-proxy/config"
	"github.com/wamphlett/ledblinky-proxy/pkg/intercepting"
	"github.com/wamphlett/ledblinky-proxy/pkg/proxying"
)

func main() {
	port := 8812

	client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", port))
	if err == nil {
		var reply string
		log.Print("handing off to existing proxy service")
		err = client.Call("InboundHandler.Handle", os.Args[1:], &reply)
		if err != nil {
			log.Fatal("inbound handler error:", err)
		}
		log.Print("arguments passed to proxy\nexiting")
		return
	}

	log.Printf("starting new proxy service on port %d", port)

	cfg, err := config.NewFromINI()
	if err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
	}
	interceptor := intercepting.New()
	proxy := proxying.New(interceptor, cfg.LEDBlinkyPath, int64(port))
	proxy.ConfigurePublishers(cfg.Receivers)

	// // TEMP MAKE AN EVENT LOGGER CALL
	// tempPublisher, err := publishing.NewEXEPublisher("./event-logger")
	// if err != nil {
	// 	log.Printf("failed to start temp publisher: %s", err.Error())
	// 	return
	// }
	// proxy.AddPublisher(tempPublisher)

	// // ADD A TEMP HTTP PUBLISHER
	// tempHTTPPublisher, err := publishing.NewHTTPPublisher("https://eop2yw3sg5mpebi.m.pipedream.net")
	// if err != nil {
	// 	log.Printf("failed to start temp http publisher: %s", err.Error())
	// 	return
	// }
	// proxy.AddPublisher(tempHTTPPublisher)

	// make sure the first set of args are handled
	proxy.Handle(os.Args[1:])

	// start the proxy and listen for further arguments
	if err := proxy.Start(); err != nil {
		log.Print("failed to start proxy service")
		return
	}
}
