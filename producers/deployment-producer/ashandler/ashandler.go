package apis

import (
	"github.com/Microsoft/kunlun/artifacts/deployments"
	"github.com/Microsoft/kunlun/common/fileio"
	"github.com/Microsoft/kunlun/common/storage"
	"github.com/Microsoft/kunlun/common/ui"
	"github.com/Microsoft/kunlun/producers/deployment-producer/ashandler/generator"
)

type ASHandler struct {
	asGenerator generator.ASGenerator
}

func NewASHandler(
	stateStore storage.Store,
	ui *ui.UI,
	fs fileio.Fs,
) ASHandler {
	return ASHandler{
		asGenerator: generator.NewASGenerator(stateStore, ui, fs),
	}
}
func (a ASHandler) Handle(hostGroups []deployments.HostGroup, deployments []deployments.Deployment) error {
	return a.asGenerator.Generate(hostGroups, deployments)
}
