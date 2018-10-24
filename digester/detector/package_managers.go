package detector

import (
	"github.com/Microsoft/kunlun/digester/common"
	"github.com/Microsoft/kunlun/digester/detector/packagemanagers/composer"
)

const (
	UnknownPackageManager common.PackageManagerName = "unknown"
)

var Composer common.PackageManagerName = composer.New().GetName()

func getPackageManagers() map[string]common.PackageManager {
	return map[string]common.PackageManager{
		string(Composer): composer.New(),
	}
}
