GOCMD=go
GOTEST=$(GOCMD) test
BINARY_NAME=cameraroll
MODE=CAMERAROLL_MODE
VERSION?=0.0.0

all: build

.PHONY: dev build

dev: ## build and run in development mode
	mkdir -p bin
	cp config.json bin/
	cp .env bin/
	cp -r migration bin/
	$(GOCMD) build -o bin/$(BINARY_NAME) cmd/main.go
	$(MODE)=dev bin/$(BINARY_NAME)

build: ## build in prod mode
	mkdir -p bin
	cp config.json bin/
	cp -r migration bin/
	$(GOCMD) build -ldflags="-s -w" -o bin/$(BINARY_NAME) cmd/main.go

clean: 
	rm -f ./bin/$(BINARY_NAME)
	rm -f ./bin/.env
	rm -f ./bin/config.json
	rm -rf ./bin/migration
