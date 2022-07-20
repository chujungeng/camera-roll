GOCMD=go
GOTEST=$(GOCMD) test
BINARY_NAME=cameraroll
MODE=CAMERAROLL_MODE
VERSION?=0.0.0

.PHONY: run

run: ## run in development mode
	$(MODE)=dev $(GOCMD) run cmd/main.go

clean: 
	rm -rf ./bin
