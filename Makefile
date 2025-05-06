CONFIG_PATH ?= sunny_5_skiers/config1.json
EVENTS ?= sunny_5_skiers/events1
BINARY = biathlon
COVER = cover
COVER_OUT = $(COVER).out
COVER_HTML = $(COVER).html

.PHONY: all build run test cover clean

all: build

build:
	@go build -o $(BINARY) ./...

run:
	@CONFIG_PATH=$(CONFIG_PATH) go run . < $(EVENTS)

test:
	@go test -v -race ./...

cover:
	@go test -covermode=atomic -coverprofile=$(COVER_OUT) -race ./...
	@go tool cover -func=$(COVER_OUT)
	@go tool cover -html=$(COVER_OUT) -o $(COVER_HTML)

clean:
	@go clean
	@rm -f $(COVER_OUT) $(COVER_HTML)