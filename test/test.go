package test

import (
	"bitbucket.org/godinezj/solid/log"
	"github.com/gobuffalo/logger"
)

func SetupTest() {
	log.Log = logger.NewLogger("debug")
	log.Info("Logging setup")
}
