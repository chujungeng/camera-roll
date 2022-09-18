GOCMD=go
GOTEST=$(GOCMD) test
BINARY_NAME=cameraroll
MODE=CAMERAROLL_MODE
VERSION?=0.0.0
COMMIT=$(shell git rev-list -1 HEAD)

all: client build

.PHONY: dev build client

client: ## build the frontend of admin-area
	npm run --prefix client build

dev: ## build and run in development mode
	mkdir -p bin
	cp config.json bin/
	cp .env bin/
	cp -r migration bin/
	cp -r client/build bin/client
	$(GOCMD) build -ldflags="-X main.commit=$(COMMIT)" -o bin/$(BINARY_NAME) . 
	$(MODE)=dev bin/$(BINARY_NAME)

build: ## build in prod mode
	mkdir -p bin
	cp config.json bin/
	cp -r migration bin/
	cp -r client/build bin/client
	$(GOCMD) build -ldflags="-s -w -X main.commit=$(COMMIT)" -o bin/$(BINARY_NAME) .

clean: 
	rm -f ./bin/$(BINARY_NAME)
	rm -f ./bin/.env
	rm -f ./bin/config.json
	rm -rf ./bin/migration
	rm -rf ./bin/client
	rm -rf ./client/build
