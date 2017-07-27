.PHONY: all dep build clean

all: build

dep:
	@go get github.com/Sirupsen/logrus
	@go get github.com/jawher/mow.cli
	@go get gopkg.in/gin-gonic/gin.v1

build: dep
	@go build -o pxepilot

clean:
	@rm -f pxepilot
