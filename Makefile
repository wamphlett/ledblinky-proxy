all: build build-logger

build:
	go build -o ledblinky-proxy .

build-logger:
	go build -o event-logger ./logger/
