BINARY_NAME=pr-checker-go

all: build

deps:
	go mod tidy
	go mod vendor

build:
	go build -mod vendor -o $(BINARY_NAME) main.go

install-local:
	go install .

lint:
	golangci-lint run
