package executor

import (
	"crypto/rand"
	"fmt"

	"github.com/Microsoft/kunlun/common/configuration"
	"github.com/Microsoft/kunlun/common/fileio"
	"github.com/Microsoft/kunlun/common/helpers"
	"github.com/Microsoft/kunlun/common/storage"
	"github.com/Microsoft/kunlun/common/ui"
	"github.com/Microsoft/kunlun/executor/commands"
)

type usage interface {
	Print()
	PrintCommandUsage(command, message string)
}

type Executor struct {
	commands      commands.CommandSet
	configuration configuration.Configuration
	usage         usage
	ui            *ui.UI
	fs            fileio.Fs
}

func NewExecutor(
	configuration configuration.Configuration,
	usage usage,
	ui *ui.UI,
	stateStore storage.Store,
	fs fileio.Fs,
) Executor {

	envIDGenerator := helpers.NewEnvIDManager(rand.Reader)

	commandSet := commands.CommandSet{}
	commandSet["help"] = commands.NewUsage(ui)
	commandSet["analyze"] = commands.NewDigest(stateStore, envIDGenerator, fs, ui)
	commandSet["interop"] = commands.NewInterop(stateStore)
	commandSet["plan_infra"] = commands.NewPlanInfra(stateStore, fs, ui)
	commandSet["apply_infra"] = commands.NewApplyInfra(stateStore, fs)
	commandSet["plan_deployment"] = commands.NewPlanDeployment(stateStore, fs, ui)
	commandSet["apply_deployment"] = commands.NewApplyDeployment(stateStore)
	commandSet["promote"] = commands.NewPromote(stateStore)

	commandSet["ssh"] = commands.NewSSH(stateStore, fs, ui)
	return Executor{
		commands:      commandSet,
		configuration: configuration,
		usage:         usage,
		fs:            fs,
	}
}

func (a Executor) Run() error {
	err := a.execute()
	return err
}

func (a Executor) getCommand(commandString string) (commands.Command, error) {
	command, ok := a.commands[commandString]
	if !ok {
		a.usage.Print()
		return nil, fmt.Errorf("unknown command: %s", commandString)
	}
	return command, nil
}

func (a Executor) execute() error {
	command, err := a.getCommand(a.configuration.Command)
	if err != nil {
		return err
	}

	if a.configuration.ShowCommandHelp {
		a.usage.PrintCommandUsage(a.configuration.Command, command.Usage())
		return nil
	}

	if a.configuration.Command == "help" && len(a.configuration.SubcommandFlags) != 0 {
		commandString := a.configuration.SubcommandFlags[0]
		command, err = a.getCommand(commandString)
		if err != nil {
			return err
		}
		a.usage.PrintCommandUsage(commandString, command.Usage())
		return nil
	}

	err = command.CheckFastFails(a.configuration.SubcommandFlags, a.configuration.State)
	if err != nil {
		return err
	}

	return command.Execute(a.configuration.SubcommandFlags, a.configuration.State)
}
