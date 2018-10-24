package apis

import (
	"fmt"

	"github.com/cppforlife/go-patch/patch"
	"github.com/Microsoft/kunlun/common/fileio"
	"gopkg.in/yaml.v2"
)

type OpsFileArg struct {
	fileReader fileio.FileReader
	Ops        patch.Ops
}

func (a *OpsFileArg) UnmarshalFlag(filePath string) error {
	if len(filePath) == 0 {
		return fmt.Errorf("file path should not be empty")
	}

	bytes, err := a.fileReader.ReadFile(filePath)
	if err != nil {
		return err
	}

	var opDefs []patch.OpDefinition

	err = yaml.Unmarshal(bytes, &opDefs)
	if err != nil {
		return err
	}

	ops, err := patch.NewOpsFromDefinitions(opDefs)
	if err != nil {
		return err
	}

	(*a).Ops = ops

	return nil
}
