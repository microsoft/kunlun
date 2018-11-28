package deployments

import (
	"github.com/Microsoft/kunlun/artifacts"
	yaml "gopkg.in/yaml.v2"
)

type Deployment struct {
	HostGroupName string
	Vars          yaml.MapSlice
	Roles         []artifacts.Role
}
