package main

import (
	"github.com/jwdev42/logger"
	"github.com/jwdev42/xtagger/internal/global"
	"github.com/jwdev42/xtagger/internal/program"
	"os"
)

func main() {
	//Setup logger
	logger.SetupDefaultLogger(os.Stderr, logger.LevelError, " - ")
	//Run program
	if err := program.Run(); err != nil {
		logger.Default().Panicf("%s", err)
		global.ExitCode = global.ExitHardError
	}
	os.Exit(int(global.ExitCode))
}
