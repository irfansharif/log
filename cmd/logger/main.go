package main

import (
	"os"

	"github.com/irfansharif/log"
)

func main() {
	logger := log.New(os.Stdout)
	logger.Info("hello, world!")
	log.SetLogMode(log.DisabledMode)
	logger.Info("another hello, world!")
}