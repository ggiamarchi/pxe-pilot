.PHONY: all dep build clean

all: build

dep:
	@go get github.com/jawher/mow.cli

build: dep
	@go build -o pxepilot

clean:
	@rm -f pxepilot
