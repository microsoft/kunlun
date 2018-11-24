package commands

import (
	"fmt"
	"strings"

	"github.com/Microsoft/kunlun/common/storage"
	"github.com/Microsoft/kunlun/common/ui"
)

const (
	UsageHeader = `
Usage:
  bbl [GLOBAL OPTIONS] %s [OPTIONS]

Global Options:
  --help       [-h]        Prints usage. Use "kl [command] --help" for more information about a command
  --state-dir  [-s]        Directory containing the kl state                                            env:"KL_STATE_DIRECTORY"
  --debug      [-d]        Prints debugging output                                                       env:"KL_DEBUG"
  --version    [-v]        Prints version
  --no-confirm [-n]        No confirm
%s
`
	CommandUsage = `
[%s command options]
  %s`
)

const GlobalUsage = `
`

type Usage struct {
	ui *ui.UI
}

func NewUsage(ui *ui.UI) Usage {
	return Usage{
		ui: ui,
	}
}

func (u Usage) CheckFastFails(subcommandFlags []string, state storage.State) error {
	return nil
}

func (u Usage) Execute(subcommandFlags []string, state storage.State) error {
	u.Print()
	return nil
}

func (u Usage) Print() {
	content := fmt.Sprintf(UsageHeader, "COMMAND", GlobalUsage)
	u.ui.Println(strings.TrimLeft(content, "\n"))
}

func (u Usage) PrintCommandUsage(command, message string) {
	commandUsage := fmt.Sprintf(CommandUsage, command, message)
	content := fmt.Sprintf(UsageHeader, command, commandUsage)
	u.ui.Println(strings.TrimLeft(content, "\n"))
}
