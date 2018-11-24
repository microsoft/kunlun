package commands

import (
	"github.com/Microsoft/kunlun/common/fileio"
	"github.com/Microsoft/kunlun/common/storage"
	"github.com/Microsoft/kunlun/common/ui"
	patching "github.com/Microsoft/kunlun/patching"
	deploymentProducer "github.com/Microsoft/kunlun/producers/deployment-producer"
)

type PlanDeployment struct {
	stateStore storage.Store
	fs         fileio.Fs
	ui         *ui.UI
}

func NewPlanDeployment(
	stateStore storage.Store,
	fs fileio.Fs,
	ui *ui.UI,
) PlanDeployment {
	return PlanDeployment{
		stateStore: stateStore,
		fs:         fs,
		ui:         ui,
	}
}

func (p PlanDeployment) CheckFastFails(args []string, state storage.State) error {
	return nil
}

func (p PlanDeployment) Execute(args []string, state storage.State) error {
	// load the provisioned manifest
	patching := patching.NewPatching(p.stateStore, p.fs)
	manifest, err := patching.ProvisionManifest()
	if err != nil {
		return err
	}
	deploymentProducer := deploymentProducer.NewDeploymentProducer(p.stateStore, p.ui, p.fs)
	err = deploymentProducer.Produce(*manifest)
	return err
}
