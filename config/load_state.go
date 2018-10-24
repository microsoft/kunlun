package config

import (
	"os"
	"path/filepath"

	flags "github.com/jessevdk/go-flags"
	"github.com/Microsoft/kunlun/common/configuration"
	"github.com/Microsoft/kunlun/common/fileio"
	"github.com/Microsoft/kunlun/common/storage"
)

type logger interface {
	Println(string)
}

type merger interface {
	MergeGlobalFlagsToState(globalflags GlobalFlags, state storage.State) (storage.State, error)
}

type fs interface {
	fileio.Stater
	fileio.TempFiler
	fileio.FileReader
	fileio.FileWriter
}

func NewConfig(bootstrap storage.StateBootstrap, merger merger, logger logger, fs fs) Config {
	return Config{
		stateBootstrap: bootstrap,
		merger:         merger,
		logger:         logger,
		fs:             fs,
	}
}

type Config struct {
	stateBootstrap storage.StateBootstrap
	merger         merger
	logger         logger
	fs             fs
}

func ParseArgs(args []string) (GlobalFlags, []string, error) {
	var globals GlobalFlags
	parser := flags.NewParser(&globals, flags.IgnoreUnknown)
	remainingArgs, err := parser.ParseArgs(args[1:])
	if err != nil {
		return GlobalFlags{}, remainingArgs, err
	}

	if !filepath.IsAbs(globals.StateDir) {
		workingDir, err := os.Getwd()
		if err != nil {
			return GlobalFlags{}, remainingArgs, err
		}
		globals.StateDir = filepath.Join(workingDir, globals.StateDir)
	}

	return globals, remainingArgs, nil
}

func (c Config) Bootstrap(globalFlags GlobalFlags, remainingArgs []string, argsLen int) (configuration.Configuration, error) {
	if argsLen == 1 { // if run kid.
		return configuration.Configuration{
			Command: "help",
		}, nil
	}

	var command string
	if len(remainingArgs) > 0 {
		command = remainingArgs[0]
	}

	if globalFlags.Version || command == "version" {
		command = "version"
		return configuration.Configuration{
			ShowCommandHelp: globalFlags.Help,
			Command:         command,
		}, nil
	}

	if len(remainingArgs) == 0 {
		return configuration.Configuration{
			Command: "help",
		}, nil
	}

	if len(remainingArgs) == 1 && command == "help" {
		return configuration.Configuration{
			Command: command,
		}, nil
	}

	if command == "help" {
		return configuration.Configuration{
			ShowCommandHelp: true,
			Command:         remainingArgs[1],
		}, nil
	}

	if globalFlags.Help {
		return configuration.Configuration{
			ShowCommandHelp: true,
			Command:         command,
		}, nil
	}

	state, err := c.stateBootstrap.GetState(globalFlags.StateDir)
	if err != nil {
		return configuration.Configuration{}, err
	}

	state, err = c.merger.MergeGlobalFlagsToState(globalFlags, state)
	if err != nil {
		return configuration.Configuration{}, err
	}

	return configuration.Configuration{
		Global: configuration.GlobalConfiguration{
			Debug:    globalFlags.Debug,
			StateDir: globalFlags.StateDir,
			Name:     globalFlags.EnvID,
		},
		State:                state,
		Command:              command,
		SubcommandFlags:      remainingArgs[1:],
		ShowCommandHelp:      false,
		CommandModifiesState: modifiesState(command),
	}, nil
}

func modifiesState(command string) bool {
	_, ok := map[string]struct{}{
		"analyze":    {}, // detect the project type and generate the draft manifests.
		"plan-lift":  {}, // parse the draft manifests and generate the infrastructure manifests. (now in terraform)
		"lift":       {}, // run the infra manifests. prepare the environment.
		"plan-shift": {}, // generate the deployment scripts, (now in ansible)
		"shift":      {}, // run the deployment scripts
		"destroy":    {}, // destroy the environment we just setup.
	}[command]
	return ok
}
