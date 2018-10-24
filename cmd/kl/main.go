package main

import (
	"log"
	"os"

	clogger "github.com/Microsoft/kunlun/common/logger"
	"github.com/Microsoft/kunlun/common/storage"
	"github.com/Microsoft/kunlun/config"
	executor "github.com/Microsoft/kunlun/executor"
	"github.com/Microsoft/kunlun/executor/commands"
	"github.com/spf13/afero"
)

var Version = "dev"

func main() {
	log.SetFlags(0)

	logger := clogger.NewLogger(os.Stdout, os.Stdin)
	stderrLogger := clogger.NewLogger(os.Stderr, os.Stdin)
	stateBootstrap := storage.NewStateBootstrap(stderrLogger, Version)

	globals, remainingArgs, err := config.ParseArgs(os.Args)
	if err != nil {
		log.Fatalf("\n\n%s\n", err)
	}

	if globals.NoConfirm {
		logger.NoConfirm()
	}
	// stateJSON, _ := json.Marshal(globals)
	// stderrLogger.Println(string(stateJSON))

	// File IO
	fs := afero.NewOsFs()
	afs := &afero.Afero{Fs: fs}

	// Configuration

	stateStore := storage.NewStore(globals.StateDir, afs)
	stateMerger := config.NewMerger(afs)
	newConfig := config.NewConfig(stateBootstrap, stateMerger, stderrLogger, afs)

	appConfig, err := newConfig.Bootstrap(globals, remainingArgs, len(os.Args))
	if err != nil {
		log.Fatalf("\n\n%s\n", err)
	}

	// // Utilities
	// envIDGenerator := helpers.NewEnvIDGenerator(rand.Reader)
	usage := commands.NewUsage(logger)

	app := executor.NewExecutor(appConfig, usage, logger, stateStore, afs)

	err = app.Run()
	if err != nil {
		log.Fatalf("\n\n%s\n", err)
	}
}
