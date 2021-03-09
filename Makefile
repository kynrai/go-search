.PHONY: all build run

all: build run
	
build:
	@go build -o bin/search cmd/main.go 

run:
	@bin/search
