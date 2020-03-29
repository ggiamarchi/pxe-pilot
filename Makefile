.PHONY: all check-fmt check-vet check-lint check build clean

EXECUTABLE := pxe-pilot

all: build

dep-dev:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.24.0

check-fmt:
	@! gofmt -d -e . | read

check-vet:
	@go vet

check-lint:
	@bin/golangci-lint run

check: check-fmt check-vet check-lint

build: check
	@go build -o $(EXECUTABLE)

clean:
	@rm -f pxe-pilot
