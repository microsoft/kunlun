package templates

import (
	"github.com/Microsoft/kunlun/common/helpers"
)

var resourceGroupTF = []byte(`
variable "resource_group_name" {}
variable "location" {}
variable "env_name" {}

resource "azurerm_resource_group" "kunlun_resource_group" {
	name     = "${var.resource_group_name}"
	location = "${var.location}"
}
`)

var resourceGroupTFVars = []byte(`
resource_group_name = "{{.resourceGroupName}}"
location = "{{.location}}"
env_name = "{{.envName}}"
`)

func NewResourceGroupTemplate() (string, error) {
	return string(resourceGroupTF), nil
}

func NewResourceGroupInput(resourceGroupName, location, envName string) (string, error) {
	return helpers.Render(resourceGroupTFVars, map[string]interface{}{
		"resourceGroupName": resourceGroupName,
		"location":          location,
		"envName":           envName,
	})
}
