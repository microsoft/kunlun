package apis

import (
	"github.com/Microsoft/kunlun/artifacts"

	"github.com/Microsoft/kunlun/artifacts/deployments"
	"github.com/Microsoft/kunlun/common/fileio"
	"github.com/Microsoft/kunlun/common/storage"
	"github.com/Microsoft/kunlun/common/ui"
	ashandler "github.com/Microsoft/kunlun/producers/deployment-producer/ashandler"
	"github.com/Microsoft/kunlun/producers/deployment-producer/dpbuilder"
)

type DeploymentProducer struct {
	stateStore storage.Store
	ui         *ui.UI
	fs         fileio.Fs
}

func NewDeploymentProducer(
	stateStore storage.Store,
	ui *ui.UI,
	fs fileio.Fs,
) DeploymentProducer {
	return DeploymentProducer{
		stateStore: stateStore,
		ui:         ui,
		fs:         fs,
	}
}

type deploymentItem struct {
	hostGroup  deployments.HostGroup
	deployment deployments.Deployment
}

func (dp DeploymentProducer) Produce(
	manifest artifacts.Manifest,
) error {
	// generate the deployments
	dpBuilder := dpbuilder.DeploymentBuilder{}
	hostGroups, deployments, err := dpBuilder.Produce(manifest)
	if err != nil {
		return err
	}

	// generate the ansible scripts based on the deployments.
	asHandler := ashandler.NewASHandler(dp.stateStore, dp.ui, dp.fs)
	err = asHandler.Handle(hostGroups, deployments)
	if err != nil {
		return err
	}
	return nil
}
