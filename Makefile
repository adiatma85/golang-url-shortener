.PHONY: swag-install
swag-install:
	@go install github.com/swaggo/swag/cmd/swag@v1.6.7

.PHONY: swaggo
swaggo:
	@/bin/rm -rf ./docs/swagger
	@`go env GOPATH`/bin/swag init -g ./src/cmd/main.go -o ./docs/swagger --parseInternal	

.PHONY: prepare
prepare: swag-install swaggo
	@go mod download

.PHONY: build
build:
	@go build -o ./build/app ./src/cmd

.PHONY: build-alpine
build-alpine:
	@go mod tidy && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/app ./src/cmd

.PHONY: run
run: swaggo build
	@./build/app
