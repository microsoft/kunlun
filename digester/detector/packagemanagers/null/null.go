package null

import (
	"github.com/Microsoft/kunlun/digester/common"
)

type packageManager struct{}

func New() common.PackageManager {
	return &packageManager{}
}

func (p *packageManager) GetName() common.PackageManagerName {
	return ""
}

func (p *packageManager) Identify(path string) bool {
	return true
}

func (p *packageManager) DetectFramework(path string) []common.FrameworkName {
	return []common.FrameworkName{}
}
