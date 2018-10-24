package handler

import (
	"github.com/Microsoft/kunlun/common/helpers"
)

type Outputs interface {
	GetString(key string) string
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]string
}

var outputsOpsFile = []byte(`
---
- type: replace
  path: /vm_groups/name=jumpbox/networks/0/outputs?
  value:
    - public_ip: {{.vm_groups_jumpbox_networks_0_outputs_0}}
- type: replace
  path: /vm_groups/name=web-servers/networks/0/outputs?
  value:
    - ip: {{.vm_groups_web_servers_networks_0_outputs_0}}
`)

func ToOutputsOpsFile(outputs Outputs) (string, error) {
	contents, err := helpers.Render(outputsOpsFile, getOutputsParams(outputs))
	if err != nil {
		return "", err
	}

	return contents, nil
}

func getOutputsParams(outputs Outputs) map[string]interface{} {
	return map[string]interface{}{
		"vm_groups_jumpbox_networks_0_outputs_0":     outputs.GetStringMap("vm_groups_jumpbox_networks_0_outputs_0")["public_ip"],
		"vm_groups_web_servers_networks_0_outputs_0": outputs.GetStringMap("vm_groups_web-servers_networks_0_outputs_0")["ip"],
	}
}
