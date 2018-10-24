package templates

import (
	artifacts "github.com/Microsoft/kunlun/artifacts"
	"github.com/Microsoft/kunlun/common/helpers"
)

var vmssTF = []byte(`
resource "azurerm_virtual_machine_scale_set" "{{.vmGroupName}}" {
	name                      = "${var.env_name}-{{.vmGroupName}}"
	location                  = "${var.location}"
	resource_group_name       = "${azurerm_resource_group.kunlun_resource_group.name}"
	upgrade_policy_mode       = "Manual"
	sku {
	  name     = "${var.{{.vmGroupName}}_avmss_sku_name}"
	  tier     = "Standard"
	  capacity = "${var.{{.vmGroupName}}_avmss_sku_capacity}"
	}

	storage_profile_image_reference {
	  publisher = "${var.{{.vmGroupName}}_avmss_storage_profile_image_reference_publisher}"
	  offer     = "${var.{{.vmGroupName}}_avmss_storage_profile_image_reference_offer}"
	  sku       = "${var.{{.vmGroupName}}_avmss_storage_profile_image_reference_sku}"
	  version   = "${var.{{.vmGroupName}}_avmss_storage_profile_image_reference_version}"
	}

	storage_profile_os_disk {
	  caching           = "${var.{{.vmGroupName}}_avmss_storage_profile_os_disk_caching}"
	  create_option     = "${var.{{.vmGroupName}}_avmss_storage_profile_os_disk_create_option}"
	  managed_disk_type = "${var.{{.vmGroupName}}_avmss_storage_profile_os_disk_managed_disk_type}"
	}

	storage_profile_data_disk {
	  lun               = 0
	  caching           = "${var.{{.vmGroupName}}_avmss_storage_profile_data_disk_caching}"
	  create_option     = "${var.{{.vmGroupName}}_avmss_storage_profile_data_disk_create_option}"
	  managed_disk_type = "${var.{{.vmGroupName}}_avmss_storage_profile_data_disk_managed_disk_type}"
	  disk_size_gb      = "${var.{{.vmGroupName}}_avmss_storage_profile_data_disk_disk_size_gb}"
	}

	os_profile {
	  computer_name_prefix = "{{.vmGroupName}}-"
	  admin_username       = "${var.{{.vmGroupName}}_avmss_os_profile_admin_username}"
	}

	os_profile_linux_config {
	  disable_password_authentication = true
	  ssh_keys {
		  path = "/home/${var.{{.vmGroupName}}_avmss_os_profile_admin_username}/.ssh/authorized_keys"
		  key_data = "${var.{{.vmGroupName}}_avmss_os_profile_linux_config_ssh_keys_key_data}"
	  }
	}

	network_profile {
	  name    = "{{.vmGroupName}}_network_profile"
	  primary = true
	  {{if .nsgName}}
	  network_security_group_id = "${azurerm_network_security_group.{{.nsgName}}.id}"

	  {{end}}
	  ip_configuration {
		name                                   = "{{.vmGroupName}}_ip_configuration"
		primary                                = true
		subnet_id                              = "${azurerm_subnet.{{.subnetName}}.id}"
		{{if .loadBalancerBackendAddressPoolName}}
		load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.{{.loadBalancerBackendAddressPoolName}}.id}"]
		{{end}}
	  }
	}
}

variable "{{.vmGroupName}}_avmss_sku_name" {
	default = "Standard_D2s_v3"
}

variable "{{.vmGroupName}}_avmss_sku_capacity" {}

variable "{{.vmGroupName}}_avmss_storage_profile_image_reference_publisher" {
	default = "Canonical"
}

variable "{{.vmGroupName}}_avmss_storage_profile_image_reference_offer" {
	default = "UbuntuServer"
}

variable "{{.vmGroupName}}_avmss_storage_profile_image_reference_sku" {
	default = "16.04-LTS"
}

variable "{{.vmGroupName}}_avmss_storage_profile_image_reference_version" {
	default = "latest"
}

variable "{{.vmGroupName}}_avmss_storage_profile_os_disk_caching" {
	default = "ReadWrite"
}

variable "{{.vmGroupName}}_avmss_storage_profile_os_disk_create_option" {
	default = "FromImage"
}

variable "{{.vmGroupName}}_avmss_storage_profile_os_disk_managed_disk_type" {
	default = "Standard_LRS"
}

variable "{{.vmGroupName}}_avmss_storage_profile_data_disk_caching" {
	default = "ReadWrite"
}

variable "{{.vmGroupName}}_avmss_storage_profile_data_disk_create_option" {
	default = "Empty"
}

variable "{{.vmGroupName}}_avmss_storage_profile_data_disk_managed_disk_type" {
	default = "Standard_LRS"
}

variable "{{.vmGroupName}}_avmss_storage_profile_data_disk_disk_size_gb" {
	default = 10
}

variable "{{.vmGroupName}}_avmss_os_profile_admin_username"{}

variable "{{.vmGroupName}}_avmss_os_profile_linux_config_ssh_keys_key_data"{}
`)

var vmssTFVars = []byte(`
{{.vmGroupName}}_avmss_sku_name = "{{.avmss_sku_name}}"
{{.vmGroupName}}_avmss_sku_capacity = "{{.avmss_sku_capacity}}"
{{.vmGroupName}}_avmss_storage_profile_image_reference_publisher = "{{.avmss_storage_profile_image_reference_publisher}}"
{{.vmGroupName}}_avmss_storage_profile_image_reference_offer = "{{.avmss_storage_profile_image_reference_offer}}"
{{.vmGroupName}}_avmss_storage_profile_image_reference_sku = "{{.avmss_storage_profile_image_reference_sku}}"
{{.vmGroupName}}_avmss_storage_profile_image_reference_version = "{{.avmss_storage_profile_image_reference_version}}"
{{.vmGroupName}}_avmss_storage_profile_os_disk_caching = "{{.avmss_storage_profile_os_disk_caching}}"
{{.vmGroupName}}_avmss_storage_profile_os_disk_create_option = "{{.avmss_storage_profile_os_disk_create_option}}"
{{.vmGroupName}}_avmss_storage_profile_os_disk_managed_disk_type = "{{.avmss_storage_profile_os_disk_managed_disk_type}}"
{{.vmGroupName}}_avmss_storage_profile_data_disk_caching = "{{.avmss_storage_profile_data_disk_caching}}"
{{.vmGroupName}}_avmss_storage_profile_data_disk_create_option = "{{.avmss_storage_profile_data_disk_create_option}}"
{{.vmGroupName}}_avmss_storage_profile_data_disk_managed_disk_type = "{{.avmss_storage_profile_data_disk_managed_disk_type}}"
{{.vmGroupName}}_avmss_storage_profile_data_disk_disk_size_gb = "{{.avmss_storage_profile_data_disk_disk_size_gb}}"
{{.vmGroupName}}_avmss_os_profile_admin_username = "{{.avmss_os_profile_admin_username}}"
{{.vmGroupName}}_avmss_os_profile_linux_config_ssh_keys_key_data = "{{.avmss_os_profile_linux_config_ssh_keys_key_data}}"

`)

func NewVMSSTemplate(vmss artifacts.VMGroup) (string, error) {
	return helpers.Render(vmssTF, getVMSSTFParams(vmss))
}

func NewVMSSInput(vmss artifacts.VMGroup) (string, error) {
	return helpers.Render(vmssTFVars, getVMSSTFVarsParams(vmss))
}

func getVMSSTFParams(vmss artifacts.VMGroup) map[string]interface{} {
	return map[string]interface{}{
		"vmGroupName":                        vmss.Name,
		"subnetName":                         vmss.NetworkInfos[0].SubnetName,
		"nsgName":                            vmss.NetworkInfos[0].NetworkSecurityGroupName,
		"loadBalancerBackendAddressPoolName": vmss.NetworkInfos[0].LoadBalancerBackendAddressPoolName,
	}
}

func getVMSSTFVarsParams(vmss artifacts.VMGroup) map[string]interface{} {
	return map[string]interface{}{
		"vmGroupName":        vmss.Name,
		"avmss_sku_name":     vmss.SKU,
		"avmss_sku_capacity": vmss.Count,
		"avmss_storage_profile_image_reference_publisher":   vmss.Storage.Image.Publisher,
		"avmss_storage_profile_image_reference_offer":       vmss.Storage.Image.Offer,
		"avmss_storage_profile_image_reference_sku":         vmss.Storage.Image.SKU,
		"avmss_storage_profile_image_reference_version":     vmss.Storage.Image.Version,
		"avmss_storage_profile_os_disk_caching":             vmss.Storage.OSDisk.Caching,
		"avmss_storage_profile_os_disk_create_option":       vmss.Storage.OSDisk.CreateOption,
		"avmss_storage_profile_os_disk_managed_disk_type":   vmss.Storage.OSDisk.ManagedDiskType,
		"avmss_storage_profile_data_disk_caching":           vmss.Storage.DataDisks[0].Caching,
		"avmss_storage_profile_data_disk_create_option":     vmss.Storage.DataDisks[0].CreateOption,
		"avmss_storage_profile_data_disk_managed_disk_type": vmss.Storage.DataDisks[0].ManagedDiskType,
		"avmss_storage_profile_data_disk_disk_size_gb":      vmss.Storage.DataDisks[0].DiskSizeGB,
		"avmss_os_profile_admin_username":                   vmss.OSProfile.AdminName,
		"avmss_os_profile_linux_config_ssh_keys_key_data":   vmss.OSProfile.LinuxConfiguration.SSH.PublicKeys[0],
	}
}
