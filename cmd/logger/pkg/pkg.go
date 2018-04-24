package pkg

import "github.com/irfansharif/log"

func Log(logger *log.Logger) {
	logger.Info("from pkg!")
}
