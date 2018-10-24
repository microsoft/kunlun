package templates

import (
	artifacts "github.com/Microsoft/kunlun/artifacts"
	"github.com/Microsoft/kunlun/common/helpers"
)

var networkSecurityGroupTemplate = []byte(`
resource "azurerm_network_security_group" "{{.nsgName}}" {
	name                = "${var.env_name}-{{.nsgName}}"
	location            = "${var.location}"
	resource_group_name = "${azurerm_resource_group.kunlun_resource_group.name}"
}

`)

var networkSecurityRuleTemplate = []byte(`
resource "azurerm_network_security_rule" "{{.nsrName}}" {
	name                        = "{{.nsrName}}"
	priority                    = "${var.{{.nsrName}}_ansr_priority}"
	direction                   = "${var.{{.nsrName}}_ansr_direction}"
	access                      = "${var.{{.nsrName}}_ansr_access}"
	protocol                    = "${var.{{.nsrName}}_ansr_protocol}"
	source_port_range           = "${var.{{.nsrName}}_ansr_source_port_range}"
	destination_port_range      = "${var.{{.nsrName}}_ansr_destination_port_range}"
	source_address_prefix       = "${var.{{.nsrName}}_ansr_source_address_prefix}"
	destination_address_prefix  = "${var.{{.nsrName}}_ansr_destination_address_prefix}"
	resource_group_name         = "${azurerm_resource_group.kunlun_resource_group.name}"
	network_security_group_name = "${azurerm_network_security_group.{{.nsgName}}.name}"
}

variable "{{.nsrName}}_ansr_priority" {}
variable "{{.nsrName}}_ansr_direction" {}
variable "{{.nsrName}}_ansr_access" {}
variable "{{.nsrName}}_ansr_protocol" {}
variable "{{.nsrName}}_ansr_source_port_range" {}
variable "{{.nsrName}}_ansr_destination_port_range" {}
variable "{{.nsrName}}_ansr_source_address_prefix" {}
variable "{{.nsrName}}_ansr_destination_address_prefix" {}
`)

var networkSecurityRuleInput = []byte(`
{{.nsrName}}_ansr_priority = "{{.ansr_priority}}"
{{.nsrName}}_ansr_direction = "{{.ansr_direction}}"
{{.nsrName}}_ansr_access = "{{.ansr_access}}"
{{.nsrName}}_ansr_protocol = "{{.ansr_protocol}}"
{{.nsrName}}_ansr_source_port_range = "{{.ansr_source_port_range}}"
{{.nsrName}}_ansr_destination_port_range = "{{.ansr_destination_port_range}}"
{{.nsrName}}_ansr_source_address_prefix = "{{.ansr_source_address_prefix}}"
{{.nsrName}}_ansr_destination_address_prefix = "{{.ansr_destination_address_prefix}}"
`)

func NewNSGTemplate(nsg artifacts.NetworkSecurityGroup) (string, error) {
	template := ""
	nsgTemplate, err := helpers.Render(networkSecurityGroupTemplate, buildNSGTemplateGoParams(nsg))
	if err != nil {
		return "", err
	}
	template += nsgTemplate

	for _, nsr := range nsg.NetworkSecurityRules {
		nsrTemplate, err := helpers.Render(networkSecurityRuleTemplate, buildNSRTemplateGoParams(nsr, nsg.Name))
		if err != nil {
			return "", err
		}
		template += nsrTemplate
	}
	return template, err
}

func NewNSGInput(nsg artifacts.NetworkSecurityGroup) (string, error) {
	input := ""
	for _, nsr := range nsg.NetworkSecurityRules {
		nsrInput, err := helpers.Render(networkSecurityRuleInput, buildNSRInputGoParams(nsr))
		if err != nil {
			return "", err
		}
		input += nsrInput
	}
	return input, nil
}

func buildNSGTemplateGoParams(nsg artifacts.NetworkSecurityGroup) map[string]interface{} {
	return map[string]interface{}{
		"nsgName": nsg.Name,
	}
}

func buildNSRTemplateGoParams(nsr artifacts.NetworkSecurityRule, nsgName string) map[string]interface{} {
	return map[string]interface{}{
		"nsgName": nsgName,
		"nsrName": nsr.Name,
	}
}

func buildNSRInputGoParams(nsr artifacts.NetworkSecurityRule) map[string]interface{} {
	return map[string]interface{}{
		"nsrName":                         nsr.Name,
		"ansr_priority":                   nsr.Priority,
		"ansr_direction":                  nsr.Direction,
		"ansr_access":                     nsr.Access,
		"ansr_protocol":                   nsr.Protocol,
		"ansr_source_port_range":          nsr.SourcePortRange,
		"ansr_destination_port_range":     nsr.DestinationPortRange,
		"ansr_source_address_prefix":      nsr.SourceAddressPrefix,
		"ansr_destination_address_prefix": nsr.DestinationAddressPrefix,
	}
}
