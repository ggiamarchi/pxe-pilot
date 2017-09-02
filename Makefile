.PHONY: all dep check-fmt check-vet check-lint check build clean

EXECUTABLE := pxe-pilot

all: build

dep:
	@go get github.com/Sirupsen/logrus
	@go get github.com/jawher/mow.cli
	@go get gopkg.in/gin-gonic/gin.v1
	@go get gopkg.in/yaml.v2
	@go get github.com/olekukonko/tablewriter

dep-dev:
	@go get github.com/golang/lint/golint

check-fmt:
	@! gofmt -d -e . | read

check-vet:
	@go vet

check-lint:
	@golint .

check: check-fmt check-vet check-lint

build: dep check
	@go build -o $(EXECUTABLE)

clean:
	@rm -f pxe-pilot
