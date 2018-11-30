package templates

import (
	artifacts "github.com/Microsoft/kunlun/artifacts"
	"github.com/Microsoft/kunlun/common/helpers"
)

var loadBalancerTF = []byte(`
resource "azurerm_public_ip" "{{.lbName}}_public_ip" {
	name                         = "${var.env_name}-{{.lbName}}-public-ip"
	location                     = "${var.location}"
	resource_group_name          = "${azurerm_resource_group.kunlun_resource_group.name}"
	public_ip_address_allocation = "static"
	sku                          = "${var.{{.lbName}}_al_sku}"
	{{if .haveDomainName -}}
	domain_name_label = "{{var.{{.lbName}}_al_domain_name_label}}"
	{{- end}}
}

resource "azurerm_lb" "{{.lbName}}" {
	name                = "${var.env_name}-{{.lbName}}"
	location            = "${var.location}"
	resource_group_name = "${azurerm_resource_group.kunlun_resource_group.name}"
	sku                 = "${var.{{.lbName}}_al_sku}"
	frontend_ip_configuration {
	  name                 = "{{.lbName}}-public-ip"
	  public_ip_address_id = "${azurerm_public_ip.{{.lbName}}_public_ip.id}"
	}
  }

variable "{{.lbName}}_al_sku" {}
{{if .haveDomainName -}}
variable "{{.lbName}}_al_domain_name_label" {}
{{- end}}
`)

var loadBalancerTFVars = []byte(`
{{.lbName}}_al_sku = "{{.al_sku}}"
{{if .haveDomainName -}}
{{.lbName}}_al_domain_name_label = "{{.al_domain_name_label}}"
{{- end}}
`)

var loadBalancerBackendAddressPoolTF = []byte(`
resource "azurerm_lb_backend_address_pool" "{{.backendAddressPoolName}}" {
	resource_group_name = "${azurerm_resource_group.kunlun_resource_group.name}"
	loadbalancer_id     = "${azurerm_lb.{{.lbName}}.id}"
	name                = "{{.lbName}}-{{.backendAddressPoolName}}"
}
`)

var loadBalancerHealthProbeTF = []byte(`
resource "azurerm_lb_probe" "{{.healthProbeName}}" {
	resource_group_name = "${azurerm_resource_group.kunlun_resource_group.name}"
	loadbalancer_id     = "${azurerm_lb.{{.lbName}}.id}"
	name                = "{{.lbName}}-{{.healthProbeName}}"
	protocol            = "${var.{{.healthProbeName}}_alp_protocol}"
	{{if .haveRequestPath -}}
	request_path        = "${var.{{.healthProbeName}}_alp_request_path}"
	{{- end}}
	port                = "${var.{{.healthProbeName}}_alp_port}"
}

variable "{{.healthProbeName}}_alp_protocol" {}
variable "{{.healthProbeName}}_alp_port" {}
{{if .haveRequestPath -}}
variable "{{.healthProbeName}}_alp_request_path" {}
{{- end}}
`)

var loadBalancerHealthProbeTFVars = []byte(`
{{.healthProbeName}}_alp_protocol = "{{.alp_protocol}}"
{{.healthProbeName}}_alp_port = "{{.alp_port}}"
{{if .haveRequestPath -}}
{{.healthProbeName}}_alp_request_path = "{{.alp_request_path}}"
{{- end}}
`)

var loadBalancerRuleTF = []byte(`
resource "azurerm_lb_rule" "{{.ruleName}}" {
	resource_group_name            = "${azurerm_resource_group.kunlun_resource_group.name}"
	loadbalancer_id                = "${azurerm_lb.{{.lbName}}.id}"
	name                           = "{{.ruleName}}"
	protocol                       = "${var.{{.ruleName}}_alr_protocol}"
	frontend_port                  = "${var.{{.ruleName}}_alr_frontend_port}"
	backend_port                   = "${var.{{.ruleName}}_alr_backend_port}"
	frontend_ip_configuration_name = "{{.lbName}}-public-ip"
	backend_address_pool_id        = "${azurerm_lb_backend_address_pool.{{.backendAddressPoolName}}.id}"
	probe_id                       = "${azurerm_lb_probe.{{.healthProbeName}}.id}"
}

variable "{{.ruleName}}_alr_protocol" {}
variable "{{.ruleName}}_alr_frontend_port" {}
variable "{{.ruleName}}_alr_backend_port" {}
`)

var loadBalancerRuleTFVars = []byte(`
{{.ruleName}}_alr_protocol = "{{.alr_protocol}}"
{{.ruleName}}_alr_frontend_port = "{{.alr_frontend_port}}"
{{.ruleName}}_alr_backend_port = "{{.alr_backend_port}}"
`)

func NewLoadBalancerTemplate(lb artifacts.LoadBalancer) (string, error) {
	tf := ""

	lbTF, err := helpers.Render(loadBalancerTF, getLoadBalancerTFParams(lb))
	if err != nil {
		return "", err
	}
	tf += lbTF

	for _, lbbap := range lb.BackendAddressPools {
		lbbapTF, err := helpers.Render(
			loadBalancerBackendAddressPoolTF,
			getLoadBalancerBackendAddressPoolTFParams(lbbap, lb.Name))
		if err != nil {
			return "", err
		}
		tf += lbbapTF
	}

	for _, lbhp := range lb.HealthProbes {
		lbhpTF, err := helpers.Render(
			loadBalancerHealthProbeTF,
			getLoadBalancerHealthProbeTFParams(lbhp, lb.Name),
		)
		if err != nil {
			return "", err
		}
		tf += lbhpTF
	}

	for _, lbr := range lb.Rules {
		lbrTF, err := helpers.Render(
			loadBalancerRuleTF,
			getLoadBalancerRuleTFParams(lbr, lb.Name),
		)
		if err != nil {
			return "", err
		}
		tf += lbrTF
	}

	return tf, nil
}

func NewLoadBalancerInput(lb artifacts.LoadBalancer) (string, error) {
	tfVars := ""

	lbVars, err := helpers.Render(loadBalancerTFVars, getLoadBalancerTFVarsParams(lb))
	if err != nil {
		return "", err
	}
	tfVars += lbVars

	for _, lbhp := range lb.HealthProbes {
		lbhpVars, err := helpers.Render(
			loadBalancerHealthProbeTFVars,
			getLoadBalancerHealthProbeTFVarsParams(lbhp),
		)
		if err != nil {
			return "", err
		}
		tfVars += lbhpVars
	}

	for _, lbr := range lb.Rules {
		lbrVars, err := helpers.Render(
			loadBalancerRuleTFVars,
			getLoadBalancerRuleTFVarsParams(lbr),
		)
		if err != nil {
			return "", err
		}
		tfVars += lbrVars
	}
	return tfVars, nil
}

func getLoadBalancerTFParams(lb artifacts.LoadBalancer) map[string]interface{} {
	return map[string]interface{}{
		"lbName":         lb.Name,
		"haveDomainName": lb.DomainName != "",
	}
}

func getLoadBalancerTFVarsParams(lb artifacts.LoadBalancer) map[string]interface{} {
	return map[string]interface{}{
		"lbName":               lb.Name,
		"al_sku":               lb.SKU,
		"haveDomainName":       lb.DomainName != "",
		"al_domain_name_label": lb.DomainName,
	}
}

func getLoadBalancerBackendAddressPoolTFParams(
	lbbap artifacts.LoadBalancerBackendAddressPool,
	lbName string,
) map[string]interface{} {
	return map[string]interface{}{
		"lbName":                 lbName,
		"backendAddressPoolName": lbbap.Name,
	}
}

func getLoadBalancerHealthProbeTFParams(
	lbhp artifacts.LoadBalancerHealthProbe,
	lbName string,
) map[string]interface{} {
	return map[string]interface{}{
		"lbName":          lbName,
		"healthProbeName": lbhp.Name,
		"haveRequestPath": lbhp.Protocol == "Http" || lbhp.Protocol == "Https",
	}
}

func getLoadBalancerHealthProbeTFVarsParams(
	lbhp artifacts.LoadBalancerHealthProbe,
) map[string]interface{} {
	return map[string]interface{}{
		"healthProbeName":  lbhp.Name,
		"haveRequestPath":  lbhp.Protocol == "Http" || lbhp.Protocol == "Https",
		"alp_protocol":     lbhp.Protocol,
		"alp_request_path": lbhp.RequestPath,
		"alp_port":         lbhp.Port,
	}
}

func getLoadBalancerRuleTFParams(
	lbr artifacts.LoadBalancerRule,
	lbName string,
) map[string]interface{} {
	return map[string]interface{}{
		"lbName":                 lbName,
		"ruleName":               lbr.Name,
		"backendAddressPoolName": lbr.BackendAddressPoolName,
		"healthProbeName":        lbr.HealthProbeName,
	}
}

func getLoadBalancerRuleTFVarsParams(
	lbr artifacts.LoadBalancerRule,
) map[string]interface{} {
	return map[string]interface{}{
		"ruleName":          lbr.Name,
		"alr_protocol":      lbr.Protocol,
		"alr_frontend_port": lbr.FrontendPort,
		"alr_backend_port":  lbr.BackendPort,
	}
}
