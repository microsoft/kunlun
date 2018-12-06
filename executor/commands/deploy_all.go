package commands

import (
	"github.com/Microsoft/kunlun/common/fileio"
	"github.com/Microsoft/kunlun/common/helpers"
	"github.com/Microsoft/kunlun/common/storage"
	"github.com/Microsoft/kunlun/common/ui"
)

type DeployAll struct {
	stateStore   storage.Store
	envIDManager helpers.EnvIDManager
	fs           fileio.Fs
	ui           *ui.UI
}

func NewDeployAll(
	stateStore storage.Store,
	envIDManager helpers.EnvIDManager,
	fs fileio.Fs,
	ui *ui.UI,
) DeployAll {
	return DeployAll{
		stateStore:   stateStore,
		envIDManager: envIDManager,
		fs:           fs,
		ui:           ui,
	}
}

func (d DeployAll) CheckFastFails(args []string, state storage.State) error {
	return nil
}

func (d DeployAll) Usage() string {
	return ""
}

func (d DeployAll) Execute(args []string, state storage.State) error {
	err := NewDigest(d.stateStore, d.envIDManager, d.fs, d.ui).Execute(args, state)
	if err != nil {
		return err
	}

	err = NewPlanInfra(d.stateStore, d.fs, d.ui).Execute(args, state)
	if err != nil {
		return err
	}

	err = NewApplyInfra(d.stateStore, d.fs).Execute(args, state)
	if err != nil {
		return err
	}

	err = NewPlanDeployment(d.stateStore, d.fs, d.ui).Execute(args, state)
	if err != nil {
		return err
	}

	err = NewApplyDeployment(d.stateStore).Execute(args, state)
	if err != nil {
		return err
	}

	return nil
}
