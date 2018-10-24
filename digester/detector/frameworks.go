package detector

import (
	"github.com/Microsoft/kunlun/digester/common"
	"github.com/Microsoft/kunlun/digester/detector/frameworks/laravel5"
)

const (
	UnknownFramework common.FrameworkName = "unknown"
)

var Laravel5 common.FrameworkName = laravel5.New().GetName()

func getFrameworks() map[string]common.Framework {
	return map[string]common.Framework{
		string(Laravel5): laravel5.New(),
	}
}
