package main

import (
	"log"
	"os"

	"github.com/Microsoft/kunlun/common/storage"
	cui "github.com/Microsoft/kunlun/common/ui"
	"github.com/Microsoft/kunlun/config"
	executor "github.com/Microsoft/kunlun/executor"
	"github.com/Microsoft/kunlun/executor/commands"
	"github.com/spf13/afero"
)

var Version = "dev"

func main() {
	log.SetFlags(0)

	ui := cui.NewUI(os.Stdout, os.Stdin)
	stderrUI := cui.NewUI(os.Stderr, os.Stdin)
	stateBootstrap := storage.NewStateBootstrap(stderrUI, Version)

	globals, remainingArgs, err := config.ParseArgs(os.Args)
	if err != nil {
		log.Fatalf("\n\n%s\n", err)
	}

	if globals.NoConfirm {
		ui.NoConfirm()
	}

	// File IO
	fs := afero.NewOsFs()
	afs := &afero.Afero{Fs: fs}

	// Configuration

	stateStore := storage.NewStore(globals.StateDir, afs)
	stateMerger := config.NewMerger(afs)
	newConfig := config.NewConfig(stateBootstrap, stateMerger, stderrUI, afs)

	appConfig, err := newConfig.Bootstrap(globals, remainingArgs, len(os.Args))
	if err != nil {
		log.Fatalf("\n\n%s\n", err)
	}

	usage := commands.NewUsage(ui)

	app := executor.NewExecutor(appConfig, usage, ui, stateStore, afs)

	err = app.Run()
	if err != nil {
		log.Fatalf("\n\n%s\n", err)
	}
}
