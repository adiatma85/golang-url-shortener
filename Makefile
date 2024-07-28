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

.PHONY: run-tests
run-tests:
	@go clean -cache
	@go test -v -failfast `go list ./... | grep -i 'business'` -cover

.PHONY: mock-install
mock-install:
	@go install go.uber.org/mock/mockgen@v0.4.0

.PHONY: mock
mock:
	@`go env GOPATH`/bin/mockgen -source src/business/domain/$(domain)/$(domain).go -destination src/business/domain/mock/$(domain)/$(domain).go

.PHONY: mock-all
mock-all:
	@make mock domain=user
	@make mock domain=url