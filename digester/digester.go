package apis

import (
	"github.com/Microsoft/kunlun/common/storage"
	"github.com/Microsoft/kunlun/digester/common"
	"github.com/Microsoft/kunlun/digester/questionnaire"
)

func Run(state storage.State, filePath string) error {
	bp := questionnaire.Run(state, filePath)
	return bp.ExposeYaml(filePath)
}

func ImportBlueprintYaml(filePath string) (common.Blueprint, error) {
	return common.ImportBlueprintYaml(filePath)
}
