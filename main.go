package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	logFile = "ledblinky-proxy.log"
)

func main() {
	logArgs(os.Args[1:])
}

func logArgs(args []string) {
	log := fmt.Sprintf("[%s] %s", time.Now().Format("2006-02-01 15:04:05"), strings.Join(os.Args[1:], " "))
	fmt.Println(log)

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer f.Close()
		_, _ = f.WriteString(fmt.Sprintf("%s\n", log))
	}
}
