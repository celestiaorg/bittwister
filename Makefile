BINARY_NAME := bittwister

all: generate build

generate:
	go generate ./...

build:
	go build -o bin/$(BINARY_NAME) -v .

docker:
	docker build -t bittwister .

lint:
	golangci-lint run ./...

test-go:
	sudo go test -v ./... -count=1 -p=1

test-packetloss:
	@bash ./scripts/tests/packetloss.sh

test-bandwidth:
	@bash ./scripts/tests/bandwidth.sh

test-latency:
	@bash ./scripts/tests/latency.sh

test-jitter:
	@bash ./scripts/tests/jitter.sh

test: test-go test-packetloss test-bandwidth test-latency test-jitter

.PHONY: all generate build run test