package handler

import (
	"fmt"
	"github.com/Microsoft/kunlun/common/helpers"
	"strings"
)

type Outputs interface {
	Keys() []string
	GetString(key string) string
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]string
}

var outputIpTemplate = []byte(`
- type: replace
  path: /vm_groups/name={{.vm_group_name}}/networks/0/outputs?
  value:
    - ip: {{.ip}}
`)
var outputPublicIpTemplate = []byte(`
- type: replace
  path: /vm_groups/name={{.vm_group_name}}/networks/0/outputs?/0/public_ip?
  value: {{.public_ip}}
`)
var outputHostTemplate = []byte(`
- type: replace
  path: /vm_groups/name={{.vm_group_name}}/networks/0/outputs?/0/host?
  value: {{.host}}
`)
var outputOptionalTemplates = map[string][]byte{
	"public_ip": outputPublicIpTemplate,
	"host":      outputHostTemplate,
}

func ToOutputsOpsFile(outputs Outputs) (string, error) {
	vm_group_outputs := ""
	for _, key := range outputs.Keys() {
		if strings.HasPrefix(key, "vm_groups") {
			vm_group_name := strings.Split(key, "_")[2]
			template, params := getTemplateAndParams(vm_group_name, outputs)
			vm_group_output, err := helpers.Render(template, params)
			if err != nil {
				return "", err
			}
			vm_group_outputs += vm_group_output
		}
	}
	return vm_group_outputs, nil
}

func getTemplateAndParams(vm_group_name string, outputs Outputs) ([]byte, map[string]interface{}) {
	output := outputs.GetStringMap(fmt.Sprintf("vm_groups_%s_networks_0_outputs_0", vm_group_name))

	template := make([]byte, 0)
	params := make(map[string]interface{})
	params["vm_group_name"] = vm_group_name
	// We assume that the internal IP address always exists
	template = append(template, outputIpTemplate...)
	params["ip"] = output["ip"]
	// public_ip and host are optional
	param_keys := []string{"public_ip", "host"}
	for _, param_key := range param_keys {
		if value, ok := output[param_key]; ok {
			template = append(template, outputOptionalTemplates[param_key]...)
			params[param_key] = value
		}
	}
	return template, params
}
