BINARY_NAME := bittwister

all: generate build

generate:
	go generate ./...

build:
	go build -o bin/$(BINARY_NAME) -v .

# `run` is used to ease the developer life
run: all
	sudo ./bin/$(BINARY_NAME) start -d wlp3s0 -b 500

test-packetloss:
	@bash ./scripts/tests/packetloss.sh

test-bandwidth:
	@bash ./scripts/tests/bandwidth.sh

test: test-packetloss test-bandwidth

.PHONY: all generate build run test