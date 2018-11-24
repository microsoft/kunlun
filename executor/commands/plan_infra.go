package commands

import (
	"github.com/Microsoft/kunlun/common/fileio"
	"github.com/Microsoft/kunlun/common/storage"
	"github.com/Microsoft/kunlun/common/ui"
	patching "github.com/Microsoft/kunlun/patching"
	infra "github.com/Microsoft/kunlun/producers/infra-producer"
	"github.com/Microsoft/kunlun/producers/infra-producer/handler"
)

type PlanInfra struct {
	stateStore storage.Store
	fs         fileio.Fs
	ui         *ui.UI
}

func NewPlanInfra(
	stateStore storage.Store,
	fs fileio.Fs,
	ui *ui.UI,
) PlanInfra {
	return PlanInfra{
		stateStore: stateStore,
		fs:         fs,
		ui:         ui,
	}
}

func (p PlanInfra) CheckFastFails(args []string, state storage.State) error {
	return nil
}

func (p PlanInfra) Execute(args []string, state storage.State) error {
	handlerType := handler.TerraformHandlerType // should get from args
	debug := true
	infraProducer, _ := infra.NewInfraProducer(p.stateStore, handlerType, debug)

	// load the provisioned manifest
	patching := patching.NewPatching(p.stateStore, p.fs)
	manifest, err := patching.ProvisionManifest()
	if err != nil {
		return err
	}

	return infraProducer.Setup(*manifest, state)
}
