BINARY_NAME=app
COVER_FILE=coverage.out

all: deps lint test build

deps:
	go mod tidy
	go mod vendor

bump-deps:
	go get -u ./...
	$(MAKE)

test:
	go test -v -count=1 -coverprofile $(COVER_FILE) -cover ./...

build:
	go build -mod vendor -a -o $(BINARY_NAME) .

run:
	go run .

install:
	go install .

lint:
	golangci-lint run