package main

import (
	"dev.splitted-desktop.com/horizon/pxe-pilot/logger"
)

func main() {
	logger.Init()
	setupCLI()
}
