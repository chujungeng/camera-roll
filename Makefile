GOCMD=go
GOTEST=$(GOCMD) test
BINARY_NAME=cameraroll
MODE=CAMERAROLL_MODE
VERSION?=0.0.0

all: build

.PHONY: run build

run: ## run in development mode
	$(MODE)=dev $(GOCMD) run cmd/main.go

build: ## build in prod mode
	mkdir -p bin
	cp config.json bin/
	cp -r migration bin/
	$(MODE)=prod $(GOCMD) build -ldflags="-s -w" -o bin/$(BINARY_NAME) cmd/main.go

clean: 
	rm -rf ./bin
