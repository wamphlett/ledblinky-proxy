package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	logFile = "event-logger.log"
)

func main() {
	logArgs(os.Args[1:])
}

func logArgs(args []string) {
	for i, a := range args {
		if strings.Contains(a, " ") {
			args[i] = fmt.Sprintf(`"%s"`, a)
		}
	}

	logMsg := fmt.Sprintf("[%s] %s", time.Now().Format("2006-02-01 15:04:05"), strings.Join(args, " "))
	log.Print(logMsg)

	f, err := os.OpenFile(filepath.Join(filepath.Dir(os.Args[0]), logFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer f.Close()
		_, _ = f.WriteString(fmt.Sprintf("%s\n", logMsg))
	}
}
