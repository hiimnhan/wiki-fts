buildPath = bin/run

all: build run

build:
	go build -o ${buildPath} main.go

run:
	./bin/run
