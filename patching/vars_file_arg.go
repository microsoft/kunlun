package apis

import (
	"fmt"

	"github.com/Microsoft/kunlun/common/fileio"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type VarsFileArg struct {
	fs fileio.Fs

	Vars StaticVariables
}

func (a *VarsFileArg) UnmarshalFlag(filePath string) error {
	if len(filePath) == 0 {
		return errors.New("Expected file path to be non-empty")
	}

	bytes, err := a.fs.ReadFile(filePath)
	if err != nil {
		return errors.New(err.Error() + fmt.Sprintf("Reading variables file '%s'", filePath))
	}

	var vars StaticVariables

	err = yaml.Unmarshal(bytes, &vars)
	if err != nil {
		return errors.New(err.Error() + fmt.Sprintf("Deserializing variables file '%s'", filePath))
	}

	(*a).Vars = vars

	return nil
}
